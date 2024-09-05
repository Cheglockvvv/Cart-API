package repository

import (
	"context"
	"fmt"
	"github.com/Cheglockvvv/Cart-API/internal/models"
	"github.com/jmoiron/sqlx"
)

type cartItemEntity struct {
	ID       string `db:"id"`
	CartID   string `db:"cart_id"`
	Product  string `db:"product"`
	Quantity int    `db:"quantity"`
}

type CartItem struct {
	DB *sqlx.DB
}

func NewCartItem(db *sqlx.DB) *CartItem {
	cartItem := &CartItem{DB: db}

	return cartItem
}

func (c *CartItem) Create(ctx context.Context, cartID, name string, quantity int) (string, error) {

	const query = `INSERT INTO cart_item (cart_id, product, quantity)
								VALUES ($1, $2, $3)
								ON CONFLICT (cart_id, product)
								DO UPDATE SET quantity = cart_item.quantity +
								    EXCLUDED.quantity
								RETURNING id`
	var itemID string
	err := c.DB.QueryRowxContext(ctx, query, cartID, name, quantity).Scan(&itemID)

	if err != nil {
		return "", fmt.Errorf("c.DB.QueryRowx: %w", err)
	}

	return itemID, nil
}

func (c *CartItem) Delete(ctx context.Context, cartID, itemID string) error {
	const query = `DELETE FROM cart_item WHERE cart_id = $1 AND id = $2`

	_, err := c.DB.ExecContext(ctx, query, cartID, itemID)
	if err != nil {
		return fmt.Errorf("c.DB.Exec: %w", err)
	}

	return nil
}

func (c *CartItem) ItemExists(ctx context.Context, id string) (bool, error) {
	const checkItem = `SELECT EXISTS(SELECT 1 FROM cart_item WHERE id = $1)`
	var exists bool
	err := c.DB.QueryRowxContext(ctx, checkItem, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("c.DB.QueryRowxContext.Scan: %w", err)
	}

	return exists, nil
}

func (c *CartItem) Read(ctx context.Context, id string) (models.CartItem, error) {
	const query = `SELECT ci.id, ci.cart_id, ci.product, ci.quantity 
								FROM cart_item ci 
								WHERE ci.id = $1`
	result := c.DB.QueryRowxContext(ctx, query, id)
	item := cartItemEntity{}
	err := result.StructScan(&item)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("result.StructScan: %w", err)
	}

	convertedItem := cartItemConvert(item)

	return convertedItem, nil
}

func cartItemConvert(item cartItemEntity) models.CartItem {
	modelItem := models.CartItem{
		ID:       item.ID,
		CartID:   item.CartID,
		Product:  item.Product,
		Quantity: item.Quantity,
	}

	return modelItem
}
