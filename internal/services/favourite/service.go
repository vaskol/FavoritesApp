package services

import (
	"log"

	"favouritesApp/internal/models"
	"favouritesApp/internal/storage"
)

type FavouriteService struct {
	store *storage.MemoryStore
}

func NewFavouriteService(store *storage.MemoryStore) *FavouriteService {
	return &FavouriteService{store: store}
}

func (s *FavouriteService) GetFavourites(userID string) []models.Asset {
	log.Printf("Service: GetFavourites called for user %s", userID)
	return s.store.Get(userID)
}

func (s *FavouriteService) AddFavourite(userID string, asset models.Asset) {
	log.Printf("Service: AddFavourite called for user %s", userID)
	s.store.Add(userID, asset)
}

func (s *FavouriteService) RemoveFavourite(userID, assetID string) bool {
	log.Printf("Service: RemoveFavourite called for user %s, asset %s", userID, assetID)
	return s.store.Remove(userID, assetID)
}

func (s *FavouriteService) EditDescription(userID, assetID, newDesc string) bool {
	log.Printf("Service: EditDescription called for user %s, asset %s", userID, assetID)
	return s.store.EditDescription(userID, assetID, newDesc)
}
