package handler

import (
	"context"
	"encoding/json"
	"github.com/Cheglockvvv/Cart-API/internal/models"
	"github.com/go-chi/chi"
	"net/http"
	"regexp"
)

type CartService interface {
	CreateCart(context.Context) (string, error)
	GetCartByID(context.Context, string) (models.Cart, error)
}

type CartItemService interface {
	AddItemToCart(context.Context, string, string, int) (models.CartItem, error)
	RemoveItemFromCart(context.Context, string, string) error
}

type Cart struct {
	cartService     CartService
	cartItemService CartItemService
}

type cartItemDTO struct {
	ID       string `json:"id"`
	CartID   string `json:"cart_id"`
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

type cartEntityDTO struct {
	ID    string        `json:"id"`
	Items []cartItemDTO `json:"items"`
}

func NewCart(cartService CartService, cartItemService CartItemService) *Cart {
	return &Cart{cartService: cartService, cartItemService: cartItemService}
}

var (
	cartItemAdd    = regexp.MustCompile(`^/cart/[0-9]+/items/*$`)
	cartItemRemove = regexp.MustCompile(`^/cart/[0-9]+/items/[0-9]+/*$`)
	cartView       = regexp.MustCompile(`^/cart/[0-9]+/*$`)
)

// CreateCart
// @Summary Create a new cart
// @Description Creates a cart and returns it
// @Tags Cart
// @Produce json
// @Success 200 {object} cartEntityDTO
// @Failure 500 "Internal Server Error"
// @Router /cart [post]
func (c *Cart) CreateCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cartID, err := c.cartService.CreateCart(ctx)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	cart := cartEntityDTO{ID: cartID, Items: []cartItemDTO{}}
	err = json.NewEncoder(w).Encode(cart)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// AddItemToCart
// @Description adds an item to a specified cart with provided details and returns it
// @Tags CartItem
// @Accept json
// @Produce json
// @Param id path string true "cart id"
// @Param item body handler.AddItemToCart.request true "Item to add to cart"
// @Success 200 {object} cartItemDTO
// @Failure 400 "Bad Request"
// @Failure 422 "Unprocessable Entity"
// @Failure 500 "Internal Server Error"
// @Router /cart/{id}/items [post]
func (c *Cart) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Product  string `json:"product"`
		Quantity int    `json:"quantity"`
	}

	ctx := r.Context()

	if !cartItemAdd.MatchString(r.URL.Path) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	cartID := chi.URLParam(r, "id")
	var parsedBody cartItemDTO
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

	item, err := c.cartItemService.AddItemToCart(ctx, cartID, parsedBody.Product,
		parsedBody.Quantity)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	convertedItem := cartItemConvert(item)
	w.Header().Set("Content-Type", "application/json") // TODO: move to middleware
	err = json.NewEncoder(w).Encode(convertedItem)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// RemoveItemFromCart
// @Summary Remove item from cart
// @Description removes a specified item from a specified cart
// @Tags CartItem
// @Produce json
// @Param cart_id path string true "CartID"
// @Param item_id path string true "ItemID"
// @Success 200 "{}"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /cart/{cart_id}/items/{item_id} [delete]
func (c *Cart) RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if !cartItemRemove.MatchString(r.URL.Path) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	cartID := chi.URLParam(r, "cart_id")
	itemID := chi.URLParam(r, "item_id")

	err := c.cartItemService.RemoveItemFromCart(ctx, cartID, itemID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // TODO: same
	_, err = w.Write([]byte("{}"))

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

// GetCartByID
// @Summary Get a cart by ID
// @Description With specified CartID returns a cart
// @Tags Cart
// @Produce json
// @Param id path string true "CartID"
// @Success 200 {object} cartEntityDTO
// @Failure 400 "Bad Request"
// @Failure 422 "Unprocessable Entity"
// @Failure 500 "Internal Server Error"
// @Router /cart/{id} [get]
func (c *Cart) GetCartByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if !cartView.MatchString(r.URL.Path) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	id := chi.URLParam(r, "id")
	cart, err := c.cartService.GetCartByID(ctx, id)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if cart.ID == "" {
		http.Error(w, "", http.StatusUnprocessableEntity)
		return
	}

	convertedCart := cartConvert(cart)

	w.Header().Set("Content-Type", "application/json") // TODO: same
	err = json.NewEncoder(w).Encode(convertedCart)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func cartItemConvert(modelItem models.CartItem) cartItemDTO {
	item := cartItemDTO{
		ID:       modelItem.ID,
		CartID:   modelItem.CartID,
		Product:  modelItem.Product,
		Quantity: modelItem.Quantity,
	}

	return item
}

func cartConvert(modelCart models.Cart) cartEntityDTO {
	cartItems := make([]cartItemDTO, 0, len(modelCart.Items))
	for _, item := range modelCart.Items {
		cartItems = append(cartItems, cartItemConvert(item))
	}

	cart := cartEntityDTO{
		ID:    modelCart.ID,
		Items: cartItems,
	}

	return cart
}
