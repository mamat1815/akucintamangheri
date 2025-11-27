package controllers

import (
	"campus-lost-and-found/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EnumerationController struct {
	Repo *repository.EnumerationRepository
}

func NewEnumerationController(repo *repository.EnumerationRepository) *EnumerationController {
	return &EnumerationController{Repo: repo}
}

// GetCategories godoc
// @Summary Get item categories
// @Description Get list of item categories
// @Tags enumerations
// @Accept json
// @Produce json
// @Success 200 {object} []models.ItemCategory
// @Router /enumerations/item-categories [get]
func (ctrl *EnumerationController) GetCategories(c *gin.Context) {
	categories, err := ctrl.Repo.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// GetLocations godoc
// @Summary Get campus locations
// @Description Get list of campus locations
// @Tags enumerations
// @Accept json
// @Produce json
// @Success 200 {object} []models.CampusLocation
// @Router /enumerations/campus-locations [get]
func (ctrl *EnumerationController) GetLocations(c *gin.Context) {
	locations, err := ctrl.Repo.GetLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locations)
}
