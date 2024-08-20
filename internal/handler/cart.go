package handler

import (
	"Cart-API/internal/models"
	"net/http"
	"regexp"
)

type CartService interface {
	CreateCart() (string, error)
	GetCartByID(cartID string) (*models.Cart, error)
	AddItemToCart(cartID, name string, quantity int) (string, error)
	RemoveItemFromCart(cartID, itemID string) error
}

type Cart struct {
	cartService CartService
}

func NewHandler(cartService CartService) *Cart {
	return &Cart{cartService: cartService}
}

func (c *Cart) CreateCart(w http.ResponseWriter, r *http.Request) {}

func (c *Cart) AddItemToCart(w http.ResponseWriter, r *http.Request) {}

func (c *Cart) RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {}

func (c *Cart) GetCartByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cart, err := c.cartService.GetCartByID(id)
}

func (c *Cart) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		CartCreate = regexp.MustCompile(`^/cart/*$`)
		CartAdd    = regexp.MustCompile(`^/cart/[0-9]+/items/*$`)
		CartRemove = regexp.MustCompile(`^/cart/[0-9]+/items/[0-9]+/*$`)
		CartView   = regexp.MustCompile(`^/cart/[0-9]+/*$`)
	)

	switch {
	case r.Method == http.MethodPost && CartCreate.MatchString(r.URL.Path):
		c.CreateCart(w, r)
		return
	case r.Method == http.MethodPost && CartAdd.MatchString(r.URL.Path):
		c.AddItemToCart(w, r)
		return
	case r.Method == http.MethodDelete && CartRemove.MatchString(r.URL.Path):
		c.
	}
}
