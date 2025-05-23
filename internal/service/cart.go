package service

import (
	"context"
	"fmt"
	"github.com/Cheglockvvv/Cart-API/internal/errs"
	"github.com/Cheglockvvv/Cart-API/internal/models"
)

type CartRepository interface {
	Create(context.Context) (string, error)
	Read(context.Context, string) (models.Cart, error)
	CartExists(context.Context, string) (bool, error)
}

type Cart struct {
	cartRepository CartRepository
}

func NewCart(repository CartRepository) *Cart {
	return &Cart{cartRepository: repository}
}

func (c *Cart) CreateCart(ctx context.Context) (string, error) {
	id, err := c.cartRepository.Create(ctx)
	if err != nil {
		return "", fmt.Errorf("c.cartRepository.CreateCart: %w", err)
	}
	return id, nil
}

func (c *Cart) GetCartByID(ctx context.Context, id string) (models.Cart, error) {
	ok, err := c.cartRepository.CartExists(ctx, id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.cartRepository.CartExists: %w", err)
	}

	if !ok {
		return models.Cart{}, errs.ErrCartNotFound // TODO: cart is available return this error
	}

	cart, err := c.cartRepository.Read(ctx, id)
	if err != nil {
		return models.Cart{}, fmt.Errorf("c.cartRepository.GetCartByID: %w", err)
	}

	return cart, nil
}
