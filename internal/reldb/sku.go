package reldb

import "fmt"

type SKUAttrib struct {
	VAttribName  string
	VAttribValue string
}
type SKU struct {
	Attributes    []SKUAttrib
	FRetailPrice  float64
	FRetailOPrice float64
	FPrice        float64
	FOPrice       float64
	CDefault      string
	CStock        string
}

func (m *Model) ProductSKUs(iProdID uint32) ([]SKU, error) {

	attribs, err := m.ProductAttributes(iProdID, true)
	if err != nil {
		return nil, fmt.Errorf("error fetching attributes for product %d: %s", iProdID, err)
	}

	colorAttribs, err := m.ProductColorAttributes(iProdID, true)
	if err != nil {
		return nil, fmt.Errorf("error fetching color attributes for product %d: %s", iProdID, err)
	}

	skus := []SKU{}
	for _, attrib := range attribs {

		a := SKUAttrib{*attrib.VAttribName, *attrib.VValue}
		skus = append(skus, SKU{
			Attributes:    []SKUAttrib{a},
			FRetailPrice:  attrib.FRetailPrice,
			FRetailOPrice: attrib.FRetailOPrice,
			FPrice:        attrib.FPrice,
			FOPrice:       attrib.FOPrice,
			CDefault:      *attrib.CDefault,
			CStock:        *attrib.CStock,
		})
	}

	for _, cattrib := range colorAttribs {
		skuAttrib := SKUAttrib{"Color", *cattrib.VColorName}
		for _, attrib := range cattrib.ProductAttributes {
			d := "N"
			if *cattrib.CColorDefault == "Y" && *attrib.CDefault == "Y" {
				d = "Y"
			}
			s := "N"
			if *cattrib.CStatus == "A" && *attrib.CStock == "Y" {
				s = "Y"
			}
			skus = append(skus, SKU{
				Attributes:    append([]SKUAttrib{skuAttrib}, SKUAttrib{*attrib.VAttribName, *attrib.VValue}),
				FRetailPrice:  attrib.FRetailPrice,
				FRetailOPrice: attrib.FRetailOPrice,
				FPrice:        attrib.FPrice,
				FOPrice:       attrib.FOPrice,
				CDefault:      d,
				CStock:        s,
			})
		}
	}

	return skus, nil
}
