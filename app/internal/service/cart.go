package service

import (
	"Cart-API/app/internal/models"
	"fmt"
)

type CartRepository interface {
	CreateCart() (string, error)
	GetCartByID(id string) (models.Cart, error)
	AddItemToCart(cartID, name string, quantity int) (string, error)
	RemoveItemFromCart(cartID, itemID string) error
	CartIsAvailable(id string) (bool, error)
	GetItemByID(id string) (models.CartItem, error)
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

func (c *Cart) AddItemToCart(cartID, name string, quantity int) (models.CartItem, error) {
	ok, err := c.cartRepository.CartIsAvailable(cartID)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("c.cartRepository.CartIsAvailable: %w", err)
	}

	if !ok {
		return models.CartItem{}, fmt.Errorf("cart is not available")
	}

	item, err := c.cartRepository.AddItemToCart(cartID, name, quantity)
	if err != nil {
		return item, fmt.Errorf("c.cartRepository.AddItemToCart: %w", err)
	}

	return item, nil
}

func (c *Cart) RemoveItemFromCart(cartID, itemID string) error {
	ok, err := c.cartRepository.CartIsAvailable(cartID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.CartIsAvailable: %w", err)
	}

	if !ok {
		return fmt.Errorf("cart is not available")
	}

	err = c.cartRepository.RemoveItemFromCart(cartID, itemID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.RemoveItemFromCart: %w", err)
	}

	return nil
}
