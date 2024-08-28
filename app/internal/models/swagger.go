package models

type AddItemToCartRequest struct {
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}
