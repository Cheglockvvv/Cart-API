package repository

import (
	"context"
	"fmt"
	"github.com/Cheglockvvv/Cart-API/internal/models"
	"github.com/jmoiron/sqlx"
)

type cartDTO struct {
	id    string           `db:"id"`
	items []cartItemEntity `db:"items"`
}

type Cart struct {
	DB *sqlx.DB // TODO: switch to lowercase and make getter btw
}

func NewCart(db *sqlx.DB) *Cart {
	cart := &Cart{DB: db}

	return cart
}

func (c *Cart) Create(ctx context.Context) (string, error) {
	const query = `INSERT INTO cart VALUES (DEFAULT) RETURNING id`

	var id string
	err := c.DB.QueryRowxContext(ctx, query).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("c.DB.QueryRowxContext.Scan: %w", err)
	}

	return id, nil
}

func (c *Cart) Read(ctx context.Context, id string) (models.Cart, error) { // TODO: rename to read
	const query = `SELECT ci.id, ci.cart_id, ci.product, ci.quantity 
								FROM cart_item ci 
								WHERE ci.cart_id = $1`

	rows, err := c.DB.QueryxContext(ctx, query, id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.DB.QueryxContext: %w", err)
	}
	defer rows.Close()

	cartEn := cartDTO{id: id}

	items := make([]cartItemEntity, 0)

	for rows.Next() {
		row := cartItemEntity{}
		err = rows.StructScan(&row)

		if err != nil {
			return models.Cart{}, fmt.Errorf("rows.StructScan: %w", err)
		}
		items = append(items, row)
	}

	convertedItems := make([]models.CartItem, 0, len(items))
	for i := range items {
		convertedItems = append(convertedItems, cartItemConvert(items[i]))
	}

	cart := models.Cart{ID: cartEn.id, Items: convertedItems}

	return cart, nil
}

func (c *Cart) CartExists(ctx context.Context, id string) (bool, error) {

	const query = `SELECT EXISTS(SELECT 1 FROM cart WHERE id = $1)`
	var exists bool
	err := c.DB.QueryRowxContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("c.DB.QueryRowxContext.Scan: %w", err)
	}

	return exists, nil
}
