package migrations

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Up(db *sqlx.DB, migrationsLocation string) error {
	m, err := initMigrator(db, migrationsLocation)
	if err != nil {
		return fmt.Errorf("initMigrator: %w", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("db is up to date: %w", err)
	}

	return nil
}

func Down(db *sqlx.DB, migrationsLocation string) error {
	m, err := initMigrator(db, migrationsLocation)
	if err != nil {
		return fmt.Errorf("initMigrator: %w", err)
	}

	err = m.Down()
	if err != nil {
		return fmt.Errorf("db is already down: %w", err)
	}

	return nil
}

func initMigrator(db *sqlx.DB, migrationsLocation string) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("postgres.WithInstance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsLocation),
		"postgres", driver)

	if err != nil {
		return nil, fmt.Errorf("migrate.NewWithDatabaseInstance: %w", err)
	}

	return m, nil
}
