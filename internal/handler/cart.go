package handler

import (
	"Cart-API/internal/models"
	"net/http"
)

type CartService interface {
	CreateCart() (string, error)
	GetCartByID(cartID string) (*models.Cart, error)
	AddItemToCart(cartID, name string, quantity int) error
	RemoveItemFromCart(cartID, itemID string) error
}

type Cart struct {
	cartService CartService
}

func (c *Cart) HelloWorld(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		return
	}
}
