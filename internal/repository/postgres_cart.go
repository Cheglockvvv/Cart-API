package repository

import (
	"Cart-API/internal/models"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresCart struct {
	DB *sqlx.DB
}

func (r *PostgresCart) Init(connectionString string) error {
	db, err := sqlx.Connect("pgx", connectionString)
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

	if !rows.Next() {
		return models.Cart{}, nil
	}

	cart := models.Cart{ID: id}
	items := make([]models.CartItem, 0)

	for rows.Next() {
		row := models.CartItem{}
		err = rows.StructScan(&row)

		if err != nil {
			return models.Cart{}, fmt.Errorf("rows.StructScan: %w", err)
		}
		items = append(items, row)
	}

	cart.Items = items
	return cart, nil
}

func (r *PostgresCart) AddItemToCart(cartID, name string, quantity int) (string, error) {
	const check = `SELECT id FROM cart WHERE id = $1`
	result, err := r.DB.Exec(check, cartID)
	if err != nil {
		return "", fmt.Errorf("r.DB.Exec: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("result.RowsAffected: %w", err)
	}

	if count != 1 {
		return "", fmt.Errorf("cart not found")
	}

	const query = `INSERT INTO cart_item (cart_id, name, quantity)
								VALUES ($1, $2, $3)
								RETURNING id`
	var itemID string
	err = r.DB.QueryRowx(query, cartID, name, quantity).Scan(&itemID)

	if err != nil {
		return "", fmt.Errorf("r.DB.QueryRowx: %w", err)
	}

	return itemID, nil
}

func (r *PostgresCart) RemoveItemFromCart(cartID, itemID string) error {

	const checkCart = `SELECT id FROM cart WHERE id = $1`
	result, err := r.DB.Exec(checkCart, cartID)
	if err != nil {
		return fmt.Errorf("r.DB.Exec: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if count != 1 {
		return fmt.Errorf("cart not found")
	}

	const checkItem = `SELECT id FROM cart_item WHERE id = $1`
	result, err = r.DB.Exec(checkItem, itemID)
	if err != nil {
		return fmt.Errorf("r.DB.Exec: %w", err)
	}
	count, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected: %w", err)
	}

	if count != 1 {
		return fmt.Errorf("item not found")
	}

	const query = `DELETE FROM cart_item WHERE cart_id = $1 AND id = $2`

	r.DB.QueryRowx(query, cartID, itemID)

	return nil
}
