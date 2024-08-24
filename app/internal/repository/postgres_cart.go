package repository

import (
	models2 "Cart-API/app/internal/models"
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

func (r *PostgresCart) GetCartByID(id string) (models2.Cart, error) {
	const query = `SELECT ci.id, ci.cart_id, ci.name, ci.quantity 
								FROM cart_item ci 
								WHERE ci.cart_id = $1`

	rows, err := r.DB.Queryx(query, id)
	if err != nil {
		return models2.Cart{}, fmt.Errorf("r.DB.Queryx: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return models2.Cart{}, nil
	}

	cart := models2.Cart{ID: id}

	//TODO: change struct
	items := make([]cartItem, 0)

	for rows.Next() {
		row := cartItem{}
		err = rows.StructScan(&row)

		if err != nil {
			return models2.Cart{}, fmt.Errorf("rows.StructScan: %w", err)
		}
		items = append(items, row)
	}

	convertedItems := make([]models2.CartItem, len(items))
	for i := range items {
		convertedItems[i] = modelConvert(items[i])
	}

	cart.Items = convertedItems
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

func modelConvert(item cartItem) models2.CartItem {
	modelItem := models2.CartItem{
		ID:       item.ID,
		CartID:   item.CartID,
		Name:     item.Name,
		Quantity: item.Quantity,
	}

	return modelItem
}
