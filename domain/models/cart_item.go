package models

type CartItem struct {
	ID       int64  `json:"id"`
	CartID   int64  `json:"cart_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}
