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

type cartItemEntity struct {
	ID       string `json:"id"`
	CartID   string `json:"cart_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type cartEntity struct {
	ID    string           `json:"id"`
	Items []cartItemEntity `json:"items"`
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

	cart := cartEntity{ID: cartID, Items: []cartItemEntity{}}
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
	var parsedBody cartItemEntity
	err := json.NewDecoder(r.Body).Decode(&parsedBody)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	switch {
	case parsedBody.Quantity <= 0:
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	case parsedBody.Product == "":
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	item, err := c.cartService.AddItemToCart(cartID, parsedBody.Product, parsedBody.Quantity)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	convertedItem := cartItemConvert(item)
	err = json.NewEncoder(w).Encode(convertedItem)

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

func cartItemConvert(modelItem models.CartItem) cartItemEntity {
	item := cartItemEntity{
		ID:       modelItem.ID,
		CartID:   modelItem.CartID,
		Product:  modelItem.Product,
		Quantity: modelItem.Quantity,
	}

	return item
}

func cartConvert(modelCart models.Cart) cartEntity {
	cartItems := make([]cartItemEntity, len(modelCart.Items))
	for i, item := range modelCart.Items {
		cartItems[i] = cartItemConvert(item)
	}

	cart := cartEntity{
		ID:    modelCart.ID,
		Items: cartItems,
	}

	return cart
}
