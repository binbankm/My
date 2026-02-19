package api

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Source    string `json:"source"`
}

// LogFile represents a log file
type LogFile struct {
	Path     string    `json:"path"`
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	ModTime  time.Time `json:"modTime"`
	IsSystem bool      `json:"isSystem"`
}

// Common log directories
var logDirs = []string{
	"/var/log",
	"/var/log/nginx",
	"/var/log/apache2",
	"/var/log/mysql",
	"/var/log/postgresql",
}

// ListLogFiles lists available log files
func ListLogFiles(c *gin.Context) {
	var logFiles []LogFile

	for _, dir := range logDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			// Skip directories and non-log files
			if info.IsDir() {
				return nil
			}

			// Check if it's a log file
			if strings.HasSuffix(info.Name(), ".log") || 
			   strings.HasSuffix(info.Name(), ".log.1") ||
			   strings.Contains(info.Name(), ".log.") {
				logFiles = append(logFiles, LogFile{
					Path:     path,
					Name:     info.Name(),
					Size:     info.Size(),
					ModTime:  info.ModTime(),
					IsSystem: true,
				})
			}

			return nil
		})
	}

	c.JSON(http.StatusOK, logFiles)
}

// ReadLogFile reads content from a log file
func ReadLogFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path is required"})
		return
	}

	// Get optional parameters
	lines := c.DefaultQuery("lines", "100")
	filter := c.Query("filter")
	tail := c.DefaultQuery("tail", "true") == "true"

	// Security: Ensure path is in allowed directories
	allowed := false
	for _, dir := range logDirs {
		if strings.HasPrefix(path, dir) {
			allowed = true
			break
		}
	}

	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this file"})
		return
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log file not found"})
		return
	}

	var content string

	if tail {
		// Use tail command for better performance on large files
		cmd := exec.Command("tail", "-n", lines, path)
		output, cmdErr := cmd.Output()
		if cmdErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read log file: " + cmdErr.Error()})
			return
		}
		content = string(output)
	} else {
		// Read first N lines
		cmd := exec.Command("head", "-n", lines, path)
		output, cmdErr := cmd.Output()
		if cmdErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read log file: " + cmdErr.Error()})
			return
		}
		content = string(output)
	}

	// Apply filter if provided
	if filter != "" {
		var filteredLines []string
		scanner := bufio.NewScanner(strings.NewReader(content))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(strings.ToLower(line), strings.ToLower(filter)) {
				filteredLines = append(filteredLines, line)
			}
		}
		content = strings.Join(filteredLines, "\n")
	}

	c.JSON(http.StatusOK, gin.H{
		"path":    path,
		"content": content,
		"lines":   strings.Count(content, "\n") + 1,
	})
}

// SearchLogs searches across multiple log files
func SearchLogs(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query is required"})
		return
	}

	dir := c.DefaultQuery("dir", "/var/log")
	
	// Security check
	allowed := false
	for _, allowedDir := range logDirs {
		if dir == allowedDir || strings.HasPrefix(dir, allowedDir+"/") {
			allowed = true
			break
		}
	}

	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this directory"})
		return
	}

	// Use grep to search
	cmd := exec.Command("grep", "-r", "-i", "-n", "--include=*.log", query, dir)
	output, err := cmd.Output()

	results := []map[string]string{}
	
	if err == nil {
		scanner := bufio.NewScanner(strings.NewReader(string(output)))
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ":", 3)
			if len(parts) >= 3 {
				results = append(results, map[string]string{
					"file":    parts[0],
					"line":    parts[1],
					"content": parts[2],
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   query,
		"results": results,
		"count":   len(results),
	})
}

// GetSystemLogs gets system logs using journalctl
func GetSystemLogs(c *gin.Context) {
	lines := c.DefaultQuery("lines", "100")
	unit := c.Query("unit")
	since := c.Query("since")

	args := []string{"-n", lines, "--no-pager"}
	
	if unit != "" {
		args = append(args, "-u", unit)
	}
	
	if since != "" {
		args = append(args, "--since", since)
	}

	cmd := exec.Command("journalctl", args...)
	output, err := cmd.Output()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get system logs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"content": string(output),
		"lines":   strings.Count(string(output), "\n"),
	})
}

// DownloadLogFile downloads a log file
func DownloadLogFile(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path is required"})
		return
	}

	// Security: Ensure path is in allowed directories
	allowed := false
	for _, dir := range logDirs {
		if strings.HasPrefix(path, dir) {
			allowed = true
			break
		}
	}

	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this file"})
		return
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log file not found"})
		return
	}

	c.FileAttachment(path, filepath.Base(path))
}

// ClearLogFile clears a log file (truncates it)
func ClearLogFile(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Security: Ensure path is in allowed directories
	allowed := false
	for _, dir := range logDirs {
		if strings.HasPrefix(req.Path, dir) {
			allowed = true
			break
		}
	}

	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this file"})
		return
	}

	// Truncate the file
	file, err := os.OpenFile(req.Path, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear log file: " + err.Error()})
		return
	}
	defer file.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Log file cleared successfully"})
}

// GetLogStats gets statistics about log files
func GetLogStats(c *gin.Context) {
	var totalSize int64
	var fileCount int

	for _, dir := range logDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if !info.IsDir() && (strings.HasSuffix(info.Name(), ".log") || 
			   strings.Contains(info.Name(), ".log.")) {
				totalSize += info.Size()
				fileCount++
			}

			return nil
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"totalSize":  totalSize,
		"fileCount":  fileCount,
		"totalSizeMB": fmt.Sprintf("%.2f", float64(totalSize)/(1024*1024)),
	})
}
