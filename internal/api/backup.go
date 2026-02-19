package api

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Backup represents a backup
type Backup struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // file, database, system
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	Created     time.Time `json:"created"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
}

const backupDir = "/var/backups/serverpanel"

// ListBackups lists all backups
func ListBackups(c *gin.Context) {
	// Ensure backup directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		// If we can't create the directory (permission issue), return empty list
		c.JSON(http.StatusOK, []Backup{})
		return
	}

	files, err := os.ReadDir(backupDir)
	if err != nil {
		// If we can't read the directory, return empty list
		c.JSON(http.StatusOK, []Backup{})
		return
	}

	var backups []Backup
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		// Parse backup filename to extract metadata
		name := file.Name()
		backupType := "file"
		if strings.Contains(name, "database") {
			backupType = "database"
		} else if strings.Contains(name, "system") {
			backupType = "system"
		}

		backups = append(backups, Backup{
			ID:      name,
			Name:    name,
			Type:    backupType,
			Path:    filepath.Join(backupDir, name),
			Size:    info.Size(),
			Created: info.ModTime(),
			Status:  "completed",
		})
	}

	c.JSON(http.StatusOK, backups)
}

// CreateBackup creates a new backup
func CreateBackup(c *gin.Context) {
	var req struct {
		Type        string `json:"type" binding:"required"` // file, database
		Name        string `json:"name" binding:"required"`
		Source      string `json:"source" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure backup directory exists
	os.MkdirAll(backupDir, 0755)

	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s-%s.tar.gz", req.Type, req.Name, timestamp)
	backupPath := filepath.Join(backupDir, filename)

	var err error
	switch req.Type {
	case "file":
		err = createFileBackup(req.Source, backupPath)
	case "database":
		err = createDatabaseBackup(req.Source, backupPath)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backup: " + err.Error()})
		return
	}

	// Get file info
	info, _ := os.Stat(backupPath)
	
	backup := Backup{
		ID:          filename,
		Name:        filename,
		Type:        req.Type,
		Path:        backupPath,
		Size:        info.Size(),
		Created:     time.Now(),
		Description: req.Description,
		Status:      "completed",
	}

	c.JSON(http.StatusOK, backup)
}

// DownloadBackup downloads a backup file
func DownloadBackup(c *gin.Context) {
	id := c.Param("id")

	backupPath := filepath.Join(backupDir, id)

	// Security check
	if !strings.HasPrefix(backupPath, backupDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid backup ID"})
		return
	}

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
		return
	}

	c.FileAttachment(backupPath, id)
}

// DeleteBackup deletes a backup
func DeleteBackup(c *gin.Context) {
	id := c.Param("id")

	backupPath := filepath.Join(backupDir, id)

	// Security check
	if !strings.HasPrefix(backupPath, backupDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid backup ID"})
		return
	}

	if err := os.Remove(backupPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete backup: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup deleted successfully"})
}

// RestoreBackup restores from a backup
func RestoreBackup(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Destination string `json:"destination" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	backupPath := filepath.Join(backupDir, id)

	// Security check
	if !strings.HasPrefix(backupPath, backupDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid backup ID"})
		return
	}

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
		return
	}

	// Restore the backup
	if err := restoreFileBackup(backupPath, req.Destination); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore backup: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup restored successfully"})
}

// GetBackupStats gets backup statistics
func GetBackupStats(c *gin.Context) {
	files, err := os.ReadDir(backupDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"totalBackups": 0,
			"totalSize":    0,
		})
		return
	}

	var totalSize int64
	count := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		totalSize += info.Size()
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"totalBackups": count,
		"totalSize":    totalSize,
		"totalSizeMB":  fmt.Sprintf("%.2f", float64(totalSize)/(1024*1024)),
	})
}

// Helper functions

// createFileBackup creates a tar.gz backup of a file or directory
func createFileBackup(source, destination string) error {
	// Create output file
	outFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Walk the source directory
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return err
		}

		// Update header name to preserve directory structure
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// If it's a file, write its content
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(tarWriter, file); err != nil {
				return err
			}
		}

		return nil
	})
}

// restoreFileBackup restores a tar.gz backup
func restoreFileBackup(source, destination string) error {
	// Open the backup file
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create gzip reader
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzReader)

	// Extract files
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Construct full path
		target := filepath.Join(destination, header.Name)

		// Create directory if needed
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			continue
		}

		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		// Create file
		outFile, err := os.Create(target)
		if err != nil {
			return err
		}

		// Copy content
		if _, err := io.Copy(outFile, tarReader); err != nil {
			outFile.Close()
			return err
		}
		outFile.Close()
	}

	return nil
}

// createDatabaseBackup creates a database backup using mysqldump or pg_dump
func createDatabaseBackup(dbConfig, destination string) error {
	// Parse database config (format: type:host:port:user:password:database)
	parts := strings.Split(dbConfig, ":")
	if len(parts) != 6 {
		return fmt.Errorf("invalid database config format")
	}

	dbType := parts[0]
	host := parts[1]
	port := parts[2]
	user := parts[3]
	password := parts[4]
	database := parts[5]

	var cmd *exec.Cmd

	if dbType == "mysql" {
		// Use MYSQL_PWD environment variable instead of command-line argument
		cmd = exec.Command("mysqldump",
			"-h", host,
			"-P", port,
			"-u", user,
			database,
		)
		cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", password))
	} else if dbType == "postgresql" {
		cmd = exec.Command("pg_dump",
			"-h", host,
			"-p", port,
			"-U", user,
			database,
		)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", password))
	} else {
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Create output file
	outFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Pipe output to file
	cmd.Stdout = outFile
	
	return cmd.Run()
}
