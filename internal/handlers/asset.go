package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"assetsApp/internal/models"
	services "assetsApp/internal/services/asset"

	"github.com/gorilla/mux"
)

type AssetHandler struct {
	service *services.AssetService
}

func NewAssetHandler(service *services.AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

func (h *AssetHandler) GetAssets(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]
	log.Printf("GetAssets called for user %s", userID)
	favs := h.service.GetAssets(userID)
	json.NewEncoder(w).Encode(favs)
	log.Printf("GetAssets completed for user %s", userID)
}

func (h *AssetHandler) AddAsset(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]
	log.Printf("AddAsset called for user %s", userID)

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
	log.Printf("Adding asset of type %s for user %s", assetType, userID)

	var asset models.Asset

	switch assetType {
	case "chart":
		data := []models.ChartData{}
		if d, ok := body["data"].([]interface{}); ok {
			for _, item := range d {
				m := item.(map[string]interface{})
				data = append(data, models.ChartData{
					DatapointCode: m["datapoint_code"].(string),
					Value:         m["value"].(float64),
				})
			}
		}
		a := &models.Chart{
			ID:          body["id"].(string),
			Title:       body["title"].(string),
			Description: body["description"].(string),
			XAxisTitle:  body["x_axis_title"].(string),
			YAxisTitle:  body["y_axis_title"].(string),
			Data:        data,
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

	h.service.AddAsset(userID, asset)
	w.WriteHeader(http.StatusCreated)

	switch a := asset.(type) {
	case *models.Chart:
		json.NewEncoder(w).Encode(a)
	case *models.Insight:
		json.NewEncoder(w).Encode(a)
	case *models.Audience:
		json.NewEncoder(w).Encode(a)
	default:
		http.Error(w, "Unknown asset type", http.StatusInternalServerError)
	}
	log.Printf("AddAsset completed for user %s", userID)
}

func (h *AssetHandler) RemoveAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	assetID := vars["assetId"]
	log.Printf("RemoveAsset called for user %s, asset %s", userID, assetID)
	if !h.service.RemoveAsset(userID, assetID) {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("RemoveAsset completed for user %s, asset %s", userID, assetID)
}

func (h *AssetHandler) EditAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	assetID := vars["assetId"]
	log.Printf("EditAsset called for user %s, asset %s", userID, assetID)
	var body struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.service.EditDescription(userID, assetID, body.Description) {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("EditAsset completed for user %s, asset %s", userID, assetID)
}
