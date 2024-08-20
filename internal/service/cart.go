package service

import "Cart-API/internal/models"

type CartRepository interface {
	Init(dataSourceName string) error
	CreateCart() (string, error)
	GetCartByID(id string) (*models.Cart, error)
	AddItemToCart(cartID, name string, quantity int) (string, error)
	RemoveItemFromCart(cartID, itemID string) error
}

type Cart struct {
	cartRepository CartRepository
}
