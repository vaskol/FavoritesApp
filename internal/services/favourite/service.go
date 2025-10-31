package services

import (
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
	return s.store.Get(userID)
}

func (s *FavouriteService) AddFavourite(userID string, asset models.Asset) {
	s.store.Add(userID, asset)
}

func (s *FavouriteService) RemoveFavourite(userID, assetID string) bool {
	return s.store.Remove(userID, assetID)
}

func (s *FavouriteService) EditDescription(userID, assetID, newDesc string) bool {
	return s.store.EditDescription(userID, assetID, newDesc)
}
