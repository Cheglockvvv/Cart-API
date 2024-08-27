package connection

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

func GetConnection(connectionString string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	return db, nil
}
