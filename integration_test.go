//go:build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	baseURL  = "http://localhost:8888"
	username = "admin"
	password = "admin123"
)

var authToken string

func TestMain(m *testing.M) {
	// Wait for server to be ready
	fmt.Println("Waiting for server to be ready...")
	for i := 0; i < 30; i++ {
		resp, err := http.Get(baseURL + "/api/auth/login")
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(1 * time.Second)
	}
	
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func getAuthToken(t *testing.T) string {
	if authToken != "" {
		return authToken
	}

	loginData := map[string]string{
		"username": username,
		"password": password,
	}
	jsonData, _ := json.Marshal(loginData)

	resp, err := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	token, ok := result["token"].(string)
	if !ok {
		t.Fatal("No token in login response")
	}

	authToken = token
	return authToken
}

func makeAuthRequest(t *testing.T, method, path string, body io.Reader) *http.Response {
	req, err := http.NewRequest(method, baseURL+path, body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	token := getAuthToken(t)
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	return resp
}

func TestAuthLogin(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		password   string
		wantStatus int
	}{
		{"Valid credentials", "admin", "admin123", http.StatusOK},
		{"Invalid username", "nonexistent", "admin123", http.StatusUnauthorized},
		{"Invalid password", "admin", "wrong", http.StatusUnauthorized},
		{"Empty username", "", "admin123", http.StatusBadRequest},
		{"Empty password", "admin", "", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginData := map[string]string{
				"username": tt.username,
				"password": tt.password,
			}
			jsonData, _ := json.Marshal(loginData)

			resp, err := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("Expected status %d, got %d. Body: %s", tt.wantStatus, resp.StatusCode, string(body))
			}
		})
	}
}

func TestSystemInfo(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/system/info", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var info map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	requiredFields := []string{"hostname", "os", "platform", "uptime", "cpuCores"}
	for _, field := range requiredFields {
		if _, ok := info[field]; !ok {
			t.Errorf("Missing field: %s", field)
		}
	}
}

func TestSystemStats(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/system/stats", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	requiredFields := []string{"cpu", "memory", "disk", "network"}
	for _, field := range requiredFields {
		if _, ok := stats[field]; !ok {
			t.Errorf("Missing field: %s", field)
		}
	}
}

func TestDockerContainers(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/docker/containers", nil)
	defer resp.Body.Close()

	// Docker might not be available, so we accept both success and service unavailable
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusServiceUnavailable {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Unexpected status %d. Body: %s", resp.StatusCode, string(body))
	}
}

func TestDockerImages(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/docker/images", nil)
	defer resp.Body.Close()

	// Docker might not be available
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusServiceUnavailable {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Unexpected status %d. Body: %s", resp.StatusCode, string(body))
	}
}

func TestListFiles(t *testing.T) {
	// Use /home which is the default allowed path
	resp := makeAuthRequest(t, "GET", "/api/files?path=/home", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
		return
	}

	var files []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
}

func TestDatabases(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/database", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}
}

func TestCronJobs(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/cron", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}
}

func TestLogFiles(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/logs/files", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}
}

func TestNginxStatus(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/nginx/status", nil)
	defer resp.Body.Close()

	// Nginx might not be installed, so we accept both success and error responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Unexpected status %d. Body: %s", resp.StatusCode, string(body))
	}
}

func TestBackups(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/backup", nil)
	defer resp.Body.Close()

	// Backup directory might not be accessible or not exist, both are acceptable
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Unexpected status %d. Body: %s", resp.StatusCode, string(body))
		return
	}

	// If successful, verify response format
	if resp.StatusCode == http.StatusOK {
		var backups []interface{}
		if err := json.NewDecoder(resp.Body).Decode(&backups); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
	}
}

func TestUsers(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/users", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}

	var users []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(users) == 0 {
		t.Error("Expected at least one user (admin)")
	}
}

func TestRoles(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/roles", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}
}

func TestPermissions(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/permissions", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}
}

func TestSettings(t *testing.T) {
	resp := makeAuthRequest(t, "GET", "/api/settings", nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status 200, got %d. Body: %s", resp.StatusCode, string(body))
	}
}
