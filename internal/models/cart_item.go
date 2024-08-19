package models

type CartItem struct {
	ID       string `json:"id"`
	CartID   string `json:"cart_id"`
	Product  string `json:"name"`
	Quantity int    `json:"quantity"`
}
