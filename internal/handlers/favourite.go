package handlers

import (
	"encoding/json"
	"net/http"

	favouriteServices "assetsApp/internal/services/favourite"

	"github.com/gorilla/mux"
)

type FavouriteHandler struct {
	service *favouriteServices.FavouriteService
}

func NewFavouriteHandler(service *favouriteServices.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{service: service}
}

func (h *FavouriteHandler) GetFavourites(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]
	favs := h.service.GetFavourites(userID)
	json.NewEncoder(w).Encode(favs)
}

func (h *FavouriteHandler) AddFavourite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	assetID := vars["assetId"]

	var body struct {
		AssetType string `json:"asset_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.service.AddFavourite(userID, assetID, body.AssetType) {
		http.Error(w, "Could not add favourite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *FavouriteHandler) RemoveFavourite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	assetID := vars["assetId"]
	if !h.service.RemoveFavourite(userID, assetID) {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
