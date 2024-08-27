package service

import (
	"Cart-API/app/internal/models"
	"context"
	"fmt"
)

type CartRepository interface {
	CreateCart(context.Context) (string, error)
	GetCartByID(context.Context, string) (models.Cart, error)
	CartIsAvailable(context.Context, string) (bool, error)
}

type Cart struct {
	cartRepository CartRepository
}

func NewCart(repository CartRepository) *Cart {
	return &Cart{cartRepository: repository}
}

func (c *Cart) CreateCart(ctx context.Context) (string, error) {
	id, err := c.cartRepository.CreateCart(ctx)
	if err != nil {
		return "", fmt.Errorf("c.cartRepository.CreateCart: %w", err)
	}
	return id, nil
}

func (c *Cart) GetCartByID(ctx context.Context, id string) (models.Cart, error) {
	ok, err := c.cartRepository.CartIsAvailable(ctx, id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.cartRepository.CartIsAvailable: %w", err)
	}

	if !ok {
		return models.Cart{}, fmt.Errorf("cart not available")
	}

	cart, err := c.cartRepository.GetCartByID(ctx, id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.cartRepository.GetCartByID: %w", err)
	}

	return cart, nil
}
