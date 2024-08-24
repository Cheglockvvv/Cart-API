package repository

import (
	"Cart-API/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type PostgresCart struct {
	DB *sqlx.DB
}

func (r *PostgresCart) Init(dataSourceName string) error {
	db, err := sqlx.Connect("postgres", dataSourceName)

	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	r.DB = db
	return nil
}

func (r *PostgresCart) CreateCart() (string, error) {
	const createCartQuery = `INSERT INTO cart VALUES (DEFAULT) RETURNING id`

	var id string
	err := r.DB.QueryRowx(createCartQuery).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create a cart: %w", err)
	}

	log.Println(id)

	return id, nil
}

func (r *PostgresCart) GetCartByID(id string) (*models.Cart, error) {
	const getCartByIDQuery = `SELECT ci.id, ci.cart_id, ci.name, ci.quantity FROM cart_item ci WHERE ci.cart_id = $1`
	rows, err := r.DB.Queryx(getCartByIDQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	defer rows.Close()

	cart := &models.Cart{ID: id}
	items := make([]models.CartItem, 0)

	for rows.Next() {
		var row models.CartItem
		rErr := rows.StructScan(&row)

		if rErr != nil {
			return nil, fmt.Errorf("failed to scan row: %w", rErr)
		}
		items = append(items, row)

		log.Printf("ID: %s, cart_ID: %s, name: %s, quantity: %d\n",
			row.ID, row.CartID, row.Name, row.Quantity)
	}

	return cart, nil
}

func (r *PostgresCart) AddItemToCart(cartID, name string, quantity int) (string, error) {
	const addItemToCartQuery = `INSERT INTO cart_item (cart_id, name, quantity) VALUES ($1, $2, $3) RETURNING id`
	var itemID string
	err := r.DB.QueryRowx(addItemToCartQuery, cartID, name, quantity).Scan(&itemID)

	if err != nil {
		return "", fmt.Errorf("failed to add item to cart: %w", err)
	}

	return itemID, nil
}

func (r *PostgresCart) RemoveItemFromCart(cartID, itemID string) error {
	const removeItemFromCartQuery = `DELETE FROM cart_item WHERE cart_id = $1 AND id = $2`
	err := r.DB.QueryRowx(removeItemFromCartQuery, cartID, itemID)

	if err != nil {
		return fmt.Errorf("failed to remove item from cart: %w", err)
	}

	return nil
}
