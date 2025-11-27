package controllers

import (
	"campus-lost-and-found/internal/dto"
	"campus-lost-and-found/internal/middleware"
	"campus-lost-and-found/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AssetController struct {
	Service *services.AssetService
}

func NewAssetController(service *services.AssetService) *AssetController {
	return &AssetController{Service: service}
}

// CreateAsset godoc
// @Summary Create a new asset
// @Description Create asset and generate QR code
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAssetRequest true "Create Asset Request"
// @Success 200 {object} dto.AssetResponse
// @Failure 400 {object} map[string]string
// @Router /assets [post]
func (ctrl *AssetController) CreateAsset(c *gin.Context) {
	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	res, err := ctrl.Service.CreateAsset(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetAsset godoc
// @Summary Get asset details
// @Description Get asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Asset ID"
// @Success 200 {object} dto.AssetResponse
// @Failure 404 {object} map[string]string
// @Router /assets/{id} [get]
func (ctrl *AssetController) GetAsset(c *gin.Context) {
	id := c.Param("id")
	asset, err := ctrl.Service.GetAsset(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Asset not found"})
		return
	}

	// Privacy check: If not owner, hide private details
	userID := middleware.GetUserID(c)
	isOwner := asset.OwnerID == userID

	res := dto.AssetResponse{
		ID:           asset.ID,
		OwnerID:      asset.OwnerID,
		CategoryID:   asset.CategoryID,
		CategoryName: asset.Category.Name,
		Description:  asset.Description,
		LostMode:     asset.LostMode,
		QRCodeURL:    asset.QRCodeURL,
		CreatedAt:    asset.CreatedAt,
	}

	if isOwner {
		res.PrivateImageURL = asset.PrivateImageURL
	}

	c.JSON(http.StatusOK, res)
}

// UpdateLostMode godoc
// @Summary Update lost mode
// @Description Enable or disable lost mode
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Asset ID"
// @Param request body dto.UpdateLostModeRequest true "Update Lost Mode Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /assets/{id}/lost-mode [put]
func (ctrl *AssetController) UpdateLostMode(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateLostModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	err := ctrl.Service.UpdateLostMode(id, req.LostMode, userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lost mode updated"})
}

// ReportFound godoc
// @Summary Report asset found (Scan QR)
// @Description Report that an asset has been found by scanning QR
// @Tags assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Param request body dto.ReportFoundRequest true "Report Found Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /assets/{id}/report-found [post]
func (ctrl *AssetController) ReportFound(c *gin.Context) {
	id := c.Param("id")
	var req dto.ReportFoundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.Service.ReportFound(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Asset reported found"})
}

// GetFoundEvents godoc
// @Summary Get found events for an asset
// @Description Get history of found events
// @Tags assets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Asset ID"
// @Success 200 {object} []models.FoundEvent
// @Failure 404 {object} map[string]string
// @Router /assets/{id}/found-events [get]
func (ctrl *AssetController) GetFoundEvents(c *gin.Context) {
	id := c.Param("id")
	events, err := ctrl.Service.GetFoundEvents(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

// GetLostAssets godoc
// @Summary Get all lost assets
// @Description Get a list of assets reported as lost
// @Tags assets
// @Accept json
// @Produce json
// @Success 200 {object} []dto.AssetResponse
// @Failure 500 {object} map[string]string
// @Router /assets/lost [get]
func (ctrl *AssetController) GetLostAssets(c *gin.Context) {
	assets, err := ctrl.Service.GetLostAssets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}
