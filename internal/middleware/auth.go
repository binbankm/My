package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/binbankm/My/internal/util"
	"github.com/gin-gonic/gin"
)

// getAllowedOrigins returns configured CORS origins
func getAllowedOrigins() string {
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		// Default to localhost for development
		return "http://localhost:3000,http://localhost:8888"
	}
	return origins
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := getAllowedOrigins()
		origin := c.GetHeader("Origin")
		
		// Check if origin is allowed
		if origin != "" && contains(strings.Split(allowedOrigins, ","), origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else if os.Getenv("GIN_MODE") != "release" {
			// DEVELOPMENT ONLY: Allow all origins in dev mode
			// WARNING: This should NEVER be used in production - set GIN_MODE=release
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// contains checks if a string slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == item {
			return true
		}
	}
	return false
}

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
