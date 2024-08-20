package service

import "Cart-API/internal/models"

type CartRepository interface {
	Init(dataSourceName string) error
	CreateCart() (string, error)
	GetCartByID(id string) (*models.Cart, error)
	AddItemToCart(cartID, name string, quantity int) (string, error)
	RemoveItemFromCart(cartID, itemID string) error
}

type Cart struct {
	cartRepository CartRepository
}

func NewService(repository CartRepository) *Cart {
	return &Cart{cartRepository: repository}
}

func (service *Cart) CreateCart() (string, error) {
	id, err := service.cartRepository.CreateCart()
	if err != nil {
		return "", err
	}
	return id, nil
}

func (service *Cart) GetCartByID(id string) (*models.Cart, error) {
	cart, err := service.cartRepository.GetCartByID(id)

	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (service *Cart) AddItemToCart(cartID, name string, quantity int) (string, error) {
	itemID, err := service.cartRepository.AddItemToCart(cartID, name, quantity)

	if err != nil {
		return "", err
	}

	return itemID, nil
}

func (service *Cart) RemoveItemFromCart(cartID, itemID string) error {
	err := service.cartRepository.RemoveItemFromCart(cartID, itemID)
	return err
}
