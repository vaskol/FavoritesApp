package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"assetsApp/internal/models"
)

// CachedStore is a decorator that adds caching via Redis for favourites.
// It implements the same AssetStore interface by delegating to the inner db store
// except for favourites operations which have cache behaviour.
type CachedStore struct {
	db    AssetStore
	cache *RedisClient
}

func NewCachedStore(db AssetStore, cache *RedisClient) *CachedStore {
	return &CachedStore{
		db:    db,
		cache: cache,
	}
}

func favsCacheKey(userID string) string {
	return fmt.Sprintf("favourites:%s", userID)
}

// ----- Delegate methods for assets (unchanged behaviour) -----

func (c *CachedStore) Get(userID string) []models.Asset {
	return c.db.Get(userID)
}

func (c *CachedStore) Add(userID string, asset models.Asset) {
	c.db.Add(userID, asset)
}

func (c *CachedStore) Remove(userID, assetID string) bool {
	// If removing an asset, also best-effort invalidate favourites cache for the owner
	res := c.db.Remove(userID, assetID)
	if res && c.cache != nil {
		_ = c.cache.Del(context.Background(), favsCacheKey(userID))
	}
	return res
}

func (c *CachedStore) EditDescription(userID, assetID, newDesc string) bool {
	res := c.db.EditDescription(userID, assetID, newDesc)
	if res && c.cache != nil {
		_ = c.cache.Del(context.Background(), favsCacheKey(userID))
	}
	return res
}

// ----- Favourites with caching -----

func (c *CachedStore) GetFavourites(userID string) []models.Favourite {
	// Try cache first
	ctx := context.Background()
	if c.cache != nil {
		if cached, err := c.cache.Get(ctx, favsCacheKey(userID)); err == nil && cached != "" {
			var favs []models.Favourite
			if err := json.Unmarshal([]byte(cached), &favs); err == nil {
				return favs
			}
			// if unmarshal fails, fallthrough to DB fetch
			log.Printf("cached_store: failed to unmarshal favourites cache for user %s: %v", userID, err)
		}
	}

	// Fallback to DB
	favs := c.db.GetFavourites(userID)

	// Write back to cache (best-effort)
	if c.cache != nil {
		if b, err := json.Marshal(favs); err == nil {
			_ = c.cache.Set(ctx, favsCacheKey(userID), string(b))
		} else {
			log.Printf("cached_store: failed to marshal favourites for caching: %v", err)
		}
	}

	return favs
}

func (c *CachedStore) AddFavourite(userID, assetID, assetType string) bool {
	res := c.db.AddFavourite(userID, assetID, assetType)
	if res && c.cache != nil {
		_ = c.cache.Del(context.Background(), favsCacheKey(userID))
	}
	return res
}

func (c *CachedStore) RemoveFavourite(userID, assetID string) bool {
	res := c.db.RemoveFavourite(userID, assetID)
	if res && c.cache != nil {
		_ = c.cache.Del(context.Background(), favsCacheKey(userID))
	}
	return res
}
