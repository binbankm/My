package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

// getDockerClient creates a new Docker client
func getDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// ListContainers lists all Docker containers
func ListContainers(c *gin.Context) {
	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Docker: " + err.Error()})
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list containers: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, containers)
}

// GetContainer gets details of a specific container
func GetContainer(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Docker: " + err.Error()})
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	containerJSON, err := cli.ContainerInspect(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Container not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, containerJSON)
}

// StartContainer starts a container
func StartContainer(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Docker: " + err.Error()})
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := cli.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start container: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container started successfully"})
}

// StopContainer stops a container
func StopContainer(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Docker: " + err.Error()})
		return
	}
	defer cli.Close()

	timeout := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout+10)*time.Second)
	defer cancel()

	if err := cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop container: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container stopped successfully"})
}

// RestartContainer restarts a container
func RestartContainer(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Docker: " + err.Error()})
		return
	}
	defer cli.Close()

	timeout := 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout+20)*time.Second)
	defer cancel()

	if err := cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &timeout}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restart container: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container restarted successfully"})
}

// DeleteContainer deletes a container
func DeleteContainer(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Docker: " + err.Error()})
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := cli.ContainerRemove(ctx, id, container.RemoveOptions{Force: true}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete container: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container deleted successfully"})
}

// ListImages lists all Docker images
func ListImages(c *gin.Context) {
	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Docker: " + err.Error()})
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	images, err := cli.ImageList(ctx, image.ListOptions{All: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list images: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, images)
}

// DeleteImage deletes a Docker image
func DeleteImage(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Docker: " + err.Error()})
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := cli.ImageRemove(ctx, id, image.RemoveOptions{Force: true}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}

// GetContainerLogs gets logs from a container
func GetContainerLogs(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Docker: " + err.Error()})
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "100",
	}

	logs, err := cli.ContainerLogs(ctx, id, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get logs: " + err.Error()})
		return
	}
	defer logs.Close()

	logContent, err := io.ReadAll(logs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read logs: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": string(logContent)})
}

// GetContainerStats gets real-time stats from a container
func GetContainerStats(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Docker: " + err.Error()})
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stats, err := cli.ContainerStats(ctx, id, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats: " + err.Error()})
		return
	}
	defer stats.Body.Close()

	statsData, err := io.ReadAll(stats.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read stats: " + err.Error()})
		return
	}

	// Parse and return the stats
	var statsJSON map[string]interface{}
	if err := json.Unmarshal(statsData, &statsJSON); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse stats: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, statsJSON)
}
