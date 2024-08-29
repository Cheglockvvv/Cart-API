package connection

import (
	"fmt"
	"github.com/Cheglockvvv/Cart-API/config"
	"github.com/jmoiron/sqlx"
)

func GetConnection(config config.DBConfig) (*sqlx.DB, error) {

	postgresConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.DBName, config.SSLMode)

	db, err := sqlx.Connect("pgx", postgresConnectionString)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	return db, nil
}
