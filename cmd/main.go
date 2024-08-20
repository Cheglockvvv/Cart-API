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
	err := cartRepository.Init("fasd")
	if err != nil {
		log.Fatal(err)
	}

	cartService := service.NewCart(&cartRepository)

	cartHandler := handler.NewHandler(cartService)

	mux := http.NewServeMux()
	mux.Handle("POST /cart/{id}", cartHandler)
	mux.HandleFunc("lksgjaslg", cartHandler.CreateCart)

	http.ListenAndServe(":8080", mux)
}
