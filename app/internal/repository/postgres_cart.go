package repository

import (
	models "Cart-API/app/internal/models"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type cartItem struct {
	ID       string `db:"id"`
	CartID   string `db:"cart_id"`
	Name     string `db:"name"`
	Quantity int    `db:"quantity"`
}

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
	ok, err := r.cartIsAvailable(id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("r.cartIsAvailable: %w", err)
	}

	if !ok {
		return models.Cart{}, fmt.Errorf("cart not found")
	}
	const query = `SELECT ci.id, ci.cart_id, ci.name, ci.quantity 
								FROM cart_item ci 
								WHERE ci.cart_id = $1`

	rows, err := r.DB.Queryx(query, id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("r.DB.Queryx: %w", err)
	}
	defer rows.Close()

	cart := models.Cart{ID: id}

	items := make([]cartItem, 0)

	for rows.Next() {
		row := cartItem{}
		err = rows.StructScan(&row)

		if err != nil {
			return models.Cart{}, fmt.Errorf("rows.StructScan: %w", err)
		}
		items = append(items, row)
	}

	convertedItems := make([]models.CartItem, len(items))
	for i := range items {
		convertedItems[i] = modelConvert(items[i])
	}

	cart.Items = convertedItems
	return cart, nil
}

func (r *PostgresCart) AddItemToCart(cartID, name string, quantity int) (models.CartItem, error) {

	ok, err := r.cartIsAvailable(cartID)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("r.cartIsAvailable: %w", err)
	}

	if !ok {
		return models.CartItem{}, fmt.Errorf("cart not found")
	}

	const query = `INSERT INTO cart_item (cart_id, name, quantity)
								VALUES ($1, $2, $3)
								ON CONFLICT (cart_id, name)
								DO UPDATE SET quantity = cart_item.quantity +
								    EXCLUDED.quantity
								RETURNING id`
	var itemID string
	err = r.DB.QueryRowx(query, cartID, name, quantity).Scan(&itemID)

	if err != nil {
		return models.CartItem{}, fmt.Errorf("r.DB.QueryRowx: %w", err)
	}

	item, err := r.GetItem(itemID)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("r.GetItem: %w", err)
	}

	return item, nil
}

func (r *PostgresCart) RemoveItemFromCart(cartID, itemID string) error {

	ok, err := r.cartIsAvailable(cartID)
	if err != nil {
		return fmt.Errorf("r.cartIsAvailable: %w", err)
	}

	if !ok {
		return fmt.Errorf("cart not found")
	}

	const checkItem = `SELECT id FROM cart_item WHERE id = $1`
	result, err := r.DB.Exec(checkItem, itemID)
	if err != nil {
		return fmt.Errorf("r.DB.Exec: %w", err)
	}
	count, err := result.RowsAffected()
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

func (r *PostgresCart) cartIsAvailable(id string) (bool, error) {

	const checkCart = `SELECT id FROM cart WHERE id = $1`
	result, err := r.DB.Exec(checkCart, id)
	if err != nil {
		return false, fmt.Errorf("r.DB.Exec: %w", err)
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

func (r *PostgresCart) GetItem(id string) (models.CartItem, error) {
	const query = `SELECT ci.id, ci.cart_id, ci.name, ci.quantity 
								FROM cart_item ci 
								WHERE ci.id = $1`
	result := r.DB.QueryRowx(query, id)
	item := cartItem{}
	err := result.StructScan(&item)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("result.StructScan: %w", err)
	}

	convertedItem := modelConvert(item)

	return convertedItem, nil
}

func modelConvert(item cartItem) models.CartItem {
	modelItem := models.CartItem{
		ID:       item.ID,
		CartID:   item.CartID,
		Name:     item.Name,
		Quantity: item.Quantity,
	}

	return modelItem
}
