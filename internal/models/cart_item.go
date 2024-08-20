package models

type CartItem struct {
	ID       string
	CartID   string `db:"cart_id"`
	Name     string
	Quantity int
}
