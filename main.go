package main

import (
	"log"
	"net/http"

	"favouritesApp/internal/handlers"
	services "favouritesApp/internal/services/favourite"
	"favouritesApp/internal/storage"

	"github.com/gorilla/mux"
)

func main() {
	store := storage.NewMemoryStore()
	service := services.NewFavouriteService(store)
	handler := handlers.NewFavouriteHandler(service)

	r := mux.NewRouter()
	r.HandleFunc("/users/{userId}/favourites", handler.GetFavourites).Methods("GET")
	r.HandleFunc("/users/{userId}/favourites", handler.AddFavourite).Methods("POST")
	r.HandleFunc("/users/{userId}/favourites/{assetId}", handler.EditFavourite).Methods("PUT")
	r.HandleFunc("/users/{userId}/favourites/{assetId}", handler.RemoveFavourite).Methods("DELETE")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
