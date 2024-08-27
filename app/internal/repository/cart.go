package repository

import (
	models "Cart-API/app/internal/models"
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type cartEntity struct {
	ID    string           `db:"id"`
	Items []cartItemEntity `db:"items"`
}

type Cart struct {
	DB *sqlx.DB
}

func InitCart(db *sqlx.DB) *Cart {
	cart := &Cart{DB: db}

	return cart
}

func (c *Cart) CreateCart(ctx context.Context) (string, error) {
	const query = `INSERT INTO cart VALUES (DEFAULT) RETURNING id`

	var id string
	err := c.DB.QueryRowx(query).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("c.DB.QueryRowx.Scan: %w", err)
	}

	return id, nil
}

func (c *Cart) GetCartByID(ctx context.Context, id string) (models.Cart, error) {
	const query = `SELECT ci.id, ci.cart_id, ci.product, ci.quantity 
								FROM cart_item ci 
								WHERE ci.cart_id = $1`

	rows, err := c.DB.Queryx(query, id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.DB.Queryx: %w", err)
	}
	defer rows.Close()

	cartEn := cartEntity{ID: id}

	items := make([]cartItemEntity, 0)

	for rows.Next() {
		row := cartItemEntity{}
		err = rows.StructScan(&row)

		if err != nil {
			return models.Cart{}, fmt.Errorf("rows.StructScan: %w", err)
		}
		items = append(items, row)
	}

	convertedItems := make([]models.CartItem, len(items))
	for i := range items {
		convertedItems[i] = cartItemConvert(items[i])
	}

	cart := models.Cart{ID: cartEn.ID, Items: convertedItems}
	return cart, nil
}

func (c *Cart) CartIsAvailable(ctx context.Context, id string) (bool, error) {

	const checkCart = `SELECT id FROM cart WHERE id = $1`
	result, err := c.DB.Exec(checkCart, id)
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
