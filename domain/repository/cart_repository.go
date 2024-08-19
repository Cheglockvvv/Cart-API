package repository

import "Cart-API/domain/models"

type CartRepository interface {
	CreateCart(cart *models.Cart) error
	GetCartByID(id int64) (*models.Cart, error)
	AddItemToCart(id int64, item *models.CartItem) error
	RemoveItemFromCart(cartID int64, itemID int64) error
}
