package storage

import (
	"assetsApp/internal/models"
	"log"
	"sync"
)

type MemoryStore struct {
	mu         sync.RWMutex
	store      map[string][]models.Asset
	favourites map[string][]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store:      make(map[string][]models.Asset),
		favourites: make(map[string][]string),
	}
}

func (m *MemoryStore) Get(userID string) []models.Asset {
	log.Printf("Storage: Get called for user %s", userID)
	m.mu.RLock()
	defer m.mu.RUnlock()

	assets := make([]models.Asset, len(m.store[userID]))
	copy(assets, m.store[userID])
	return assets
}

func (m *MemoryStore) Add(userID string, asset models.Asset) {
	log.Printf("Storage: Add called for user %s", userID)
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, a := range m.store[userID] {
		if a.GetID() == asset.GetID() {
			log.Printf("Storage: Add called for user %s, asset %s already exists", userID, asset.GetID())
			// already exists, ignore or return an error
			return
		}
	}
	m.store[userID] = append(m.store[userID], asset)
}

func (m *MemoryStore) Remove(userID, assetID string) bool {
	log.Printf("Storage: Remove called for user %s, asset %s", userID, assetID)
	m.mu.Lock()
	defer m.mu.Unlock()
	assets, ok := m.store[userID]
	if !ok {
		return false
	}
	for i := range assets {
		if assets[i].GetID() == assetID {
			m.store[userID] = append(assets[:i], assets[i+1:]...)
			return true
		}
	}
	return false
}

func (m *MemoryStore) EditDescription(userID, assetID, desc string) bool {
	log.Printf("Storage: EditDescription called for user %s, asset %s", userID, assetID)
	m.mu.Lock()
	defer m.mu.Unlock()
	assets, ok := m.store[userID]
	if !ok {
		return false
	}
	for i := range assets {
		if assets[i].GetID() == assetID {
			assets[i].SetDescription(desc)
			m.store[userID][i] = assets[i]
			return true
		}
	}
	return false
}
func (m *MemoryStore) AddFavourite(userID, assetID, assetType string) bool {
	log.Printf("Storage: AddFavourite called for user %s, asset %s", userID, assetID)
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, fav := range m.favourites[userID] {
		if fav == assetID {
			return false // already exists
		}
	}
	m.favourites[userID] = append(m.favourites[userID], assetID)
	return true
}

func (m *MemoryStore) RemoveFavourite(userID, assetID string) bool {
	log.Printf("Storage: RemoveFavourite called for user %s, asset %s", userID, assetID)
	m.mu.Lock()
	defer m.mu.Unlock()

	favs := m.favourites[userID]
	for i, fav := range favs {
		if fav == assetID {
			m.favourites[userID] = append(favs[:i], favs[i+1:]...)
			return true
		}
	}
	return false
}

func (m *MemoryStore) GetFavourites(userID string) []models.Favourite {
	log.Printf("Storage: GetFavourites called for user %s", userID)
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []models.Favourite
	allAssets := m.store[userID]
	for _, favID := range m.favourites[userID] {
		for _, asset := range allAssets {
			if asset.GetID() == favID {
				result = append(result, models.Favourite{UserID: userID, Asset: asset})
			}
		}
	}
	return result
}
