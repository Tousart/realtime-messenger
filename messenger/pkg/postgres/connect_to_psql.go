package pkg

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectToPSQL(addr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, fmt.Errorf("pkg: ConnectToPSQL: %s", err.Error())
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pkg: ConnectToPSQL: %s", err.Error())
	}
	return db, nil
}
