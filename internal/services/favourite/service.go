package services

import (
	"log"

	"assetsApp/internal/models"
	"assetsApp/internal/storage"
)

type FavouriteService struct {
	store *storage.MemoryStore
}

func NewFavouriteService(store *storage.MemoryStore) *FavouriteService {
	return &FavouriteService{store: store}
}

func (s *FavouriteService) GetFavourites(userID string) []models.Asset {
	log.Printf("Service: GetFavourites called for user %s", userID)
	return s.store.GetFavourites(userID)
}

func (s *FavouriteService) AddFavourite(userID, assetID string) bool {
	log.Printf("Service: AddFavourite called for user %s, asset %s", userID, assetID)
	assets := s.store.Get(userID)
	var found models.Asset
	for _, a := range assets {
		if a.GetID() == assetID {
			found = a
			break
		}
	}
	if found == nil {
		return false
	}
	s.store.AddFavourite(userID, assetID)
	return true
}


func (s *FavouriteService) RemoveFavourite(userID, assetID string) bool {
	log.Printf("Service: RemoveFavourite called for user %s, asset %s", userID, assetID)
	return s.store.Remove(userID, assetID)
}
