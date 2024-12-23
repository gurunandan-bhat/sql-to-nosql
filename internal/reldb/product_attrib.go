package reldb

import (
	"fmt"
)

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
