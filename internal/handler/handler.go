package handler

import "Cart-API/internal/models"

type CartService interface {
	CreateCart() (string, error)
	GetCartByID(cartID string) (*models.Cart, error)
	AddItemToCart(cartID, name string, quantity int) error
	RemoveItemFromCart(cartID, itemID string) error
}

type Cart struct {
	cartService CartService
}
