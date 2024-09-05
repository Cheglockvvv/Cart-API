package service

import (
	"context"
	"fmt"
	"github.com/Cheglockvvv/Cart-API/internal/errs"
	"github.com/Cheglockvvv/Cart-API/internal/models"
)

type CartItemRepository interface {
	Create(context.Context, string, string, int) (string, error)
	Read(context.Context, string) (models.CartItem, error)
	Delete(context.Context, string, string) error
	ItemExists(context.Context, string) (bool, error)
}

type CartItem struct {
	cartRepository     CartRepository
	cartItemRepository CartItemRepository
}

func NewCartItem(cartRepository CartRepository, cartItemRepository CartItemRepository) *CartItem {
	return &CartItem{cartRepository: cartRepository, cartItemRepository: cartItemRepository}
}

func (c *CartItem) AddItemToCart(ctx context.Context, cartID, product string, quantity int) (models.CartItem, error) {
	ok, err := c.cartRepository.CartExists(ctx, cartID)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("c.cartRepository.CartExists: %w", err)
	}

	if !ok {
		return models.CartItem{}, errs.ErrCartNotFound // TODO: same
	}

	// TODO: add transaction blyaaat jackpot
	itemID, err := c.cartItemRepository.Create(ctx, cartID, product, quantity)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("c.cartRepository.Create: %w", err)
	}

	item, err := c.cartItemRepository.Read(ctx, itemID)
	if err != nil {
		return models.CartItem{}, fmt.Errorf("c.cartRepository.Read: %w", err)
	}

	return item, nil
}

func (c *CartItem) RemoveItemFromCart(ctx context.Context, cartID, itemID string) error {
	ok, err := c.cartRepository.CartExists(ctx, cartID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.CartExists: %w", err)
	}

	if !ok {
		return errs.ErrCartNotFound // TODO: same
	}

	ok, err = c.cartItemRepository.ItemExists(ctx, itemID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.ItemExists: %w", err)
	}

	if !ok {
		return errs.ErrItemNotFound
	}

	err = c.cartItemRepository.Delete(ctx, cartID, itemID)
	if err != nil {
		return fmt.Errorf("c.cartRepository.Delete: %w", err)
	}

	return nil
}
