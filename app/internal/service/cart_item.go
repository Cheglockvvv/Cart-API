package service

import (
	"Cart-API/app/internal/models"
	"context"
	"fmt"
)

type CartItemRepository interface {
	AddItemToCart(context.Context, string, string, int) (string, error)
	GetItemByID(context.Context, string) (models.CartItem, error)
	RemoveItemFromCart(context.Context, string, string) error
	ItemIsAvailable(context.Context, string) (bool, error)
}

type CartItem struct {
	cartRepository     CartRepository
	cartItemRepository CartItemRepository
}

func NewCartItem(cartRepository CartRepository, cartItemRepository CartItemRepository) *CartItem {
	return &CartItem{cartRepository: cartRepository, cartItemRepository: cartItemRepository}
}

func (c *CartItem) AddItemToCart(ctx context.Context, cartID, product string, quantity int) (models.CartItem, error) {
	ok, err := c.cartRepository.CartIsAvailable(ctx, cartID)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("c.cartRepository.CartIsAvailable: %w", err)
	}

	if !ok {
		return models.CartItem{}, fmt.Errorf("cart is not available")
	}

	itemID, err := c.cartItemRepository.AddItemToCart(ctx, cartID, product, quantity)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("c.cartRepository.AddItemToCart: %w", err)
	}

	item, err := c.cartItemRepository.GetItemByID(ctx, itemID)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("c.cartRepository.GetItemByID: %w", err)
	}

	return item, nil
}

func (c *CartItem) RemoveItemFromCart(ctx context.Context, cartID, itemID string) error {
	ok, err := c.cartRepository.CartIsAvailable(ctx, cartID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.CartIsAvailable: %w", err)
	}

	if !ok {
		return fmt.Errorf("cart is not available")
	}

	ok, err = c.cartItemRepository.ItemIsAvailable(ctx, itemID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.ItemIsAvailable: %w", err)
	}

	if !ok {
		return fmt.Errorf("cart item is not available")
	}

	err = c.cartItemRepository.RemoveItemFromCart(ctx, cartID, itemID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.RemoveItemFromCart: %w", err)
	}

	return nil
}
