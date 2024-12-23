package reldb

type Attribute struct {
	IAttribID      int32  `db:"iAttribID" json:"iAttribID"`
	VName          string `db:"vName" json:"vName"`
	IRank          int32  `db:"iRank" json:"iRank"`
	CFilterDisplay string `db:"cFilterDisplay" json:"cFilterDisplay"`
	CStatus        string `db:"cStatus" json:"cStatus"`
}

func (m *Model) AttributeMaster() ([]Attribute, error) {

	query := `SELECT * FROM attribute ORDER BY vName`

	var attributes []Attribute
	if err := m.Select(&attributes, query); err != nil {
		return nil, err
	}

	return attributes, nil
}
