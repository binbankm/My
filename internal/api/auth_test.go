package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/binbankm/My/internal/model"
	"github.com/gin-gonic/gin"
)

func setupTestDB(t *testing.T) {
	if err := model.InitDB(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
}

func TestLogin(t *testing.T) {
	setupTestDB(t)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
	}{
		{
			name:           "Valid credentials",
			username:       "admin",
			password:       "admin123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid username",
			username:       "nonexistent",
			password:       "admin123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid password",
			username:       "admin",
			password:       "wrongpassword",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty username",
			username:       "",
			password:       "admin123",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty password",
			username:       "admin",
			password:       "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Username with spaces trimmed",
			username:       "  admin  ",
			password:       "admin123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Case insensitive username",
			username:       "ADMIN",
			password:       "admin123",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			loginData := map[string]string{
				"username": tt.username,
				"password": tt.password,
			}
			jsonData, _ := json.Marshal(loginData)
			c.Request, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			Login(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Code == http.StatusOK {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to parse response: %v", err)
				}
				if _, ok := response["token"]; !ok {
					t.Error("Expected token in response")
				}
			}
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	setupTestDB(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate authenticated user
	c.Set("user_id", uint(1))
	c.Request, _ = http.NewRequest("GET", "/api/auth/info", nil)

	GetUserInfo(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if user, ok := response["user"].(map[string]interface{}); ok {
		if user["username"] != "admin" {
			t.Errorf("Expected username 'admin', got %v", user["username"])
		}
	} else {
		t.Error("Expected user in response")
	}
}

func TestLogout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/auth/logout", nil)

	Logout(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
