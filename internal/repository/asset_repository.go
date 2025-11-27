package repository

import (
	"campus-lost-and-found/internal/models"

	"gorm.io/gorm"
)

type AssetRepository struct {
	DB *gorm.DB
}

func NewAssetRepository(db *gorm.DB) *AssetRepository {
	return &AssetRepository{DB: db}
}

func (r *AssetRepository) Create(asset *models.Asset) error {
	return r.DB.Create(asset).Error
}

func (r *AssetRepository) FindByID(id string) (*models.Asset, error) {
	var asset models.Asset
	err := r.DB.Preload("Category").Preload("Owner").First(&asset, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepository) FindByOwnerID(ownerID string) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.DB.Preload("Category").Where("owner_id = ?", ownerID).Find(&assets).Error
	return assets, err
}

func (r *AssetRepository) Update(asset *models.Asset) error {
	return r.DB.Save(asset).Error
}

func (r *AssetRepository) CreateFoundEvent(event *models.FoundEvent) error {
	return r.DB.Create(event).Error
}

func (r *AssetRepository) GetFoundEvents(assetID string) ([]models.FoundEvent, error) {
	var events []models.FoundEvent
	err := r.DB.Preload("Location").Preload("Finder").Where("asset_id = ?", assetID).Find(&events).Error
	return events, err
}

func (r *AssetRepository) FindLostAssets() ([]models.Asset, error) {
	var assets []models.Asset
	err := r.DB.Where("lost_mode = ?", true).Find(&assets).Error
	return assets, err
}
