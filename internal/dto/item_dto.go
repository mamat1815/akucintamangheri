package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateFoundItemRequest struct {
	Title                string    `json:"title" binding:"required"`
	CategoryID           uuid.UUID `json:"category_id" binding:"required"`
	LocationID           uuid.UUID `json:"location_id" binding:"required"`
	ImageURL             string    `json:"image_url"`
	VerificationQuestion string    `json:"verification_question" binding:"required"`
	VerificationAnswer   string    `json:"verification_answer" binding:"required"`
}

type CreateLostItemRequest struct {
	Title            string    `json:"title" binding:"required"`
	CategoryID       uuid.UUID `json:"category_id" binding:"required"`
	Description      string    `json:"description"`
	LocationLastSeen string    `json:"location_last_seen" binding:"required"`
	DateLost         string    `json:"date_lost" binding:"required"` // Format YYYY-MM-DD
	ImageURL         string    `json:"image_url"`
}

type ItemResponse struct {
	ID                   uuid.UUID `json:"id"`
	Title                string    `json:"title"`
	CategoryID           uuid.UUID `json:"category_id"`
	CategoryName         string    `json:"category_name,omitempty"`
	LocationID           uuid.UUID `json:"location_id"`
	LocationName         string    `json:"location_name,omitempty"`
	ImageURL             string    `json:"image_url"`
	VerificationQuestion string    `json:"verification_question"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
}

type CreateClaimRequest struct {
	AnswerInput string `json:"answer_input" binding:"required"`
}

type ClaimResponse struct {
	ID          uuid.UUID `json:"id"`
	ItemID      uuid.UUID `json:"item_id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	AnswerInput string    `json:"answer_input"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type DecideClaimRequest struct {
	Status string `json:"status" binding:"required,oneof=APPROVED REJECTED"`
}
