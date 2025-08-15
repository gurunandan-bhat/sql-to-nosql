package reldb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
)

type ProductImgForm struct {
	IProdID           int      `json:"iProdID,omitempty" schema:"iProdID,omitempty" db:"iProdID"`
	VImage            string   `json:"vImage,omitempty" schema:"vImage,omitempty" db:"vImage"`
	VImage_AltTag     string   `json:"vImage_AltTag,omitempty" schema:"vImage_AltTag,omitempty" db:"vImage_AltTag"`
	VAddlImages       []string `json:"vAddlImages,omitempty" schema:"vAddlImages,omitempty" db:"vAddlImages"`
	VAddlImage_AltTag string   `json:"vAddlImage_AltTag,omitempty" schema:"vAddlImage_AltTag,omitempty" db:"vAddlImage_AltTag"`
	ToDelete          []int    `json:"toDelete,omitempty" schema:"toDelete,omitempty" db:"toDelete"`
}

type DbImage struct {
	IProdImageID int32   `json:"iProdImageID,omitempty" db:"iProdImageID"`
	VName        *string `json:"vName,omitempty" db:"vName"`
	VAltTag      *string `json:"vAlt_Tag,omitempty" db:"vAltTag"`
	CStatus      *string `json:"cStatus,omitempty" db:"cStatus"`
}

type ProductImage struct {
	VType string `json:"vType,omitempty" db:"vType"`
	DbImage
}

var decoder schema.Decoder = *schema.NewDecoder()

func (m *Model) ProductImages(productID int32) ([]ProductImage, error) {

	pImages := []ProductImage{}

	query := `SELECT
					'main'        as vType,
					0             as iProdImageID,
					vImage        as vName,
					vImage_AltTag as vAltTag,
					'A'           as cStatus
				FROM product
				WHERE iProdID = ?
			UNION
				SELECT
					'additional' as vType,
					iProdImageID,
					vPic         as vName,
					vTitle       as vAltTag,
					cStatus
				FROM product_images
				WHERE iProdID = ?`

	if err := m.Select(&pImages, query, productID, productID); err != nil {
		return nil, err
	}

	return pImages, nil
}

var delOtherQry = "DELETE FROM product_images WHERE iProdImageID in (?)"

func (m *Model) DeleteOtherImages(ctx context.Context, delIDs []int32) error {

	qry, args, err := sqlx.In(delOtherQry, delIDs)
	if err != nil {
		return err
	}

	qry = m.Rebind(qry)
	_, err = m.ExecContext(ctx, qry, args...)

	return err
}

func (m *Model) SaveProductImages(formValues map[string][]string) (ProductImgForm, error) {

	productImages := ProductImgForm{}
	if err := decoder.Decode(&productImages, formValues); err != nil {
		return ProductImgForm{}, err
	}

	newParams := []string{}
	newValues := map[string]any{}
	if productImages.VImage != "" {
		newParams = append(newParams, "vImage = :vImage, vSmallImage = :vImage")
		newValues["vImage"] = productImages.VImage
	}
	if productImages.VImage_AltTag != "" {
		newParams = append(newParams, "vImage_AltTag = :vImage_AltTag, vSmallImage_AltTag = :vImage_AltTag")
		newValues["vImage_AltTag"] = productImages.VImage_AltTag
	}

	tx, err := m.Beginx()
	if err != nil {
		return ProductImgForm{}, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Fatal(err)
		}
	}()

	if len(newParams) > 0 {

		mainImgQry := fmt.Sprintf("UPDATE product SET %s WHERE iProdID = :iProdID", strings.Join(newParams, ", "))
		newValues["iProdID"] = uint(productImages.IProdID)

		_, err = tx.NamedExec(mainImgQry, newValues)
		if err != nil {
			return ProductImgForm{}, err
		}
	}

	if len(productImages.ToDelete) > 0 {

		delOtherQry := `DELETE FROM product_images WHERE iProdImageID in (?)`
		qry, args, _ := sqlx.In(delOtherQry, productImages.ToDelete)
		qry = m.Rebind(qry)

		_, err := m.Exec(qry, args...)
		if err != nil {
			return ProductImgForm{}, err
		}
	}

	if len(productImages.VAddlImages) > 0 {

		replaceQry := `REPLACE INTO product_images 
						(iProdID, vPic, vTitle, cStatus)
					VALUES (?, ?, ?, ?)`

		altTag := "Additional Image"
		if productImages.VAddlImage_AltTag != "" {
			altTag = productImages.VAddlImage_AltTag
		}

		for _, img := range productImages.VAddlImages {
			_, err := tx.Exec(replaceQry, productImages.IProdID, img, altTag, "A")
			if err != nil {
				return ProductImgForm{}, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return ProductImgForm{}, err
	}

	return productImages, nil
}
