package router

import (
	"favouritesApp/internal/handlers"
	services "favouritesApp/internal/services/favourite"
	"favouritesApp/internal/storage"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	store := storage.NewMemoryStore()
	service := services.NewFavouriteService(store)
	handler := handlers.NewFavouriteHandler(service)

	r := mux.NewRouter()
	r.HandleFunc("/users/{userId}/favourites", handler.GetFavourites).Methods("GET")
	r.HandleFunc("/users/{userId}/favourites", handler.AddFavourite).Methods("POST")
	r.HandleFunc("/users/{userId}/favourites/{assetId}", handler.EditFavourite).Methods("PUT")
	r.HandleFunc("/users/{userId}/favourites/{assetId}", handler.RemoveFavourite).Methods("DELETE")

	// health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r
}
