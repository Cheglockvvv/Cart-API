package handler

import (
	"Cart-API/internal/models"
	"encoding/json"
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

func (c *Cart) CreateCart(w http.ResponseWriter, r *http.Request) {
	cartID, err := c.cartService.CreateCart()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	cart := models.Cart{ID: cartID, Items: []models.CartItem{}}
	err = json.NewEncoder(w).Encode(cart)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Cart) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("id")
	var parsedBody models.CartItem
	err := json.NewDecoder(r.Body).Decode(&parsedBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	itemID, err := c.cartService.AddItemToCart(cartID, parsedBody.Name, parsedBody.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	parsedBody.CartID = cartID
	parsedBody.ID = itemID

	err = json.NewEncoder(w).Encode(parsedBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (c *Cart) RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cartID")
	itemID := r.PathValue("itemID")

	err := c.cartService.RemoveItemFromCart(cartID, itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte("{}"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (c *Cart) GetCartByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cart, err := c.cartService.GetCartByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cart)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (c *Cart) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

var (
	CartCreate = regexp.MustCompile(`^/cart/*$`)
	CartAdd    = regexp.MustCompile(`^/cart/[0-9]+/items/*$`)
	CartRemove = regexp.MustCompile(`^/cart/[0-9]+/items/[0-9]+/*$`)
	CartView   = regexp.MustCompile(`^/cart/[0-9]+/*$`)
)
