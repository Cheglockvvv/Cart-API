package service

import (
	"Cart-API/internal/models"
	"fmt"
)

type CartRepository interface {
	CreateCart() (string, error)
	GetCartByID(id string) (models.Cart, error)
	AddItemToCart(cartID, name string, quantity int) (string, error)
	RemoveItemFromCart(cartID, itemID string) error
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
	cart, err := c.cartRepository.GetCartByID(id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.cartRepository.GetCartByID: %w", err)
	}

	return cart, nil
}

func (c *Cart) AddItemToCart(cartID, name string, quantity int) (string, error) {
	itemID, err := c.cartRepository.AddItemToCart(cartID, name, quantity)
	if err != nil {
		return "", fmt.Errorf("c.cartRepository.AddItemToCart: %w", err)
	}

	return itemID, nil
}

func (c *Cart) RemoveItemFromCart(cartID, itemID string) error {
	err := c.cartRepository.RemoveItemFromCart(cartID, itemID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.RemoveItemFromCart: %w", err)
	}

	return nil
}
