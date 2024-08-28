package handler

import (
	"context"
	"encoding/json"
	models "github.com/Cheglockvvv/Cart-API/app/internal/models"
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

func NewHandler(cartService CartService, cartItemService CartItemService) *Cart {
	return &Cart{cartService: cartService, cartItemService: cartItemService}
}

// CreateCart godoc
// @Summary Create a new cart
// @Description Creates a cart and returns it
// @Tags Cart
// @Produce json
// @Success 200 {object} cartEntity
// @Failure 500 "Internal Server Error"
// @Router /cart [post]
func (c *Cart) CreateCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cartID, err := c.cartService.CreateCart(ctx)
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

// AddItemToCart godoc
// @Summary Add item to cart
// @Description adds an item to a specified cart with provided details and returns it
// @Tags CartItem
// @Accept json
// @Produce json
// @Param id path string true "cart id"
// @Param item body models.AddItemToCartRequest true "Item to add to cart"
// @Success 200 {object} cartItemEntity
// @Failure 400 "Bad Request"
// @Failure 422 "Unprocessable Entity"
// @Failure 500 "Internal Server Error"
// @Router /cart/{id}/items [post]
func (c *Cart) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	item, err := c.cartItemService.AddItemToCart(ctx, cartID, parsedBody.Product, parsedBody.Quantity)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	convertedItem := cartItemConvert(item)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(convertedItem)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// RemoveItemFromCart godoc
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

	cartRemove := regexp.MustCompile(`^/cart/[0-9]+/items/[0-9]+/*$`)
	if !cartRemove.MatchString(r.URL.Path) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	cartID := r.PathValue("cartID")
	itemID := r.PathValue("itemID")

	err := c.cartItemService.RemoveItemFromCart(ctx, cartID, itemID)
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

// GetCartByID godoc
// @Summary Get a cart by ID
// @Description With specified CartID returns a cart
// @Tags Cart
// @Produce json
// @Param id path string true "CartID"
// @Success 200 {object} cartEntity
// @Failure 400 "Bad Request"
// @Failure 422 "Unprocessable Entity"
// @Failure 500 "Internal Server Error"
// @Router /cart/{id} [get]
func (c *Cart) GetCartByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cartView := regexp.MustCompile(`^/cart/[0-9]+/*$`)
	if !cartView.MatchString(r.URL.Path) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
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

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(convertedCart)

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
