package main

import (
	"fmt"
	"github.com/Cheglockvvv/Cart-API/config"
	_ "github.com/Cheglockvvv/Cart-API/docs"
	"github.com/Cheglockvvv/Cart-API/internal/db/connection"
	"github.com/Cheglockvvv/Cart-API/internal/db/migrations"
	"github.com/Cheglockvvv/Cart-API/internal/handler"
	"github.com/Cheglockvvv/Cart-API/internal/repository"
	"github.com/Cheglockvvv/Cart-API/internal/service"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
)

// @title Swagger Example API
// @version 1.0
// @description This is a documentation to Cart-API
// @termsOfService http://swagger.io/terms/

// @Contact.name API Support
// @Contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := connection.GetConnection(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	cartRepository := repository.NewCart(db)
	cartItemRepository := repository.NewCartItem(db)

	//err = migrations.Down(cartRepository.db)
	err = migrations.Up(cartRepository.DB)
	if err != nil {
		log.Fatal(err)
	}

	cartService := service.NewCart(cartRepository)
	cartItemService := service.NewCartItem(cartRepository, cartItemRepository)

	cartHandler := handler.NewCart(cartService, cartItemService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /cart", cartHandler.CreateCart)
	mux.HandleFunc("POST /cart/{id}/items", cartHandler.AddItemToCart)
	mux.HandleFunc("DELETE /cart/{cartID}/items/{itemID}", cartHandler.RemoveItemFromCart)
	mux.HandleFunc("GET /cart/{id}", cartHandler.GetCartByID)

	//mux.Handle("/swagger/*", http.StripPrefix("/swagger/*",
	//	http.FileServer(http.Dir("./docs.json")))) // TODO: remove comments

	mux.Handle("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //TODO: chi, port
	))

	//r := chi.NewRouter() TODO: remove comments
	//
	//r.Get("/swagger/*", httpSwagger.Handler(
	//	httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.API.Port), mux)

	if err != nil {
		log.Fatal(err)
	}
}
