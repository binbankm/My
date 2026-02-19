package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetSystemInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/system/info", nil)

	GetSystemInfo(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	// Verify essential fields exist
	requiredFields := []string{"hostname", "os", "platform", "kernelVersion", "uptime", "cpuCores"}
	for _, field := range requiredFields {
		if _, ok := response[field]; !ok {
			t.Errorf("Expected field '%s' in response", field)
		}
	}
}

func TestGetSystemStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/system/stats", nil)

	GetSystemStats(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	// Verify essential stats fields exist
	requiredFields := []string{"cpu", "memory", "disk", "network"}
	for _, field := range requiredFields {
		if _, ok := response[field]; !ok {
			t.Errorf("Expected field '%s' in response", field)
		}
	}

	// Verify CPU stats (array of percentages)
	if _, ok := response["cpu"].([]interface{}); !ok {
		t.Error("Expected CPU to be an array")
	}

	// Verify Memory stats
	if memory, ok := response["memory"].(map[string]interface{}); ok {
		if _, ok := memory["total"]; !ok {
			t.Error("Expected 'total' in memory stats")
		}
		if _, ok := memory["used"]; !ok {
			t.Error("Expected 'used' in memory stats")
		}
	} else {
		t.Error("Expected memory stats in response")
	}
}
