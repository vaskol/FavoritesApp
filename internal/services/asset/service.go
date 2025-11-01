package assetServices

import (
	"log"

	"assetsApp/internal/models"
	"assetsApp/internal/storage"
)

type AssetService struct {
	store storage.AssetStore
}

func NewAssetService(store storage.AssetStore) *AssetService {
	return &AssetService{store: store}
}

func (s *AssetService) GetAssets(userID string) []models.Asset {
	log.Printf("Service: GetAssets called for user %s", userID)
	return s.store.Get(userID)
}

func (s *AssetService) AddAsset(userID string, asset models.Asset) {
	log.Printf("Service: AddAsset called for user %s", userID)
	s.store.Add(userID, asset)
}

func (s *AssetService) RemoveAsset(userID, assetID string) bool {
	log.Printf("Service: RemoveAsset called for user %s, asset %s", userID, assetID)
	return s.store.Remove(userID, assetID)
}

func (s *AssetService) EditDescription(userID, assetID, description string) bool {
	log.Printf("Service: EditDescription called for user %s, asset %s", userID, assetID)
	return s.store.EditDescription(userID, assetID, description)
}
