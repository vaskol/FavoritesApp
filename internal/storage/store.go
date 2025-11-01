package storage

import "assetsApp/internal/models"

type AssetStore interface {
	Get(userID string) []models.Asset
	Add(userID string, asset models.Asset)
	Remove(userID, assetID string) bool
	EditDescription(userID, assetID, newDesc string) bool

	GetFavourites(userID string) []models.Favourite
	AddFavourite(userID, assetID, assetType string) bool
	RemoveFavourite(userID, assetID string) bool
}
