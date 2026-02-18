package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	IsDir   bool   `json:"isDir"`
	ModTime string `json:"modTime"`
}

// getAllowedBasePath returns the allowed base path for file operations
func getAllowedBasePath() string {
	basePath := os.Getenv("FILE_MANAGER_BASE_PATH")
	if basePath == "" {
		basePath = "/home" // Default to /home directory
	}
	return basePath
}

// validatePath ensures the path is within allowed directory
func validatePath(requestPath string) (string, error) {
	basePath := getAllowedBasePath()
	
	// Clean the paths
	cleanPath := filepath.Clean(requestPath)
	cleanBase := filepath.Clean(basePath)
	
	// Make absolute
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", err
	}
	
	// Check if path is within base path
	if !strings.HasPrefix(absPath, cleanBase) {
		return "", fmt.Errorf("access denied: path outside allowed directory")
	}
	
	return absPath, nil
}

// ListFiles lists files in a directory
func ListFiles(c *gin.Context) {
	dir := c.Query("path")
	if dir == "" {
		dir = getAllowedBasePath()
	}

	validPath, err := validatePath(dir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}

	files, err := os.ReadDir(validPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot read directory"})
		return
	}

	var fileList []FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		fileList = append(fileList, FileInfo{
			Name:    file.Name(),
			Path:    filepath.Join(validPath, file.Name()),
			Size:    info.Size(),
			IsDir:   file.IsDir(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, fileList)
}

// CreateFile creates a new file or directory
func CreateFile(c *gin.Context) {
	var req struct {
		Path  string `json:"path" binding:"required"`
		IsDir bool   `json:"isDir"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validPath, err := validatePath(req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}

	if req.IsDir {
		if err := os.MkdirAll(validPath, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		file, err := os.Create(validPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		file.Close()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created successfully"})
}

// UpdateFile updates file content
func UpdateFile(c *gin.Context) {
	var req struct {
		Path    string `json:"path" binding:"required"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validPath, err := validatePath(req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}

	if err := os.WriteFile(validPath, []byte(req.Content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

// DeleteFile deletes a file or directory
func DeleteFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path is required"})
		return
	}

	validPath, err := validatePath(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}

	// Additional check: don't allow deletion of base path itself
	if validPath == getAllowedBasePath() {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete base directory"})
		return
	}

	if err := os.RemoveAll(validPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

// DownloadFile downloads a file
func DownloadFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path is required"})
		return
	}

	validPath, err := validatePath(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}

	c.FileAttachment(validPath, filepath.Base(validPath))
}

// UploadFile handles file upload
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	path := c.PostForm("path")
	if path == "" {
		path = getAllowedBasePath()
	}

	validPath, err := validatePath(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid path"})
		return
	}

	dst := filepath.Join(validPath, file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
