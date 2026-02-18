package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListDatabases returns list of databases (placeholder)
func ListDatabases(c *gin.Context) {
	// This is a simplified implementation
	// In production, you would connect to MySQL/PostgreSQL and list actual databases
	c.JSON(http.StatusOK, []gin.H{
		{"name": "example_db", "type": "mysql", "size": "10MB"},
	})
}

// CreateDatabase creates a new database (placeholder)
func CreateDatabase(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Type string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In production, execute actual database creation commands
	c.JSON(http.StatusOK, gin.H{"message": "Database created successfully"})
}

// DeleteDatabase deletes a database (placeholder)
func DeleteDatabase(c *gin.Context) {
	name := c.Param("name")
	
	// In production, execute actual database deletion commands
	c.JSON(http.StatusOK, gin.H{"message": "Database " + name + " deleted successfully"})
}
