package main

import (
	"Cart-API/app/config"
	"Cart-API/app/internal/db/connection"
	"Cart-API/app/internal/db/migrations"
	"Cart-API/app/internal/handler"
	"Cart-API/app/internal/repository"
	"Cart-API/app/internal/service"
	_ "Cart-API/docs"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
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

	mux.Handle("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.API.Port), mux)
	if err != nil {
		log.Fatal(err)
	}
}
