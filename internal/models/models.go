package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base model with UUID
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserRole string

const (
	RoleUser     UserRole = "USER"
	RoleAdmin    UserRole = "ADMIN"
	RoleSecurity UserRole = "SECURITY"
)

type User struct {
	Base
	Name         string   `json:"name"`
	Email        string   `gorm:"uniqueIndex:idx_users_email" json:"email"`
	PasswordHash string   `json:"-"`
	Phone        string   `json:"phone"`
	Role         UserRole `gorm:"default:'USER'" json:"role"`
}

type ItemCategory struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name string    `json:"name"`
}

type CampusLocation struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type Asset struct {
	Base
	OwnerID         uuid.UUID    `json:"owner_id"`
	Owner           User         `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	CategoryID      uuid.UUID    `json:"category_id"`
	Category        ItemCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Description     string       `json:"description"`
	PrivateImageURL string       `json:"private_image_url"` // Not shown publicly
	LostMode        bool         `json:"lost_mode"`
	QRCodeURL       string       `json:"qr_code_url"`
}

type FoundEvent struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AssetID    uuid.UUID      `json:"asset_id"`
	Asset      Asset          `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	FinderID   *uuid.UUID     `json:"finder_id"` // Can be nil if anonymous scan? No, prompt says FinderID FK. Assuming authenticated finder or security.
	Finder     *User          `gorm:"foreignKey:FinderID" json:"finder,omitempty"`
	LocationID uuid.UUID      `json:"location_id"`
	Location   CampusLocation `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Note       string         `json:"note"`
	ImageURL   string         `json:"image_url"`
	CreatedAt  time.Time      `json:"created_at"`
}

type ItemStatus string

const (
	ItemStatusOpen     ItemStatus = "OPEN"
	ItemStatusClaimed  ItemStatus = "CLAIMED"
	ItemStatusResolved ItemStatus = "RESOLVED"
)

// Finder-First Item
type Item struct {
	Base
	Title                string         `json:"title"`
	CategoryID           uuid.UUID      `json:"category_id"`
	Category             ItemCategory   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	LocationID           uuid.UUID      `json:"location_id"`
	Location             CampusLocation `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	ImageURL             string         `json:"image_url"`
	VerificationQuestion string         `json:"verification_question"`
	VerificationAnswer   string         `json:"-"` // Plaintext per request, but hidden from JSON
	FinderID             uuid.UUID      `json:"finder_id"`
	Finder               User           `gorm:"foreignKey:FinderID" json:"finder,omitempty"`
	Status               ItemStatus     `gorm:"default:'OPEN'" json:"status"`
}

type ClaimStatus string

const (
	ClaimStatusPending  ClaimStatus = "PENDING"
	ClaimStatusApproved ClaimStatus = "APPROVED"
	ClaimStatusRejected ClaimStatus = "REJECTED"
)

type Claim struct {
	ID          uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ItemID      uuid.UUID   `json:"item_id"`
	Item        Item        `gorm:"foreignKey:ItemID" json:"item,omitempty"`
	OwnerID     uuid.UUID   `json:"owner_id"`
	Owner       User        `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	AnswerInput string      `json:"answer_input"`
	Status      ClaimStatus `gorm:"default:'PENDING'" json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	RefType   string    `json:"ref_type"` // e.g., "ASSET_FOUND", "CLAIM_UPDATE"
	RefID     uuid.UUID `json:"ref_id"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
