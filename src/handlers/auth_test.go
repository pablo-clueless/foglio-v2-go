package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/lib"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	if err := config.InitializeEnvFile(); err != nil {
		panic(err)
	}
	config.InitializeConfig()
	if err := database.InitializeDatabase(); err != nil {
		panic(err)
	}
	lib.InitialiseJWT(string(config.AppConfig.JWTTokenSecret))

	router := gin.New()
	return router
}

func TestRegisterHandler(t *testing.T) {
	router := setupTestRouter()
	defer func() {
		if err := database.CloseDatabase(); err != nil {
			// Log error but continue
		}
	}()

	router.POST("/register", NewAuthHandler().CreateUser())

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid registration",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
				"name":     "Test User",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing email",
			payload: map[string]interface{}{
				"password": "password123",
				"name":     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid email format",
			payload: map[string]interface{}{
				"email":    "invalid-email",
				"password": "password123",
				"name":     "Test User",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonPayload, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonPayload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
