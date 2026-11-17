package reldb

import (
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
