package pkg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectToPSQL(address string) (*sql.DB, error) {
	db, err := sql.Open("postgres", address)
	if err != nil {
		return nil, fmt.Errorf("pkg: ConnectToPSQL: %s", err.Error())
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pkg: ConnectToPSQL: %s", err.Error())
	}
	return db, nil
}
