package main

import (
	"Cart-API/app/config"
	"Cart-API/app/internal/db/connection"
	"Cart-API/app/internal/db/migrations"
	"Cart-API/app/internal/handler"
	"Cart-API/app/internal/repository"
	"Cart-API/app/internal/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	const DSN = "postgres://%s:%s@%s:%s/%s?sslmode=%s"
	filledDsn := fmt.Sprintf(DSN, cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port,
		cfg.DB.DBName, cfg.DB.SSLMode)

	DB, err := connection.GetConnection(filledDsn)
	if err != nil {
		log.Fatal(err)
	}

	cartRepository := repository.InitCart(DB)
	cartItemRepository := repository.InitCartItem(DB)

	//err = migrations.Down(cartRepository.DB)
	err = migrations.Up(cartRepository.DB)
	if err != nil {
		log.Fatal(err)
	}

	cartService := service.NewCart(cartRepository)
	cartItemService := service.NewCartItem(cartRepository, cartItemRepository)

	cartHandler := handler.NewHandler(cartService, cartItemService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /cart", cartHandler.CreateCart)
	mux.HandleFunc("POST /cart/{id}/items", cartHandler.AddItemToCart)
	mux.HandleFunc("DELETE /cart/{cartID}/items/{itemID}", cartHandler.RemoveItemFromCart)
	mux.HandleFunc("GET /cart/{id}", cartHandler.GetCartByID)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.API.Port), mux)
	if err != nil {
		log.Fatal(err)
	}
}
