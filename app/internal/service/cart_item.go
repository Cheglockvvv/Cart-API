package service

import "Cart-API/app/internal/models"

type CartItemRepository interface {
	AddItemToCart(cartID, name string, quantity int) (string, error)
	GetItemByID(id string) (models.CartItem, error)
	RemoveItemFromCart(cartID, itemID string) error
	ItemIsAvailable(id string) (bool, error)
}

type CartItem struct {
}
