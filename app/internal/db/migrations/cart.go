package migrations

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Up(db *sqlx.DB) error {
	m, err := initMigrator(db)
	if err != nil {
		return fmt.Errorf("initMigrator")
	}
	m.Up()

	return nil
}

func Down(db *sqlx.DB) error {
	m, err := initMigrator(db)
	if err != nil {
		return fmt.Errorf("initMigrator")
	}
	m.Down()

	return nil
}

func initMigrator(db *sqlx.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres.WithInstance")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://.",
		"postgres", driver)

	if err != nil {
		return nil, fmt.Errorf("migrate.NewWithDatabaseInstance")
	}

	return m, nil
}