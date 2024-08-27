package repository

import (
	"Cart-API/app/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CartItem struct {
	DB *sqlx.DB
}

func InitCartItem(db *sqlx.DB) *CartItem {
	cartItem := &CartItem{DB: db}

	return cartItem
}

func (c *CartItem) AddItemToCart(cartID, name string, quantity int) (string, error) {

	const query = `INSERT INTO cart_item (cart_id, product, quantity)
								VALUES ($1, $2, $3)
								ON CONFLICT (cart_id, product)
								DO UPDATE SET quantity = cart_item.quantity +
								    EXCLUDED.quantity
								RETURNING id`
	var itemID string
	err := c.DB.QueryRowx(query, cartID, name, quantity).Scan(&itemID)

	if err != nil {
		return "", fmt.Errorf("c.DB.QueryRowx: %w", err)
	}

	return itemID, nil
}

func (c *CartItem) RemoveItemFromCart(cartID, itemID string) error {
	const query = `DELETE FROM cart_item WHERE cart_id = $1 AND id = $2`

	_, err := c.DB.Exec(query, cartID, itemID)
	if err != nil {
		return fmt.Errorf("c.DB.Exec: %w", err)
	}

	return nil
}

func (c *CartItem) ItemIsAvailable(id string) (bool, error) {
	const checkItem = `SELECT id FROM cart_item WHERE id = $1`
	result, err := c.DB.Exec(checkItem, id)
	if err != nil {
		return false, fmt.Errorf("c.DB.Exec: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("result.RowsAffected: %w", err)
	}

	if count != 1 {
		return false, nil
	}

	return true, nil
}

func (c *CartItem) GetItemByID(id string) (models.CartItem, error) {
	const query = `SELECT ci.id, ci.cart_id, ci.product, ci.quantity 
								FROM cart_item ci 
								WHERE ci.id = $1`
	result := c.DB.QueryRowx(query, id)
	item := cartItemEntity{}
	err := result.StructScan(&item)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("result.StructScan: %w", err)
	}

	convertedItem := cartItemConvert(item)

	return convertedItem, nil
}
