package api

import (
	"net/http"

	"github.com/binbankm/My/internal/model"
	"github.com/gin-gonic/gin"
)

// GetSettings returns all settings
func GetSettings(c *gin.Context) {
	var settings []model.Settings
	if err := model.DB.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSettings updates settings
func UpdateSettings(c *gin.Context) {
	var req []struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, setting := range req {
		model.DB.Where("key = ?", setting.Key).
			Assign(model.Settings{Value: setting.Value}).
			FirstOrCreate(&model.Settings{Key: setting.Key})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}
