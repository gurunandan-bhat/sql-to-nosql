package reldb

import (
	"errors"
	"fmt"
)

type Color struct {
	IColorID uint32 `db:"iColorID" json:"iColorID"`
	VName    string `db:"vName" json:"vName"`
	VColor   string `db:"vColor" json:"vColor"`
	IRank    int64  `db:"iRank" json:"iRank"`
	CStatus  string `db:"cStatus" json:"cStatus"`
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

func (m *Model) ColorMaster() ([]Color, error) {

	query := `SELECT * FROM color ORDER BY vName`

	var colors []Color
	if err := m.Select(&colors, query); err != nil {
		return nil, err
	}

	return colors, nil
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
				product_colors pc
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

var ErrDefaultColorNotActive = errors.New("default color is not active")
var ErrMultipleDefaultColors = errors.New("more than one default color")
var ErrDefaultColorNotFound = errors.New("no default color found")

func (m *Model) ValidateProductColors(newColors []ProductColor) error {

	if len(newColors) == 0 {
		return nil
	}

	var defaultCount int
	for _, color := range newColors {
		if *color.CColorDefault == "Y" {
			defaultCount++
			if *color.CStatus != "A" {
				return ErrDefaultColorNotActive
			}
		}
	}
	if defaultCount > 1 {
		return ErrMultipleDefaultColors
	}
	if defaultCount == 0 {
		return ErrDefaultColorNotFound
	}
	// TODO: More validations required

	return nil
}

func (m *Model) ValidateProductColorAttributes(newColorAttribs []ProductColorAttribute) error {

	// TODO: Validations required
	return nil
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
				product_colors pc
			JOIN product_attrib pa ON
				(pa.iPCID = pc.iPCID AND pa.iProdID = pc.iProdID)
			JOIN color c ON
				pc.iColorID = c.iColorID
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
