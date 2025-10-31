package storage

import (
	"favouritesApp/internal/models"
	"sync"
)

type MemoryStore struct {
	mu    sync.RWMutex
	store map[string][]models.Asset
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[string][]models.Asset),
	}
}

func (m *MemoryStore) Get(userID string) []models.Asset {
	m.mu.RLock()
	defer m.mu.RUnlock()

	assets := make([]models.Asset, len(m.store[userID]))
	copy(assets, m.store[userID])
	return assets
}

func (m *MemoryStore) Add(userID string, asset models.Asset) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[userID] = append(m.store[userID], asset)
}

func (m *MemoryStore) Remove(userID, assetID string) bool {
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
