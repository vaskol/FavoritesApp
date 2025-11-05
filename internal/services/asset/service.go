package assetServices

import (
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
	return s.store.Get(userID)
}

func (s *AssetService) AddAsset(userID string, asset models.Asset) {
	s.store.Add(userID, asset)
}

func (s *AssetService) RemoveAsset(userID, assetID string) bool {
	return s.store.Remove(userID, assetID)
}

func (s *AssetService) EditDescription(userID, assetID, description string) bool {
	return s.store.EditDescription(userID, assetID, description)
}
