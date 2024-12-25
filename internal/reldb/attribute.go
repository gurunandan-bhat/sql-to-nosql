package reldb

type Attribute struct {
	IAttribID      int32  `db:"iAttribID" json:"iAttribID"`
	VName          string `db:"vName" json:"vName"`
	IRank          int32  `db:"iRank" json:"iRank"`
	CFilterDisplay string `db:"cFilterDisplay" json:"cFilterDisplay"`
	CStatus        string `db:"cStatus" json:"cStatus"`
}
