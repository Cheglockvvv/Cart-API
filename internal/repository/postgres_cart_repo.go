package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostgresCartRepo struct {
	DB *sqlx.DB
}

const (
	createCartQuery         = `INSERT INTO cart VALUES (DEFAULT) RETURNING id`
	getCartByIDQuery        = `SELECT ci.id, ci.cart_id, ci.name, ci.quantity FROM cart_item ci WHERE ci.cart_id = $1`
	addItemToCartQuery      = `INSERT INTO cart_item (cart_id, name, quantity) VALUES ($1, $2, $3) RETURNING id`
	removeItemFromCartQuery = `DELETE FROM cart_item WHERE cart_id = $1 AND id = $2`
)

func (repo *PostgresCartRepo) Init(dataSourceName string) error {
	db, err := sqlx.Connect("postgres", dataSourceName)

	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	repo.DB = db
	return nil
}

func (repo *PostgresCartRepo) CreateCart() (string, error) {
	err := repo.DB.QueryRowx(createCartQuery)
	if err != nil {
		return fmt.Errorf("Failed to create a cart: %w", err)
	}

}
