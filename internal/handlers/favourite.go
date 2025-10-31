package handlers

import (
	"encoding/json"
	"net/http"

	services "favouritesApp/internal/favourite"
	"favouritesApp/internal/models"

	"github.com/gorilla/mux"
)

type FavouriteHandler struct {
	service *services.FavouriteService
}

func NewFavouriteHandler(service *services.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{service: service}
}

func (h *FavouriteHandler) GetFavourites(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]
	favs := h.service.GetFavourites(userID)
	json.NewEncoder(w).Encode(favs)
}

func (h *FavouriteHandler) AddFavourite(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]

	// Simple approach: decode into a map to detect type
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	assetType, ok := body["type"].(string)
	if !ok {
		http.Error(w, "Asset type required", http.StatusBadRequest)
		return
	}

	var asset models.Asset

	switch assetType {
	case "chart":
		a := &models.Chart{
			ID:          body["id"].(string),
			Description: body["description"].(string),
		}
		asset = a
	case "insight":
		a := &models.Insight{
			ID:          body["id"].(string),
			Description: body["description"].(string),
		}
		asset = a
	case "audience":
		a := &models.Audience{
			ID:          body["id"].(string),
			Description: body["description"].(string),
			Gender:      body["gender"].(string),
			Country:     body["country"].(string),
			AgeGroup:    body["age_group"].(string),
			SocialHours: int(body["social_hours"].(float64)),
			Purchases:   int(body["purchases"].(float64)),
		}
		asset = a
	default:
		http.Error(w, "Unknown asset type", http.StatusBadRequest)
		return
	}

	h.service.AddFavourite(userID, asset)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(asset)
}

func (h *FavouriteHandler) RemoveFavourite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !h.service.RemoveFavourite(vars["userId"], vars["assetId"]) {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *FavouriteHandler) EditFavourite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var body struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.service.EditDescription(vars["userId"], vars["assetId"], body.Description) {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
