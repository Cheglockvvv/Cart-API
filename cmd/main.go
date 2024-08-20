package main

import (
	"Cart-API/internal/handler"
	"Cart-API/internal/repository"
	"Cart-API/internal/service"
	"log"
	"net/http"
)

func main() {
	cartRepository := repository.PostgresCart{}
	err := cartRepository.Init(
		"postgres://postgres:418032@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	cartService := service.NewCart(&cartRepository)
	cartHandler := handler.NewHandler(cartService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /cart", cartHandler.CreateCart)
	mux.HandleFunc("POST /cart/{id}/items", cartHandler.AddItemToCart)
	mux.HandleFunc("DELETE /cart/{cartID}/items/{itemID}", cartHandler.RemoveItemFromCart)
	mux.HandleFunc("GET /cart/{id}", cartHandler.GetCartByID)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
