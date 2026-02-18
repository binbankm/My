package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Placeholder functions for Docker management
// Full Docker integration requires docker socket access

// ListContainers returns placeholder data
func ListContainers(c *gin.Context) {
	c.JSON(http.StatusOK, []gin.H{})
}

// GetContainer returns placeholder data
func GetContainer(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Docker integration requires docker socket access"})
}

// StartContainer placeholder
func StartContainer(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Docker integration requires docker socket access"})
}

// StopContainer placeholder
func StopContainer(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Docker integration requires docker socket access"})
}

// RestartContainer placeholder
func RestartContainer(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Docker integration requires docker socket access"})
}

// DeleteContainer placeholder
func DeleteContainer(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Docker integration requires docker socket access"})
}

// ListImages returns placeholder data
func ListImages(c *gin.Context) {
	c.JSON(http.StatusOK, []gin.H{})
}

// DeleteImage placeholder
func DeleteImage(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Docker integration requires docker socket access"})
}
