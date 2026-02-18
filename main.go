package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/binbankm/My/internal/api"
	"github.com/binbankm/My/internal/middleware"
	"github.com/binbankm/My/internal/model"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	// Initialize database
	if err := model.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Set gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Configure trusted proxies
	// In production, set TRUSTED_PROXIES env var to your proxy IPs
	// For direct access (no proxy), use an empty list
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		// Parse comma-separated list of proxy IPs
		proxies := []string{}
		for _, proxy := range strings.Split(trustedProxies, ",") {
			proxy = strings.TrimSpace(proxy)
			if proxy != "" {
				proxies = append(proxies, proxy)
			}
		}
		if err := r.SetTrustedProxies(proxies); err != nil {
			log.Printf("Warning: Failed to set trusted proxies: %v", err)
		}
	} else {
		// No trusted proxies - direct access
		if err := r.SetTrustedProxies(nil); err != nil {
			log.Printf("Warning: Failed to set trusted proxies: %v", err)
		}
	}

	// CORS middleware
	r.Use(middleware.CORS())

	// API routes
	apiGroup := r.Group("/api")
	{
		// Auth routes
		auth := apiGroup.Group("/auth")
		{
			auth.POST("/login", api.Login)
			auth.POST("/logout", middleware.AuthMiddleware(), api.Logout)
			auth.GET("/info", middleware.AuthMiddleware(), api.GetUserInfo)
		}

		// Protected routes
		protected := apiGroup.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// System monitoring
			system := protected.Group("/system")
			{
				system.GET("/info", api.GetSystemInfo)
				system.GET("/stats", api.GetSystemStats)
			}

			// Docker management
			docker := protected.Group("/docker")
			{
				docker.GET("/containers", api.ListContainers)
				docker.GET("/containers/:id", api.GetContainer)
				docker.POST("/containers/:id/start", api.StartContainer)
				docker.POST("/containers/:id/stop", api.StopContainer)
				docker.POST("/containers/:id/restart", api.RestartContainer)
				docker.DELETE("/containers/:id", api.DeleteContainer)
				docker.GET("/images", api.ListImages)
				docker.DELETE("/images/:id", api.DeleteImage)
			}

			// File management
			files := protected.Group("/files")
			{
				files.GET("", api.ListFiles)
				files.POST("", api.CreateFile)
				files.PUT("", api.UpdateFile)
				files.DELETE("", api.DeleteFile)
				files.GET("/download", api.DownloadFile)
				files.POST("/upload", api.UploadFile)
			}

			// Database management
			database := protected.Group("/database")
			{
				database.GET("", api.ListDatabases)
				database.POST("", api.CreateDatabase)
				database.DELETE("/:name", api.DeleteDatabase)
			}

			// Settings
			settings := protected.Group("/settings")
			{
				settings.GET("", api.GetSettings)
				settings.PUT("", api.UpdateSettings)
			}
		}
	}

	// Serve frontend
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("Warning: Frontend not embedded, using API only mode: %v", err)
	} else {
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// Serve static files or index.html
			if _, err := distFS.Open(path[1:]); err != nil {
				c.FileFromFS("/", http.FS(distFS))
			} else {
				c.FileFromFS(path, http.FS(distFS))
			}
		})
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
