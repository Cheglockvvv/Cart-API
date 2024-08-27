package service

import (
	"Cart-API/app/internal/models"
	"fmt"
)

type CartRepository interface {
	CreateCart() (string, error)
	GetCartByID(id string) (models.Cart, error)
	CartIsAvailable(id string) (bool, error)
}

type Cart struct {
	cartRepository CartRepository
}

func NewCart(repository CartRepository) *Cart {
	return &Cart{cartRepository: repository}
}

func (c *Cart) CreateCart() (string, error) {
	id, err := c.cartRepository.CreateCart()
	if err != nil {
		return "", fmt.Errorf("c.cartRepository.CreateCart: %w", err)
	}
	return id, nil
}

func (c *Cart) GetCartByID(id string) (models.Cart, error) {
	ok, err := c.cartRepository.CartIsAvailable(id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.cartRepository.CartIsAvailable: %w", err)
	}

	if !ok {
		return models.Cart{}, fmt.Errorf("cart not available")
	}

	cart, err := c.cartRepository.GetCartByID(id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.cartRepository.GetCartByID: %w", err)
	}

	return cart, nil
}
