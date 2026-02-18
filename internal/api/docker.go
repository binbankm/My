package api

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func getDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// ListContainers returns all Docker containers
func ListContainers(c *gin.Context) {
	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker not available: " + err.Error()})
		return
	}
	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, containers)
}

// GetContainer returns details of a specific container
func GetContainer(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker not available"})
		return
	}
	defer cli.Close()

	containerJSON, err := cli.ContainerInspect(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Container not found"})
		return
	}

	c.JSON(http.StatusOK, containerJSON)
}

// StartContainer starts a Docker container
func StartContainer(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker not available"})
		return
	}
	defer cli.Close()

	if err := cli.ContainerStart(context.Background(), id, container.StartOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container started"})
}

// StopContainer stops a Docker container
func StopContainer(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker not available"})
		return
	}
	defer cli.Close()

	timeout := 10
	if err := cli.ContainerStop(context.Background(), id, container.StopOptions{Timeout: &timeout}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container stopped"})
}

// RestartContainer restarts a Docker container
func RestartContainer(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker not available"})
		return
	}
	defer cli.Close()

	timeout := 10
	if err := cli.ContainerRestart(context.Background(), id, container.StopOptions{Timeout: &timeout}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container restarted"})
}

// DeleteContainer removes a Docker container
func DeleteContainer(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker not available"})
		return
	}
	defer cli.Close()

	if err := cli.ContainerRemove(context.Background(), id, container.RemoveOptions{Force: true}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container removed"})
}

// ListImages returns all Docker images
func ListImages(c *gin.Context) {
	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker not available"})
		return
	}
	defer cli.Close()

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, images)
}

// DeleteImage removes a Docker image
func DeleteImage(c *gin.Context) {
	id := c.Param("id")

	cli, err := getDockerClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Docker not available"})
		return
	}
	defer cli.Close()

	_, err = cli.ImageRemove(context.Background(), id, types.ImageRemoveOptions{Force: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image removed"})
}
