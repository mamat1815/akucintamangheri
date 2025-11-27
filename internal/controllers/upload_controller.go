package controllers

import (
	"campus-lost-and-found/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	Service *services.UploadService
}

func NewUploadController(service *services.UploadService) *UploadController {
	return &UploadController{Service: service}
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a file to storage
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /upload [post]
func (ctrl *UploadController) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	url, err := ctrl.Service.UploadFile(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
