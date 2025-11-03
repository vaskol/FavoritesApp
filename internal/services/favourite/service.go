package favouriteServices

import (
	"log"

	"assetsApp/internal/models"
	"assetsApp/internal/storage"
)

type FavouriteService struct {
	store storage.AssetStore
}

func NewFavouriteService(store storage.AssetStore) *FavouriteService {
	return &FavouriteService{store: store}
}

func (s *FavouriteService) GetFavourites(userID string) []models.Favourite {
	log.Printf("Service: GetFavourites called for user %s", userID)
	return s.store.GetFavourites(userID)
}

func (s *FavouriteService) AddFavourite(userID, assetID, assetType string) bool {
	log.Printf("Service: AddFavourite called for user %s, asset %s", userID, assetID)
	return s.store.AddFavourite(userID, assetID, assetType)
}

func (s *FavouriteService) RemoveFavourite(userID, assetID string) bool {
	log.Printf("Service: RemoveFavourite called for user %s, asset %s", userID, assetID)
	return s.store.RemoveFavourite(userID, assetID)
}
