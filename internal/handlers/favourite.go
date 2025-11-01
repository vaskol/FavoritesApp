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

	if !h.service.AddFavourite(userID, assetID) {
		http.Error(w, "Asset not found", http.StatusNotFound)
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
