package api

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// NginxSite represents an Nginx site configuration
type NginxSite struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Enabled   bool   `json:"enabled"`
	Content   string `json:"content,omitempty"`
	ServerName string `json:"serverName,omitempty"`
	Port      string `json:"port,omitempty"`
}

const (
	nginxSitesAvailable = "/etc/nginx/sites-available"
	nginxSitesEnabled   = "/etc/nginx/sites-enabled"
	nginxConfDir        = "/etc/nginx"
)

// ListNginxSites lists all Nginx sites
func ListNginxSites(c *gin.Context) {
	// Check if Nginx is installed
	if _, err := exec.LookPath("nginx"); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nginx is not installed"})
		return
	}

	var sites []NginxSite

	// List sites-available
	if _, err := os.Stat(nginxSitesAvailable); err == nil {
		files, err := os.ReadDir(nginxSitesAvailable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read sites-available: " + err.Error()})
			return
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			sitePath := filepath.Join(nginxSitesAvailable, file.Name())
			enabledPath := filepath.Join(nginxSitesEnabled, file.Name())
			
			_, err := os.Stat(enabledPath)
			enabled := err == nil

			// Read basic info from the config file
			content, _ := os.ReadFile(sitePath)
			serverName, port := parseNginxConfig(string(content))

			sites = append(sites, NginxSite{
				Name:       file.Name(),
				Path:       sitePath,
				Enabled:    enabled,
				ServerName: serverName,
				Port:       port,
			})
		}
	}

	c.JSON(http.StatusOK, sites)
}

// GetNginxSite gets a specific Nginx site configuration
func GetNginxSite(c *gin.Context) {
	name := c.Param("name")
	
	sitePath := filepath.Join(nginxSitesAvailable, name)
	
	content, err := os.ReadFile(sitePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	enabledPath := filepath.Join(nginxSitesEnabled, name)
	_, err = os.Stat(enabledPath)
	enabled := err == nil

	serverName, port := parseNginxConfig(string(content))

	site := NginxSite{
		Name:       name,
		Path:       sitePath,
		Enabled:    enabled,
		Content:    string(content),
		ServerName: serverName,
		Port:       port,
	}

	c.JSON(http.StatusOK, site)
}

// CreateNginxSite creates a new Nginx site
func CreateNginxSite(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required"`
		Content    string `json:"content"`
		ServerName string `json:"serverName"`
		Port       string `json:"port"`
		Root       string `json:"root"`
		ProxyPass  string `json:"proxyPass"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sitePath := filepath.Join(nginxSitesAvailable, req.Name)

	// Check if site already exists
	if _, err := os.Stat(sitePath); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Site already exists"})
		return
	}

	// Generate config if not provided
	content := req.Content
	if content == "" {
		content = generateNginxConfig(req.ServerName, req.Port, req.Root, req.ProxyPass)
	}

	// Write config file
	if err := os.WriteFile(sitePath, []byte(content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create site: " + err.Error()})
		return
	}

	// Test Nginx configuration
	if err := testNginxConfig(); err != nil {
		os.Remove(sitePath)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Nginx configuration: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Site created successfully"})
}

// UpdateNginxSite updates an Nginx site configuration
func UpdateNginxSite(c *gin.Context) {
	name := c.Param("name")

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sitePath := filepath.Join(nginxSitesAvailable, name)

	// Check if site exists
	if _, err := os.Stat(sitePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	// Backup current config
	backup, _ := os.ReadFile(sitePath)

	// Write new config
	if err := os.WriteFile(sitePath, []byte(req.Content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update site: " + err.Error()})
		return
	}

	// Test Nginx configuration
	if err := testNginxConfig(); err != nil {
		// Restore backup
		os.WriteFile(sitePath, backup, 0644)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Nginx configuration: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Site updated successfully"})
}

// DeleteNginxSite deletes an Nginx site
func DeleteNginxSite(c *gin.Context) {
	name := c.Param("name")

	sitePath := filepath.Join(nginxSitesAvailable, name)
	enabledPath := filepath.Join(nginxSitesEnabled, name)

	// Disable site first
	os.Remove(enabledPath)

	// Delete site
	if err := os.Remove(sitePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete site: " + err.Error()})
		return
	}

	// Reload Nginx
	reloadNginx()

	c.JSON(http.StatusOK, gin.H{"message": "Site deleted successfully"})
}

// EnableNginxSite enables an Nginx site
func EnableNginxSite(c *gin.Context) {
	name := c.Param("name")

	sitePath := filepath.Join(nginxSitesAvailable, name)
	enabledPath := filepath.Join(nginxSitesEnabled, name)

	// Check if site exists
	if _, err := os.Stat(sitePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	// Create symlink
	if err := os.Symlink(sitePath, enabledPath); err != nil {
		if !os.IsExist(err) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable site: " + err.Error()})
			return
		}
	}

	// Test and reload Nginx
	if err := testNginxConfig(); err != nil {
		os.Remove(enabledPath)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Nginx configuration: " + err.Error()})
		return
	}

	if err := reloadNginx(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload Nginx: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Site enabled successfully"})
}

// DisableNginxSite disables an Nginx site
func DisableNginxSite(c *gin.Context) {
	name := c.Param("name")

	enabledPath := filepath.Join(nginxSitesEnabled, name)

	// Remove symlink
	if err := os.Remove(enabledPath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Site is not enabled"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable site: " + err.Error()})
		return
	}

	// Reload Nginx
	if err := reloadNginx(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload Nginx: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Site disabled successfully"})
}

// TestNginxConfig tests the Nginx configuration
func TestNginxConfig(c *gin.Context) {
	if err := testNginxConfig(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "valid": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration is valid", "valid": true})
}

// ReloadNginx reloads Nginx
func ReloadNginx(c *gin.Context) {
	if err := reloadNginx(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Nginx reloaded successfully"})
}

// GetNginxStatus gets Nginx status
func GetNginxStatus(c *gin.Context) {
	cmd := exec.Command("systemctl", "status", "nginx")
	output, _ := cmd.CombinedOutput()

	running := strings.Contains(string(output), "active (running)")

	c.JSON(http.StatusOK, gin.H{
		"running": running,
		"status":  string(output),
	})
}

// Helper functions

func testNginxConfig() error {
	cmd := exec.Command("nginx", "-t")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func reloadNginx() error {
	cmd := exec.Command("systemctl", "reload", "nginx")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(output))
	}
	return nil
}

func parseNginxConfig(content string) (serverName, port string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "server_name") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				serverName = strings.TrimSuffix(parts[1], ";")
			}
		}
		if strings.HasPrefix(line, "listen") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				port = strings.TrimSuffix(parts[1], ";")
			}
		}
	}
	return
}

func generateNginxConfig(serverName, port, root, proxyPass string) string {
	if port == "" {
		port = "80"
	}

	config := fmt.Sprintf(`server {
    listen %s;
    server_name %s;

`, port, serverName)

	if root != "" {
		config += fmt.Sprintf(`    root %s;
    index index.html index.htm;

    location / {
        try_files $uri $uri/ =404;
    }
`, root)
	} else if proxyPass != "" {
		config += fmt.Sprintf(`    location / {
        proxy_pass %s;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
`, proxyPass)
	}

	config += "}\n"
	return config
}
