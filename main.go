package main

import (
	"assetsApp/internal/handlers"
	assetServices "assetsApp/internal/services/asset"
	favouriteServices "assetsApp/internal/services/favourite"
	"assetsApp/internal/storage"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting the application...")
	store := storage.NewMemoryStore()

	// Asset related initializations
	assetService := assetServices.NewAssetService(store)
	assetHandler := handlers.NewAssetHandler(assetService)

	// Favourite related initializations
	favouriteService := favouriteServices.NewFavouriteService(store)
	favouriteHandler := handlers.NewFavouriteHandler(favouriteService)

	r := mux.NewRouter()

	// Asset routes
	r.HandleFunc("/users/{userId}/assets", assetHandler.GetAssets).Methods("GET")
	r.HandleFunc("/users/{userId}/assets", assetHandler.AddAsset).Methods("POST")
	r.HandleFunc("/users/{userId}/assets/{assetId}", assetHandler.EditAsset).Methods("PUT")
	r.HandleFunc("/users/{userId}/assets/{assetId}", assetHandler.RemoveAsset).Methods("DELETE")

	// Favourites routes
	r.HandleFunc("/users/{userId}/favourites", favouriteHandler.GetFavourites).Methods("GET")
	r.HandleFunc("/users/{userId}/favourites/{assetId}", favouriteHandler.AddFavourite).Methods("POST")
	r.HandleFunc("/users/{userId}/favourites/{assetId}", favouriteHandler.RemoveFavourite).Methods("DELETE")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
