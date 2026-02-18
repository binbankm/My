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
		// Default origins for development
		// In production, this will fall back to allowing same-origin requests
		return "http://localhost:3000,http://localhost:8888,http://127.0.0.1:3000,http://127.0.0.1:8888"
	}
	return origins
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// If no Origin header, this is a same-origin request, allow it
		if origin == "" {
			c.Next()
			return
		}
		
		// Check if origin is explicitly allowed in configuration
		allowedOrigins := getAllowedOrigins()
		originsList := strings.Split(allowedOrigins, ",")
		
		if contains(originsList, origin) {
			// Origin is in the allowed list
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else if os.Getenv("GIN_MODE") != "release" {
			// DEVELOPMENT ONLY: Allow all origins in dev mode
			// WARNING: This should NEVER be used in production - set GIN_MODE=release
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// Production mode: Allow the request origin if it matches the host
			// This handles the case where the frontend is served from the same server
			scheme := "http"
			if c.Request.TLS != nil {
				scheme = "https"
			}
			host := c.Request.Host
			requestOrigin := scheme + "://" + host
			
			// If the origin matches where we're serving from, allow it
			if origin == requestOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				// Origin not allowed - do not set CORS headers
				// This will cause the browser to block the request
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
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
