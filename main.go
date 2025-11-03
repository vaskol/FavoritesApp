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

	// // -------------------- REDIS --------------------
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

	// ctx := context.Background()

	// err := rdb.Set(ctx, "key", "value", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// val, err := rdb.Get(ctx, "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("key", val)

	// val2, err := rdb.Get(ctx, "key2").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("key2", val2)
	// }

	// -------------------- STORAGE --------------------
	// Uncomment one depending on which store you want

	// Memory store
	// store := storage.NewMemoryStore()

	// Postgres store
	db, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432/assetdb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	dbStore := storage.NewPostgresStore(db)
	redisClient := storage.NewRedisClient("localhost:6379")
	store := storage.NewCachedStore(dbStore, redisClient)

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
