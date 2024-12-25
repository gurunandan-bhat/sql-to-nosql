package reldb

import (
	"fmt"
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

type Color struct {
	IColorID uint32 `db:"iColorID" json:"iColorID"`
	VName    string `db:"vName" json:"vName"`
	VColor   string `db:"vColor" json:"vColor"`
	IRank    int64  `db:"iRank" json:"iRank"`
	CStatus  string `db:"cStatus" json:"cStatus"`
}

type ProductAttribute struct {
	IProdAttribID uint32  `db:"iProdAttribID" json:"iProdAttribID" diff:"iProdAttribID"`
	IProdID       uint32  `db:"iProdID" json:"iProdID" diff:"iProdID"`
	IAttribID     uint32  `db:"iAttribID" json:"iAttribID" diff:"iAttribID"`
	VAttribName   *string `db:"vAttribName" json:"vAttribName" diff:"-"`
	VValue        *string `db:"vValue" json:"vValue" diff:"vValue"`
	IPCID         uint32  `db:"iAttribPCID" json:"iAttribPCID" diff:"iPCID"`
	FRetailPrice  float64 `db:"fRetailPrice" json:"fRetailPrice" diff:"fRetailPrice"`
	FRetailOPrice float64 `db:"fRetailOPrice" json:"fRetailOPrice" diff:"fRetailOPrice"`
	FPrice        float64 `db:"fPrice" json:"fPrice" diff:"fPrice"`
	FOPrice       float64 `db:"fOPrice" json:"fOPrice" diff:"fOPrice"`
	CDefault      *string `db:"cDefault" json:"cDefault" diff:"cDefault"`
	CStock        *string `db:"cStock" json:"cStock" diff:"cStock"`
}

type ProductColor struct {
	IPCID              uint32  `json:"iPCID" db:"iPCID"`
	IProdID            uint32  `json:"iColorProdID" db:"iColorProdID"`
	IColorID           uint32  `json:"iColorID" db:"iColorID"`
	VColorName         *string `json:"vColorName" db:"vColorName" diff:"-"`
	FColorRetailPrice  float64 `json:"fColorRetailPrice" db:"fColorRetailPrice"`
	FColorRetailOPrice float64 `json:"fColorRetailOPrice" db:"fColorRetailOPrice"`
	FColorPrice        float64 `json:"fColorPrice" db:"fColorPrice"`
	FColorOPrice       float64 `json:"fColorOPrice" db:"fColorOPrice"`
	CColorDefault      *string `json:"cColorDefault" db:"cColorDefault"`
	CStatus            *string `json:"cStatus" db:"cStatus"`
}

type ProductColorAttribute struct {
	ProductColor
	ProductAttributes []ProductAttribute `json:"productAttributes"`
}

type dbPCARow struct {
	ProductColor
	ProductAttribute
}

func (m *Model) ProductAttributes(iProdID uint32) ([]ProductAttribute, error) {

	if iProdID == 0 {
		return nil, nil
	}

	var productAttribs []ProductAttribute
	query := `SELECT
			    pa.iProdAttribID,
				pa.iProdID,
				pa.iAttribID,
    			a.vName as vAttribName,
				pa.vValue,
				pa.iPCID iAttribPCID,
				pa.fRetailPrice,
				pa.fRetailOPrice,
				pa.fPrice,
				pa.fOPrice,
				pa.cDefault,
				pa.cStock
			FROM
    			product_attrib pa
			JOIN attribute a
			ON pa.iAttribID = a.iAttribID
			WHERE
    			iProdID = ? AND
    			NOT (vValue = '' AND fPrice = 0) AND
    			iPCID = 0
			ORDER BY pa.iProdAttribID`

	if err := m.Select(&productAttribs, query, iProdID); err != nil {
		fmt.Printf("\n\nerror fetching product attributes: %s\n\n", err)
		return nil, err
	}

	return productAttribs, nil
}

func (m *Model) ProductColors(iProdID uint32) ([]ProductColor, error) {

	query := `SELECT
				pc.iPCID,
				pc.iProdID iColorProdID,
				pc.iColorID,
				c.vName vColorName,
				pc.fColorRetailPrice,
				pc.fColorRetailOPrice,
				pc.fColorPrice,
				pc.fColorOPrice,
				pc.cColorDefault,
				pc.cStatus
			FROM
				product_color pc
			JOIN color c ON
				pc.iColorID = c.iColorID
			WHERE
				pc.iProdID = ? AND 
				pc.iPCID NOT IN (SELECT iPCID FROM product_attrib where iProdID = ?)
			ORDER BY
				pc.iColorID`

	pcRows := []ProductColor{}

	// Tell and exit if no rows found
	if err := m.Select(&pcRows, query, iProdID, iProdID); err != nil {
		return pcRows, fmt.Errorf("error scanning rows: %w", err)
	}

	return pcRows, nil
}

func (m *Model) ProductColorAttributes(iProdID uint32) ([]ProductColorAttribute, error) {

	query := `SELECT
				pc.iPCID,
				pc.iProdID iColorProdID,
				pc.iColorID,
				c.vName vColorName,
				pc.fColorRetailPrice,
				pc.fColorRetailOPrice,
				pc.fColorPrice,
				pc.fColorOPrice,
				pc.cColorDefault,
				pc.cStatus,
				pa.iProdAttribID,
				pa.iProdID,
				pa.iAttribID,
				a.vName vAttribName,
				pa.vValue vValue,
				pa.iPCID iAttribPCID,
				pa.fRetailPrice,
				pa.fRetailOPrice,
				pa.fPrice,
				pa.fOPrice,
				pa.cDefault,
				pa.cStock
			FROM
				product_color pc
			JOIN color c ON
				pc.iColorID = c.iColorID
			JOIN product_attrib pa ON
				(pa.iPCID = pc.iPCID AND pa.iProdID = pc.iProdID)
			JOIN attribute a ON
				pa.iAttribID = a.iAttribID
			WHERE
				pc.iProdID = ?
			ORDER BY
				pc.iColorID,
				pa.iProdAttribID`

	pcRows := []dbPCARow{}

	// Tell and exit if no rows found
	cas := []ProductColorAttribute{}
	if err := m.Select(&pcRows, query, iProdID); err != nil {
		return cas, fmt.Errorf("error scanning rows: %w", err)
	}

	if len(pcRows) == 0 {
		return cas, nil
	}

	var lastPCID uint32 = 0
	for _, pcRow := range pcRows {
		cas = addProductColorAttribute(cas, lastPCID, pcRow)
		lastPCID = pcRow.ProductColor.IPCID
	}

	return cas, nil
}

func addProductColorAttribute(cas []ProductColorAttribute, lastPCID uint32, pcRow dbPCARow) []ProductColorAttribute {

	if lastPCID != pcRow.ProductColor.IPCID {
		// We have a new color
		pAttribs := make([]ProductAttribute, 0)
		if pcRow.IProdAttribID > 0 {
			pAttribs = append(pAttribs, pcRow.ProductAttribute)
		}
		cas = append(cas, ProductColorAttribute{
			ProductColor:      pcRow.ProductColor,
			ProductAttributes: pAttribs,
		})

		return cas
	}

	// The color has not changed so just add the new attribute to the existing
	// attributes slice
	idx := len(cas) - 1
	cas[idx].ProductAttributes = append(cas[idx].ProductAttributes, pcRow.ProductAttribute)

	return cas
}
