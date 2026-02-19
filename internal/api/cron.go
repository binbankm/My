package api

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CronJob represents a cron job
type CronJob struct {
	ID       string    `json:"id"`
	Schedule string    `json:"schedule"` // e.g., "0 2 * * *"
	Command  string    `json:"command"`
	User     string    `json:"user"`
	Enabled  bool      `json:"enabled"`
	Comment  string    `json:"comment"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// ListCronJobs lists all cron jobs for the current user
func ListCronJobs(c *gin.Context) {
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		// If crontab is empty, it returns an error
		if strings.Contains(string(output), "no crontab") {
			c.JSON(http.StatusOK, []CronJob{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list cron jobs: " + err.Error()})
		return
	}

	jobs := parseCronTab(string(output))
	c.JSON(http.StatusOK, jobs)
}

// CreateCronJob creates a new cron job
func CreateCronJob(c *gin.Context) {
	var req struct {
		Schedule string `json:"schedule" binding:"required"`
		Command  string `json:"command" binding:"required"`
		Comment  string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate cron schedule
	if !isValidCronSchedule(req.Schedule) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cron schedule format"})
		return
	}

	// Get current crontab
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()
	
	currentCrontab := ""
	if err == nil {
		currentCrontab = string(output)
	}

	// Add new cron job
	newJob := ""
	if req.Comment != "" {
		newJob = fmt.Sprintf("# %s\n", req.Comment)
	}
	newJob += fmt.Sprintf("%s %s\n", req.Schedule, req.Command)

	updatedCrontab := currentCrontab + newJob

	// Write updated crontab
	if err := writeCrontab(updatedCrontab); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cron job: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cron job created successfully"})
}

// UpdateCronJob updates an existing cron job
func UpdateCronJob(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Schedule string `json:"schedule" binding:"required"`
		Command  string `json:"command" binding:"required"`
		Comment  string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate cron schedule
	if !isValidCronSchedule(req.Schedule) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cron schedule format"})
		return
	}

	// Get current crontab
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read crontab: " + err.Error()})
		return
	}

	jobs := parseCronTab(string(output))
	
	// Find and update the job
	found := false
	for i, job := range jobs {
		if job.ID == id {
			jobs[i].Schedule = req.Schedule
			jobs[i].Command = req.Command
			jobs[i].Comment = req.Comment
			jobs[i].Modified = time.Now()
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cron job not found"})
		return
	}

	// Rebuild crontab
	newCrontab := buildCrontab(jobs)

	// Write updated crontab
	if err := writeCrontab(newCrontab); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cron job: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cron job updated successfully"})
}

// DeleteCronJob deletes a cron job
func DeleteCronJob(c *gin.Context) {
	id := c.Param("id")

	// Get current crontab
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read crontab: " + err.Error()})
		return
	}

	jobs := parseCronTab(string(output))
	
	// Find and remove the job
	found := false
	var newJobs []CronJob
	for _, job := range jobs {
		if job.ID == id {
			found = true
			continue
		}
		newJobs = append(newJobs, job)
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cron job not found"})
		return
	}

	// Rebuild crontab
	newCrontab := buildCrontab(newJobs)

	// Write updated crontab
	if err := writeCrontab(newCrontab); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cron job: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cron job deleted successfully"})
}

// GetCronJob gets a specific cron job
func GetCronJob(c *gin.Context) {
	id := c.Param("id")

	// Get current crontab
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read crontab: " + err.Error()})
		return
	}

	jobs := parseCronTab(string(output))
	
	for _, job := range jobs {
		if job.ID == id {
			c.JSON(http.StatusOK, job)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Cron job not found"})
}

// Helper functions

// parseCronTab parses crontab output into CronJob structs
func parseCronTab(crontab string) []CronJob {
	var jobs []CronJob
	scanner := bufio.NewScanner(strings.NewReader(crontab))
	
	var currentComment string
	lineNum := 0
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++
		
		// Skip empty lines
		if line == "" {
			continue
		}
		
		// Handle comments
		if strings.HasPrefix(line, "#") {
			currentComment = strings.TrimPrefix(line, "#")
			currentComment = strings.TrimSpace(currentComment)
			continue
		}
		
		// Parse cron job line
		parts := strings.Fields(line)
		if len(parts) < 6 {
			continue
		}
		
		schedule := strings.Join(parts[0:5], " ")
		command := strings.Join(parts[5:], " ")
		
		job := CronJob{
			ID:       fmt.Sprintf("cron-%d", lineNum),
			Schedule: schedule,
			Command:  command,
			Comment:  currentComment,
			Enabled:  true,
			Created:  time.Now(),
			Modified: time.Now(),
		}
		
		jobs = append(jobs, job)
		currentComment = ""
	}
	
	return jobs
}

// buildCrontab builds a crontab string from CronJob structs
func buildCrontab(jobs []CronJob) string {
	var lines []string
	
	for _, job := range jobs {
		if job.Comment != "" {
			lines = append(lines, "# "+job.Comment)
		}
		lines = append(lines, fmt.Sprintf("%s %s", job.Schedule, job.Command))
	}
	
	return strings.Join(lines, "\n") + "\n"
}

// writeCrontab writes the crontab
func writeCrontab(crontab string) error {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "crontab-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	
	// Write crontab to temp file
	if _, err := tmpFile.WriteString(crontab); err != nil {
		return err
	}
	tmpFile.Close()
	
	// Install the new crontab
	cmd := exec.Command("crontab", tmpFile.Name())
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %s", err.Error(), string(output))
	}
	
	return nil
}

// isValidCronSchedule validates a cron schedule format
func isValidCronSchedule(schedule string) bool {
	parts := strings.Fields(schedule)
	if len(parts) != 5 {
		return false
	}
	
	// Basic validation - Unix cron format
	for _, part := range parts {
		if part != "*" {
			// Check if it's a number, range, list, or step
			if !strings.ContainsAny(part, "0123456789,-/") {
				return false
			}
		}
	}
	
	return true
}
