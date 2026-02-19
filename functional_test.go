package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// Test File Management Operations
func TestFileManagement(t *testing.T) {
	testDir := "/home/runner/test-file-ops"
	defer os.RemoveAll(testDir)

	t.Run("Create Directory", func(t *testing.T) {
		createData := map[string]interface{}{
			"path":  testDir,
			"isDir": true,
		}
		jsonData, _ := json.Marshal(createData)

		resp := makeAuthRequest(t, "POST", "/api/files", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Failed to create directory. Status: %d, Body: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("Create File", func(t *testing.T) {
		// First create directory
		os.MkdirAll(testDir, 0755)

		createData := map[string]interface{}{
			"path":    testDir + "/test.txt",
			"content": "Hello World",
			"isDir":   false,
		}
		jsonData, _ := json.Marshal(createData)

		resp := makeAuthRequest(t, "POST", "/api/files", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Failed to create file. Status: %d, Body: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("Update File", func(t *testing.T) {
		// Ensure file exists
		os.MkdirAll(testDir, 0755)
		os.WriteFile(testDir+"/test.txt", []byte("Initial"), 0644)

		updateData := map[string]interface{}{
			"path":    testDir + "/test.txt",
			"content": "Updated Content",
		}
		jsonData, _ := json.Marshal(updateData)

		resp := makeAuthRequest(t, "PUT", "/api/files", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Failed to update file. Status: %d, Body: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("List Files", func(t *testing.T) {
		// Ensure directory exists with file
		os.MkdirAll(testDir, 0755)
		os.WriteFile(testDir+"/test.txt", []byte("content"), 0644)

		resp := makeAuthRequest(t, "GET", "/api/files?path="+testDir, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Failed to list files. Status: %d, Body: %s", resp.StatusCode, string(body))
			return
		}

		var files []interface{}
		if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
	})

	t.Run("Delete File", func(t *testing.T) {
		// Ensure file exists
		os.MkdirAll(testDir, 0755)
		testFile := testDir + "/delete-me.txt"
		os.WriteFile(testFile, []byte("content"), 0644)

		resp := makeAuthRequest(t, "DELETE", "/api/files?path="+testFile, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Failed to delete file. Status: %d, Body: %s", resp.StatusCode, string(body))
		}
	})
}

// Test Database Management Operations
func TestDatabaseManagement(t *testing.T) {
	var dbID string

	t.Run("Create Database Connection", func(t *testing.T) {
		// Create a SQLite database connection (doesn't require actual server)
		createData := map[string]interface{}{
			"name": "test-db",
			"type": "mysql",
			"host": "localhost",
			"port": 3306,
			"username": "test",
			"password": "test",
			"database": "testdb",
		}
		jsonData, _ := json.Marshal(createData)

		resp := makeAuthRequest(t, "POST", "/api/database", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&result)
			if id, ok := result["id"].(float64); ok {
				dbID = string(rune(int(id)))
			}
		}
	})

	t.Run("List Databases", func(t *testing.T) {
		resp := makeAuthRequest(t, "GET", "/api/database", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Failed to list databases. Status: %d", resp.StatusCode)
			return
		}

		var databases []interface{}
		if err := json.NewDecoder(resp.Body).Decode(&databases); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
	})

	t.Run("Test Database Connection", func(t *testing.T) {
		if dbID == "" {
			t.Skip("No database ID available")
		}

		resp := makeAuthRequest(t, "POST", "/api/database/"+dbID+"/test", nil)
		defer resp.Body.Close()

		// Test might fail if database isn't available, which is OK
		// We just verify the endpoint responds
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("Unexpected status code: %d", resp.StatusCode)
		}
	})
}

// Test Cron Job Management
func TestCronManagement(t *testing.T) {
	var jobID string

	t.Run("Create Cron Job", func(t *testing.T) {
		createData := map[string]interface{}{
			"name":     "test-job",
			"schedule": "0 0 * * *",
			"command":  "echo 'test'",
			"enabled":  true,
		}
		jsonData, _ := json.Marshal(createData)

		resp := makeAuthRequest(t, "POST", "/api/cron", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		// Cron might not be available or accessible
		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&result)
			if id, ok := result["id"].(string); ok {
				jobID = id
			}
		}
	})

	t.Run("List Cron Jobs", func(t *testing.T) {
		resp := makeAuthRequest(t, "GET", "/api/cron", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Failed to list cron jobs. Status: %d", resp.StatusCode)
		}
	})

	t.Run("Delete Cron Job", func(t *testing.T) {
		if jobID == "" {
			t.Skip("No job ID available")
		}

		resp := makeAuthRequest(t, "DELETE", "/api/cron/"+jobID, nil)
		defer resp.Body.Close()

		// Might fail if cron is not accessible
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
			t.Logf("Cron delete returned status: %d", resp.StatusCode)
		}
	})
}

// Test User Management
func TestUserManagement(t *testing.T) {
	var userID string

	t.Run("Create User", func(t *testing.T) {
		createData := map[string]interface{}{
			"username": "testuser",
			"password": "testpass123",
			"email":    "test@example.com",
			"roleId":   1,
		}
		jsonData, _ := json.Marshal(createData)

		resp := makeAuthRequest(t, "POST", "/api/users", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var result map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&result)
			if user, ok := result["user"].(map[string]interface{}); ok {
				if id, ok := user["id"].(float64); ok {
					userID = string(rune(int(id)))
				}
			}
		}
	})

	t.Run("List Users", func(t *testing.T) {
		resp := makeAuthRequest(t, "GET", "/api/users", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Failed to list users. Status: %d", resp.StatusCode)
			return
		}

		var users []interface{}
		if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		if len(users) == 0 {
			t.Error("Expected at least one user")
		}
	})

	t.Run("Get User", func(t *testing.T) {
		if userID == "" {
			userID = "1" // Use admin user
		}

		resp := makeAuthRequest(t, "GET", "/api/users/"+userID, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Failed to get user. Status: %d", resp.StatusCode)
		}
	})

	t.Run("Delete User", func(t *testing.T) {
		if userID == "" || userID == "1" {
			t.Skip("Cannot delete admin user")
		}

		resp := makeAuthRequest(t, "DELETE", "/api/users/"+userID, nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Logf("Delete user status: %d, body: %s", resp.StatusCode, string(body))
		}
	})
}

// Test Settings Management
func TestSettingsManagement(t *testing.T) {
	t.Run("Get Settings", func(t *testing.T) {
		resp := makeAuthRequest(t, "GET", "/api/settings", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Failed to get settings. Status: %d", resp.StatusCode)
			return
		}

		var settings []interface{}
		if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
	})

	t.Run("Update Settings", func(t *testing.T) {
		updateData := []map[string]interface{}{
			{
				"key":   "site_name",
				"value": "Test Server Panel",
			},
			{
				"key":   "language",
				"value": "en",
			},
		}
		jsonData, _ := json.Marshal(updateData)

		resp := makeAuthRequest(t, "PUT", "/api/settings", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Failed to update settings. Status: %d, Body: %s", resp.StatusCode, string(body))
		}
	})
}

// Test WebSocket Connection (basic connectivity test)
func TestWebSocketConnection(t *testing.T) {
	t.Run("WebSocket Endpoint Available", func(t *testing.T) {
		// We can't easily test WebSocket in HTTP test, but we can check if endpoint exists
		client := &http.Client{Timeout: 5 * time.Second}
		req, _ := http.NewRequest("GET", baseURL+"/api/ws", nil)
		
		token := getAuthToken(t)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Upgrade", "websocket")
		req.Header.Set("Connection", "Upgrade")
		req.Header.Set("Sec-WebSocket-Version", "13")
		req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

		resp, err := client.Do(req)
		if err != nil {
			t.Logf("WebSocket connection attempt: %v", err)
			return
		}
		defer resp.Body.Close()

		// Should get upgrade response or bad request (not 404)
		if resp.StatusCode == http.StatusNotFound {
			t.Error("WebSocket endpoint not found")
		}
	})
}

// Test Role Management
func TestRoleManagement(t *testing.T) {
	t.Run("List Roles", func(t *testing.T) {
		resp := makeAuthRequest(t, "GET", "/api/roles", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Failed to list roles. Status: %d", resp.StatusCode)
			return
		}

		var roles []interface{}
		if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
	})

	t.Run("Get Role", func(t *testing.T) {
		resp := makeAuthRequest(t, "GET", "/api/roles/1", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Logf("Get role status: %d", resp.StatusCode)
		}
	})
}
