package reldb

import (
	"fmt"
)

type CategorySummary struct {
	IPCatID           uint32              `db:"iPCatID"`
	VName             string              `db:"vName"`
	VURLName          string              `db:"vUrlName"`
	IParentID         uint32              `db:"iParentID"`
	VShortDesc        *string             `db:"vShortDesc"`
	VMenuImage        *string             `db:"vMenuImage"`
	VMenuImage_AltTag *string             `db:"vMenuImage_AltTag"`
	CStatus           string              `db:"cStatus"`
	Attributes        []CategoryAttribute `db:"-"`
	Children          []CategorySummary   `db:"-"`
}

type CategoryAttribute struct {
	IAttribDatID uint32  `db:"iAttribDatID" json:"iAttribDatID" diff:"iAttribDatID"`
	IPCatID      uint32  `db:"iPCatID" json:"iPCatID" diff:"iPCatID"`
	IAttribID    uint32  `db:"iAttribID" json:"iAttribID" diff:"iAttribID"`
	VAttribName  *string `db:"vAttribName" json:"vAttribName" diff:"-"`
	VName        *string `db:"vName" json:"vName" diff:"vName"`
	IRank        int     `db:"iRank" json:"iRank" diff:"iRank"`
}

func (m *Model) CategoryMaster() (map[uint32]*CategorySummary, error) {

	query := `SELECT
				iPCatID,
				vName,
				vUrlName,
				iParentID,
				vShortDesc,
				vMenuImage,
				vMenuImage_AltTag,
				cStatus
			FROM
				prodcat
			ORDER BY
				vName`
	var categories []CategorySummary
	if err := m.Select(&categories, query); err != nil {
		return nil, err
	}

	cMap := make(map[uint32]*CategorySummary)
	for _, category := range categories {
		cMap[category.IPCatID] = &category
	}

	// Add a root category so top level have parents
	root := CategorySummary{
		IPCatID: 0,
		VName:   "Root Category",
	}
	cMap[0] = &root

	return cMap, nil
}

func (m *Model) CategoryAttributes() (map[uint32][]CategoryAttribute, error) {

	catAttrMap := make(map[uint32][]CategoryAttribute)
	var categoryAttribs []CategoryAttribute

	query := `SELECT
			    pca.iAttribDatID,
				pca.iAttribID,
				pca.iPCatID,
    			a.vName as vAttribName,
				pca.vName
			FROM
				prodcat c 
    			LEFT JOIN prodcat_attrib_dat pca ON c.iPCatID = pca.iPCatID
				JOIN attribute a ON pca.iAttribID = a.iAttribID
			
			ORDER BY c.iPCatID`

	if err := m.Select(&categoryAttribs, query); err != nil {
		return nil, err
	}

	for _, attr := range categoryAttribs {

		attrs, exists := catAttrMap[attr.IPCatID]
		if !exists {
			attrs = make([]CategoryAttribute, 0)
		}
		attrs = append(attrs, attr)
		catAttrMap[attr.IPCatID] = attrs
	}

	return catAttrMap, nil
}

func (m *Model) Categories() ([]CategorySummary, error) {

	catSummMap, err := m.CategoryMaster()
	if err != nil {
		return nil, fmt.Errorf("error fetching category summary: %s", err)
	}

	catAttrMap, err := m.CategoryAttributes()
	if err != nil {
		return nil, fmt.Errorf("error fetching category attributes: %s", err)
	}

	// Lets fix the children and the attributes.
	for iPCatID, category := range catSummMap {

		if iPCatID > 0 {
			parent, hasParent := catSummMap[category.IParentID]
			if hasParent {
				if parent.Children == nil {
					parent.Children = make([]CategorySummary, 0)
				}
				parent.Children = append(parent.Children, *category)
			}
		}

		attr, exists := catAttrMap[iPCatID]
		if exists {
			category.Attributes = attr
		}
	}

	result := make([]CategorySummary, 0)
	for _, cat := range catSummMap {
		if cat.IPCatID > 0 { // Dont add the root category
			result = append(result, *cat)
		}
	}

	return result, nil
}
