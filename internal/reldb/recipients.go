package reldb

import "fmt"

type Recipient struct {
	PK     string `db:"-"`
	SK     string `db:"-"`
	Set    string `db:"sset"`
	EMail  string `db:"email"`
	Name   string `db:"name"`
	Status string `db:"status"`
}

func (m *Model) Recipients() ([]Recipient, error) {

	qry := `select 
		s.vRecipientSetName sset,
		r.vEmail email,
		r.vFullname name,
		r.cStatus status
	from 
		recipient_set s join
		recipient r on s.iRecipientSetID = r.iRecipientSetID
	order by 1, 4`

	rcpts := []Recipient{}
	if err := m.Select(&rcpts, qry); err != nil {
		return rcpts, fmt.Errorf("error fetching recipients: %w", err)
	}

	return rcpts, nil
}
