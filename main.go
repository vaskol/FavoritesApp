package main

import (
	"context"
	"log"
	"net/http"

	"assetsApp/internal/handlers"
	assetServices "assetsApp/internal/services/asset"
	favouriteServices "assetsApp/internal/services/favourite"
	"assetsApp/internal/storage"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Println("Starting the application...")

	// -------------------- STORAGE --------------------
	// Uncomment one depending on which store you want

	// Memory store
	// store := storage.NewMemoryStore()

	// Postgres store
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432/assetdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	store := storage.NewPostgresStore(pool)

	// -------------------- SERVICES --------------------
	assetService := assetServices.NewAssetService(store)
	favouriteService := favouriteServices.NewFavouriteService(store)

	// -------------------- HANDLERS --------------------
	assetHandler := handlers.NewAssetHandler(assetService)
	favouriteHandler := handlers.NewFavouriteHandler(favouriteService)

	// -------------------- ROUTER --------------------
	r := mux.NewRouter()

	// Asset routes
	r.HandleFunc("/users/{userId}/assets", assetHandler.GetAssets).Methods("GET")
	r.HandleFunc("/users/{userId}/assets", assetHandler.AddAsset).Methods("POST")
	r.HandleFunc("/users/{userId}/assets/{assetId}", assetHandler.EditAsset).Methods("PUT")
	r.HandleFunc("/users/{userId}/assets/{assetId}", assetHandler.RemoveAsset).Methods("DELETE")

	// Favourite routes
	r.HandleFunc("/users/{userId}/favourites", favouriteHandler.GetFavourites).Methods("GET")
	r.HandleFunc("/users/{userId}/favourites/{assetId}", favouriteHandler.AddFavourite).Methods("POST")
	r.HandleFunc("/users/{userId}/favourites/{assetId}", favouriteHandler.RemoveFavourite).Methods("DELETE")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// -------------------- START SERVER --------------------
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
