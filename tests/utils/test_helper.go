package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/lib"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestServer struct {
	Router *gin.Engine
}

func SetupTestServer() *TestServer {
	gin.SetMode(gin.TestMode)

	if err := config.InitializeEnvFile(); err != nil {
		panic(fmt.Sprintf("Failed to initialize env file: %v", err))
	}
	config.InitializeConfig()

	err := database.InitializeDatabase()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize test database: %v", err))
	}

	lib.InitialiseJWT(string(config.AppConfig.JWTTokenSecret))

	router := gin.New()

	return &TestServer{
		Router: router,
	}
}

func (ts *TestServer) Cleanup() {
	if err := database.CloseDatabase(); err != nil {
		// Log error but continue cleanup
	}
}

func MakeRequest(router *gin.Engine, method, url string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func MakeAuthenticatedRequest(router *gin.Engine, method, url, token string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

func AssertJSONResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int) map[string]interface{} {
	assert.Equal(t, expectedStatus, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	return response
}
