package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateAssetRequest struct {
	CategoryID      uuid.UUID `json:"category_id" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	PrivateImageURL string    `json:"private_image_url"`
	LostMode        bool      `json:"lost_mode"`
}

type AssetResponse struct {
	ID              uuid.UUID `json:"id"`
	OwnerID         uuid.UUID `json:"owner_id"`
	CategoryID      uuid.UUID `json:"category_id"`
	CategoryName    string    `json:"category_name,omitempty"`
	Description     string    `json:"description"`
	PrivateImageURL string    `json:"private_image_url,omitempty"` // Only for owner
	LostMode        bool      `json:"lost_mode"`
	QRCodeURL       string    `json:"qr_code_url"`
	CreatedAt       time.Time `json:"created_at"`
}

type UpdateLostModeRequest struct {
	LostMode bool `json:"lost_mode"`
}

type ReportFoundRequest struct {
	LocationID uuid.UUID `json:"location_id" binding:"required"`
	Note       string    `json:"note"`
	ImageURL   string    `json:"image_url"`
}
