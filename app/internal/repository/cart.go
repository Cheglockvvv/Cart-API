package repository

import (
	models "Cart-API/app/internal/models"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type cartItemEntity struct {
	id       string `db:"id"`
	cartID   string `db:"cart_id"`
	product  string `db:"product"`
	quantity int    `db:"quantity"`
}

type cartEntity struct {
	id    string           `db:"id"`
	items []cartItemEntity `db:"items"`
}

type Cart struct {
	DB *sqlx.DB
}

func (c *Cart) Init(connectionString string) error {
	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		return fmt.Errorf("sqlx.Connect: %w", err)
	}

	c.DB = db
	return nil
}

func (c *Cart) CreateCart() (string, error) {
	const query = `INSERT INTO cart VALUES (DEFAULT) RETURNING id`

	var id string
	err := c.DB.QueryRowx(query).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("c.DB.QueryRowx.Scan: %w", err)
	}

	return id, nil
}

func (c *Cart) GetCartByID(id string) (models.Cart, error) {
	const query = `SELECT ci.id, ci.cart_id, ci.name, ci.quantity 
								FROM cart_item ci 
								WHERE ci.cart_id = $1`

	rows, err := c.DB.Queryx(query, id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.DB.Queryx: %w", err)
	}
	defer rows.Close()

	cartEn := cartEntity{id: id}

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

	cart := models.Cart{ID: cartEn.id, Items: convertedItems}
	return cart, nil
}

func (c *Cart) AddItemToCart(cartID, name string, quantity int) (string, error) {

	const query = `INSERT INTO cart_item (cart_id, name, quantity)
								VALUES ($1, $2, $3)
								ON CONFLICT (cart_id, name)
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

func (c *Cart) RemoveItemFromCart(cartID, itemID string) error {
	const query = `DELETE FROM cart_item WHERE cart_id = $1 AND id = $2`

	_, err := c.DB.Exec(query, cartID, itemID)
	if err != nil {
		return fmt.Errorf("c.DB.Exec: %w", err)
	}

	return nil
}

func (c *Cart) CartIsAvailable(id string) (bool, error) {

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

func (c *Cart) ItemIsAvailable(id string) (bool, error) {
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

func (c *Cart) GetItemByID(id string) (models.CartItem, error) {
	const query = `SELECT ci.id, ci.cart_id, ci.name, ci.quantity 
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

func cartItemConvert(item cartItemEntity) models.CartItem {
	modelItem := models.CartItem{
		ID:       item.id,
		CartID:   item.cartID,
		Product:  item.product,
		Quantity: item.quantity,
	}

	return modelItem
}
