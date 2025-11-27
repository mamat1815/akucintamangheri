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
	RoleStudent  UserRole = "STUDENT"
	RoleStaff    UserRole = "STAFF"
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

type ItemType string

const (
	ItemTypeLost  ItemType = "LOST"
	ItemTypeFound ItemType = "FOUND"
)

// Finder-First Item (and now Owner-First Lost Item)
type Item struct {
	Base
	Title                string         `json:"title"`
	Type                 ItemType       `gorm:"default:'FOUND'" json:"type"`
	CategoryID           uuid.UUID      `json:"category_id"`
	Category             ItemCategory   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	LocationID           *uuid.UUID     `json:"location_id"` // Nullable for Lost items
	Location             *CampusLocation `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	LocationDescription  string         `json:"location_description"` // For Lost items (free text)
	ImageURL             string         `json:"image_url"`
	VerificationQuestion string         `json:"verification_question"`
	VerificationAnswer   string         `json:"-"` // Plaintext per request, but hidden from JSON
	FinderID             *uuid.UUID     `json:"finder_id"` // Nullable for Lost items
	Finder               *User          `gorm:"foreignKey:FinderID" json:"finder,omitempty"`
	OwnerID              *uuid.UUID     `json:"owner_id"` // Nullable for Found items
	Owner                *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	DateLost             *time.Time     `json:"date_lost"`
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
