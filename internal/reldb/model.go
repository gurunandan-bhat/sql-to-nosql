package reldb

import (
	"encoding/json"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Model struct {
	*sqlx.DB
}

type TableDef struct {
	Table    string `db:"Table"`
	TableDef string `db:"Create Table"`
}

func NewModel(cfg *Config) (*Model, error) {

	dbAuth := cfg.Db["site"]
	dbConfig := mysql.Config{
		User:                 dbAuth.User,
		Passwd:               dbAuth.Password,
		Net:                  "tcp",
		Addr:                 dbAuth.Host,
		DBName:               dbAuth.Name,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	dbHandle, err := sqlx.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		return nil, err
	}

	if pingErr := dbHandle.Ping(); pingErr != nil {
		return nil, pingErr
	}

	return &Model{
		dbHandle,
	}, nil
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func dumpJSON(prefix string, v any) error {

	jsonBytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Errorf("error dumping json of type %t: %w", v, err)
	}

	fmt.Printf("\n\nDumping %s: %s\n\n", prefix, string(jsonBytes))
	return nil
}
