package service

import (
	"Cart-API/app/internal/models"
	"fmt"
)

type CartItemRepository interface {
	AddItemToCart(cartID, name string, quantity int) (string, error)
	GetItemByID(id string) (models.CartItem, error)
	RemoveItemFromCart(cartID, itemID string) error
	ItemIsAvailable(id string) (bool, error)
}

type CartItem struct {
	cartItemRepository CartItemRepository
}

func (c *CartItem) AddItemToCart(cartID, name string, quantity int) (models.CartItem, error) {
	ok, err := c.cartRepository.CartIsAvailable(cartID)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("c.cartRepository.CartIsAvailable: %w", err)
	}

	if !ok {
		return models.CartItem{}, fmt.Errorf("cart is not available")
	}

	itemID, err := c.cartRepository.AddItemToCart(cartID, name, quantity)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("c.cartRepository.AddItemToCart: %w", err)
	}

	item, err := c.cartRepository.GetItemByID(itemID)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("c.cartRepository.GetItemByID: %w", err)
	}

	return item, nil
}

func (c *CartItem) RemoveItemFromCart(cartID, itemID string) error {
	ok, err := c.cartRepository.CartIsAvailable(cartID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.CartIsAvailable: %w", err)
	}

	if !ok {
		return fmt.Errorf("cart is not available")
	}

	ok, err = c.cartRepository.ItemIsAvailable(itemID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.ItemIsAvailable: %w", err)
	}

	if !ok {
		return fmt.Errorf("cart item is not available")
	}

	err = c.cartRepository.RemoveItemFromCart(cartID, itemID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.RemoveItemFromCart: %w", err)
	}

	return nil
}
