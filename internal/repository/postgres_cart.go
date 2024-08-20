package repository

import (
	"Cart-API/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostgresCart struct {
	DB *sqlx.DB
}

func (r *PostgresCart) Init(connectionString string) error {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("sqlx.Connect: %w", err)
	}

	r.DB = db
	return nil
}

func (r *PostgresCart) CreateCart() (string, error) {
	const query = `INSERT INTO cart VALUES (DEFAULT) RETURNING id`

	var id string
	err := r.DB.QueryRowx(query).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("r.DB.QueryRowx.Scan: %w", err)
	}

	return id, nil
}

func (r *PostgresCart) GetCartByID(id string) (models.Cart, error) {
	const query = `SELECT ci.id, ci.cart_id, ci.name, ci.quantity 
								FROM cart_item ci 
								WHERE ci.cart_id = $1`

	rows, err := r.DB.Queryx(query, id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("r.DB.Queryx: %w", err)
	}
	defer rows.Close()

	cart := models.Cart{ID: id}
	items := make([]models.CartItem, 0)

	for rows.Next() {
		var row models.CartItem
		err = rows.StructScan(&row)

		if err != nil {
			return models.Cart{}, fmt.Errorf("rows.StructScan: %w", err)
		}
		items = append(items, row)
	}

	return cart, nil
}

func (r *PostgresCart) AddItemToCart(cartID, name string, quantity int) (string, error) {
	const query = `INSERT INTO cart_item (cart_id, name, quantity)
								VALUES ($1, $2, $3)
								RETURNING id`
	var itemID string
	err := r.DB.QueryRowx(query, cartID, name, quantity).Scan(&itemID)

	if err != nil {
		return "", fmt.Errorf("r.DB.QueryRowx: %w", err)
	}

	return itemID, nil
}

func (r *PostgresCart) RemoveItemFromCart(cartID, itemID string) error {
	const query = `DELETE FROM cart_item WHERE cart_id = $1 AND id = $2`

	err := r.DB.QueryRowx(query, cartID, itemID)
	if err != nil {
		return fmt.Errorf("r.DB.QueryRowx: %w", err)
	}

	return nil
}
