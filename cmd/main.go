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
	"github.com/go-chi/chi"
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

	err = migrations.Up(cartRepository.DB, cfg.API.Migrations)
	if err != nil {
		log.Println(err)
	}

	cartService := service.NewCart(cartRepository)
	cartItemService := service.NewCartItem(cartRepository, cartItemRepository)

	cartHandler := handler.NewCart(cartService, cartItemService)

	router := chi.NewRouter()
	router.Post("/cart", cartHandler.CreateCart)
	router.Get("/cart/{id}", cartHandler.GetCartByID)
	router.Post("/cart/{id}/items", cartHandler.AddItemToCart)
	router.Delete("/cart/{cart_id}/items/{item_id}", cartHandler.RemoveItemFromCart)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", cfg.API.Port))))

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.API.Port), router)
	if err != nil {
		log.Fatal(err)
	}
}
