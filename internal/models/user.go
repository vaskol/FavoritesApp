package models

// Favourite links a user to an Asset
type Favourite struct {
	UserID string `json:"user_id"`
	Asset  Asset  `json:"asset"`
}
