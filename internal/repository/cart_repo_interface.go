package repository

import (
	"Cart-API/internal/models"
	"github.com/jmoiron/sqlx"
)

type CartRepoInterface interface {
	Init(dataSourceName string) (*sqlx.DB, error)
	CreateCart() (string, error)
	GetCartByID(cartID string) (*models.Cart, error)
	AddItemToCart(cartID string, item *models.CartItem) (string, error)
	RemoveItemFromCart(cartID, itemID string) error
}
