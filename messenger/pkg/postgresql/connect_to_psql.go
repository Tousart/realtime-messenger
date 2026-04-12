package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectToPSQL(user, password, host, dbname, sslmode string, port int) (*sql.DB, error) {
	const op = "pkg: postgresql: ConnectToPSQL:"

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode))
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return db, nil
}
