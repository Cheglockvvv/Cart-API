package repository

import (
	"context"
	"fmt"
	"github.com/Cheglockvvv/Cart-API/internal/models"
	"github.com/jmoiron/sqlx"
)

type cartEntity struct { //TODO: no entity - dto
	ID    string           `db:"id"`
	Items []cartItemEntity `db:"items"`
}

type Cart struct {
	DB *sqlx.DB // TODO: switch to lowercase
}

func NewCart(db *sqlx.DB) *Cart { // TODO: fsffsklfj
	cart := &Cart{DB: db}

	return cart
}

func (c *Cart) CreateCart(ctx context.Context) (string, error) {
	const query = `INSERT INTO cart VALUES (DEFAULT) RETURNING id`

	var id string
	err := c.DB.QueryRowxContext(ctx, query).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("c.DB.QueryRowx.Scan: %w", err)
	} // TODO:

	return id, nil
}

func (c *Cart) GetCartByID(ctx context.Context, id string) (models.Cart, error) {
	const query = `SELECT ci.id, ci.cart_id, ci.product, ci.quantity 
								FROM cart_item ci 
								WHERE ci.cart_id = $1`

	rows, err := c.DB.QueryxContext(ctx, query, id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.DB.QueryxContext: %w", err)
	}
	defer rows.Close()

	cartEn := cartEntity{ID: id}

	items := make([]cartItemEntity, 0)

	for rows.Next() {
		row := cartItemEntity{}
		err = rows.StructScan(&row)

		if err != nil {
			return models.Cart{}, fmt.Errorf("rows.StructScan: %w", err)
		} // TODO:
		items = append(items, row)
	}

	convertedItems := make([]models.CartItem, len(items)) // TODO: unnecessary allocations
	for i := range items {
		convertedItems[i] = cartItemConvert(items[i])
	}

	cart := models.Cart{ID: cartEn.ID, Items: convertedItems}

	return cart, nil
}

func (c *Cart) CartIsAvailable(ctx context.Context, id string) (bool, error) { // TODO: rename to exists returns error

	const query = `SELECT EXISTS(SELECT 1 FROM cart_item WHERE cart_id = $1)` // TODO: eбаааать
	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return false, fmt.Errorf("c.DB.Exec: %w", err)
	} // TODO:
	count, err := result.RowsAffected() //TODO: remove
	if err != nil {
		return false, fmt.Errorf("result.RowsAffected: %w", err)
	} // TODO:

	if count != 1 {
		return false, nil
	}

	return true, nil
}
