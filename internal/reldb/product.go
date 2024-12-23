package reldb

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/r3labs/diff/v3"
)

type Product struct {
	IProdID           uint32  `db:"iProdID" json:"iProdID"`
	IPCatID           uint32  `db:"iPCatID" json:"iPCatID"`
	CCode             *string `db:"cCode" json:"cCode"`
	VName             string  `db:"vName" json:"vName"`
	VURLName          string  `db:"vUrlName" json:"vUrlName"`
	VShortDesc        *string `db:"vShortDesc" json:"vShortDesc"`
	VDescription      *string `db:"vDescription" json:"vDescription"`
	FRetailPrice      float64 `db:"fRetailPrice" json:"fRetailPrice"`
	FRetailOPrice     float64 `db:"fRetailOPrice" json:"fRetailOPrice"`
	FShipping         float64 `db:"fShipping" json:"fShipping"`
	FPrice            float64 `db:"fPrice" json:"fPrice"`
	FOPrice           float64 `db:"fOPrice" json:"fOPrice"`
	FActualWeight     float64 `db:"fActualWeight" json:"fActualWeight"`
	FVolumetricWeight float64 `db:"fVolumetricWeight" json:"fVolumetricWeight"`
	VSmallImage       *string `db:"vSmallImage" json:"vSmallImage"`
	VSmallImageAltTag *string `db:"vSmallImage_AltTag" json:"vSmallImage_AltTag"`
	VImage            *string `db:"vImage" json:"vImage"`
	VImageAltTag      *string `db:"vImage_AltTag" json:"vImage_AltTag"`
	CStatus           *string `db:"cStatus" json:"cStatus"`
	VYTID             *string `db:"vYTID" json:"vYTID"`
}

type ProductStatus struct {
	IProdID int32   `json:"iProdID,omitempty"`
	CStatus *string `json:"cStatus,omitempty"`
}

func (m *Model) ProductByID(iProdID uint32) (Product, error) {
	if iProdID == 0 {
		return Product{}, errors.New("invalid product id")
	}

	query := `SELECT
					iProdID,
					iPCatID,
					cCode,
					vName,
					vUrlName,
					vShortDesc,
					vDescription,
					fRetailPrice,
					fRetailOPrice,
					fShipping,
					fPrice,
					fOPrice,
					fActualWeight,
					fVolumetricWeight,
					vSmallImage,
					vSmallImage_AltTag,
					vImage,
					vImage_AltTag,
					cStatus,
					vYTID
				FROM product
				WHERE iProdID = ?
				ORDER BY
					cStatus,
					vName`

	var product Product

	if err := m.QueryRowx(query, iProdID).StructScan(&product); err != nil {
		return Product{}, fmt.Errorf("error retrieving product %d: %w", iProdID, err)
	}

	return product, nil
}

func (m *Model) ProductsByCategoryID(constraint string, categoryID int32) ([]Product, error) {

	var qry_constraint string
	if constraint == "inactive" {
		qry_constraint = " AND cStatus = 'I' "
	}
	query := `SELECT
			iProdID,
			iPCatID,
			cCode,
			vName,
			vUrlName,
			vShortDesc,
			vDescription,
			fPrice,
			fOPrice,
			fActualWeight,
			fVolumetricWeight,
			vSmallImage,
			vSmallImage_AltTag,
			vImage,
			vImage_AltTag,
			cStatus,
			vYTID
		FROM product
		WHERE iPCatID = ?` + qry_constraint +
		` ORDER BY
			cStatus,
			vName`

	rows, err := m.Queryx(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.StructScan(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (m *Model) CreateProduct(newProduct Product) (uint32, error) {

	differ, _ := diff.NewDiffer(diff.SliceOrdering(true))
	changes, err := differ.Diff(Product{}, newProduct)
	if err != nil {
		return 0, fmt.Errorf("error generating diff for product: %w", err)
	}

	toCreateValues := make([]any, len(changes))
	toCreateCols := make([]string, len(changes))
	toBind := make([]string, len(changes))

	tPS := reflect.TypeOf(Product{})
	for i, c := range changes {

		fieldName := c.Path[0]
		f, hasField := tPS.FieldByName(fieldName)
		if !hasField {
			err := fmt.Errorf("field %s not found", fieldName)
			return 0, err
		}

		colName, _ := f.Tag.Lookup("db")
		toCreateCols[i] = colName
		toCreateValues[i] = c.To
		toBind[i] = "?"
	}

	// We always generate and add the URL for a product
	// so there is no need to add vUrlName.

	query := fmt.Sprintf("INSERT INTO product (%s) VALUES (%s)",
		strings.Join(toCreateCols, ", "),
		strings.Join(toBind, ", "),
	)

	result, err := m.Exec(query, toCreateValues...)
	if err != nil {
		return 0, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(insertedID), nil
}

func (m *Model) UpdateProduct(newProduct Product) (int64, error) {

	iProdID := newProduct.IProdID
	if iProdID == 0 {
		return 0, fmt.Errorf("no product ID available to update")
	}
	product, err := m.ProductByID(iProdID)
	if err != nil {
		return 0, fmt.Errorf("error fetching product by ID %d: %w", iProdID, err)
	}

	differ, _ := diff.NewDiffer(diff.SliceOrdering(true))
	changes, err := differ.Diff(product, newProduct)
	if err != nil {
		return 0, fmt.Errorf("error generating diff: %w", err)
	}
	if len(changes) == 0 {
		return int64(iProdID), nil
	}

	toUpdateValues := make([]any, len(changes)+1)
	toUpdateCols := make([]string, len(changes))

	tPS := reflect.TypeOf(Product{})
	for i, c := range changes {
		fieldName := c.Path[0]
		f, _ := tPS.FieldByName(fieldName)
		colName, _ := f.Tag.Lookup("db")
		colValue := c.To

		toUpdateValues[i] = colValue
		toUpdateCols[i] = fmt.Sprintf("%s = ?", colName)
	}

	// Check if shipping has changed and if yes,
	// update fPrices and fOPrices in product_attrib
	// and product_color_assoc

	sChange := changes.Filter([]string{"FShipping"})
	if len(sChange) > 0 {
		newShipping, ok := (sChange[0].To).(float64)
		if !ok {
			return 0, errors.New("detected shipping change but could not determine value")
		}
		if err := m.updateAttribPrices(iProdID, newShipping); err != nil {
			return 0, fmt.Errorf("error updating attribute prices due to changed shipping cost: %w", err)
		}
	}

	query := `UPDATE product SET ` +
		strings.Join(toUpdateCols, ", ") +
		` WHERE iProdID = ?`

	toUpdateValues[len(changes)] = iProdID

	result, err := m.Exec(query, toUpdateValues...)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *Model) updateAttribPrices(iProdID uint32, shipping float64) error {

	qryP := `UPDATE product_attrib
				SET fPrice = fRetailPrice + ?
			WHERE fRetailPrice > 0 AND iProdID = ?`
	_, err := m.Exec(qryP, shipping, iProdID)
	if err != nil {
		return err
	}

	qryOP := `UPDATE product_attrib
				SET fOPrice = fRetailOPrice + ?
			WHERE fRetailOPrice > 0 AND iProdID = ?`

	_, err = m.Exec(qryOP, shipping, iProdID)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) MakeUrlName(name, table string) (string, error) {

	urlName := strcase.ToSnake(name)

	found := false
	url := urlName
	postfix := 0
	maxTries := 100
	for !found {

		exists, err := m.hasUrlName(url, table)
		if !exists || (err != nil) {
			return url, err
		}

		found = false
		postfix++
		url = fmt.Sprintf("%s_%d", urlName, postfix)

		if postfix >= maxTries {
			break
		}
	}

	return "", errors.New("urlName: Max tries exceeded")
}

func (m *Model) hasUrlName(name, table string) (bool, error) {

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE vUrlName = ?", table)
	if err := m.QueryRowx(query, name).Scan(&count); err != nil {
		return true, err
	}

	return count > 0, nil

}

func (m *Model) UpdateProductStatus(status ProductStatus) (int64, error) {

	query := `UPDATE product SET cStatus = ? WHERE iProdID = ?`
	result, err := m.Exec(query, status.CStatus, status.IProdID)
	if err != nil {
		return 0, err
	}

	updated, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return updated, nil
}

func (m *Model) DeleteProduct(iProdID int32) (int64, error) {

	deleteQry := `DELETE FROM product WHERE iProdId = ?`
	result, err := m.Exec(deleteQry, iProdID)
	if err != nil {
		return 0, err
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}
