package services

import (
	"campus-lost-and-found/internal/dto"
	"campus-lost-and-found/internal/models"
	"campus-lost-and-found/internal/repository"
	"fmt"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type AssetService struct {
	Repo          *repository.AssetRepository
	UploadService *UploadService
	NotifService  *NotificationService
}

func NewAssetService(repo *repository.AssetRepository, uploadService *UploadService, notifService *NotificationService) *AssetService {
	return &AssetService{
		Repo:          repo,
		UploadService: uploadService,
		NotifService:  notifService,
	}
}

func (s *AssetService) CreateAsset(req dto.CreateAssetRequest, ownerID uuid.UUID) (*dto.AssetResponse, error) {
	asset := &models.Asset{
		OwnerID:         ownerID,
		CategoryID:      req.CategoryID,
		Description:     req.Description,
		PrivateImageURL: req.PrivateImageURL,
		LostMode:        req.LostMode,
	}

	if err := s.Repo.Create(asset); err != nil {
		return nil, err
	}

	// Generate QR Code
	qrContent := fmt.Sprintf("https://campus-lost-found.app/scan/%s", asset.ID.String())
	png, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	// Upload QR Code
	filename := fmt.Sprintf("qr_%s.png", asset.ID.String())
	qrURL, err := s.UploadService.UploadBytes(png, filename)
	if err != nil {
		return nil, err
	}

	asset.QRCodeURL = qrURL
	s.Repo.Update(asset)

	return &dto.AssetResponse{
		ID:              asset.ID,
		OwnerID:         asset.OwnerID,
		CategoryID:      asset.CategoryID,
		Description:     asset.Description,
		PrivateImageURL: asset.PrivateImageURL,
		LostMode:        asset.LostMode,
		QRCodeURL:       asset.QRCodeURL,
		CreatedAt:       asset.CreatedAt,
	}, nil
}

func (s *AssetService) GetAsset(id string) (*models.Asset, error) {
	return s.Repo.FindByID(id)
}

func (s *AssetService) UpdateLostMode(id string, lostMode bool, userID uuid.UUID) error {
	asset, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}

	if asset.OwnerID != userID {
		return fmt.Errorf("unauthorized")
	}

	asset.LostMode = lostMode
	return s.Repo.Update(asset)
}

func (s *AssetService) ReportFound(assetID string, req dto.ReportFoundRequest) error {
	asset, err := s.Repo.FindByID(assetID)
	if err != nil {
		return err
	}

	event := &models.FoundEvent{
		AssetID:    asset.ID,
		LocationID: req.LocationID,
		Note:       req.Note,
		ImageURL:   req.ImageURL,
	}

	if err := s.Repo.CreateFoundEvent(event); err != nil {
		return err
	}

	// Notify Owner
	s.NotifService.CreateNotification(
		asset.OwnerID,
		"Asset Scanned!",
		fmt.Sprintf("Your asset '%s' was scanned at a location.", asset.Description),
		"ASSET_FOUND",
		event.ID,
	)

	return nil
}

func (s *AssetService) GetFoundEvents(assetID string) ([]models.FoundEvent, error) {
	return s.Repo.GetFoundEvents(assetID)
}

func (s *AssetService) GetLostAssets() ([]dto.AssetResponse, error) {
	assets, err := s.Repo.FindLostAssets()
	if err != nil {
		return nil, err
	}

	var responses []dto.AssetResponse
	for _, asset := range assets {
		responses = append(responses, dto.AssetResponse{
			ID:           asset.ID,
			OwnerID:      asset.OwnerID,
			CategoryID:   asset.CategoryID,
			CategoryName: asset.Category.Name, // Preloaded in Repo
			Description:  asset.Description,
			LostMode:     asset.LostMode,
			QRCodeURL:    asset.QRCodeURL,
			CreatedAt:    asset.CreatedAt,
			// PrivateImageURL is intentionally omitted for public feed
		})
	}
	return responses, nil
}

func (s *AssetService) GetUserAssets(userID uuid.UUID) ([]dto.AssetResponse, error) {
	assets, err := s.Repo.FindByOwnerID(userID.String())
	if err != nil {
		return nil, err
	}

	var responses []dto.AssetResponse
	for _, asset := range assets {
		responses = append(responses, dto.AssetResponse{
			ID:              asset.ID,
			OwnerID:         asset.OwnerID,
			CategoryID:      asset.CategoryID,
			CategoryName:    asset.Category.Name,
			Description:     asset.Description,
			PrivateImageURL: asset.PrivateImageURL, // Owner can see private image
			LostMode:        asset.LostMode,
			QRCodeURL:       asset.QRCodeURL,
			CreatedAt:       asset.CreatedAt,
		})
	}
	return responses, nil
}
