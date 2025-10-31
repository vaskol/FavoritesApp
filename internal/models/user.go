package models

// Favourite links a user to an Asset
type Favourite struct {
	UserID string
	Asset  Asset
}
