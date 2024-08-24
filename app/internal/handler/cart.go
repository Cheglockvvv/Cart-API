package handler

import (
	models "Cart-API/app/internal/models"
	"encoding/json"
	"net/http"
	"regexp"
)

type CartService interface {
	CreateCart() (string, error)
	GetCartByID(cartID string) (models.Cart, error)
	AddItemToCart(cartID, name string, quantity int) (models.CartItem, error)
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
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	cart := models.Cart{ID: cartID, Items: []models.CartItem{}}
	err = json.NewEncoder(w).Encode(cart)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (c *Cart) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	cartAdd := regexp.MustCompile(`^/cart/[0-9]+/items/*$`)
	if !cartAdd.MatchString(r.URL.Path) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	cartID := r.PathValue("id")
	var parsedBody models.CartItem
	err := json.NewDecoder(r.Body).Decode(&parsedBody)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	switch {
	case parsedBody.Quantity <= 0:
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	case parsedBody.Name == "":
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	item, err := c.cartService.AddItemToCart(cartID, parsedBody.Name, parsedBody.Quantity)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(item)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (c *Cart) RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	cartRemove := regexp.MustCompile(`^/cart/[0-9]+/items/[0-9]+/*$`)
	if !cartRemove.MatchString(r.URL.Path) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	cartID := r.PathValue("cartID")
	itemID := r.PathValue("itemID")

	err := c.cartService.RemoveItemFromCart(cartID, itemID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte("{}"))

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (c *Cart) GetCartByID(w http.ResponseWriter, r *http.Request) {
	cartView := regexp.MustCompile(`^/cart/[0-9]+/*$`)
	if !cartView.MatchString(r.URL.Path) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	cart, err := c.cartService.GetCartByID(id)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if cart.ID == "" {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cart)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
