package e2e

import (
	"net/http"
	"testing"

	"foglio/v2/src/middlewares"
	"foglio/v2/src/routes"
	"foglio/v2/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	server *utils.TestServer
}

func (suite *E2ETestSuite) SetupSuite() {
	suite.server = utils.SetupTestServer()
	
	// Setup middlewares
	suite.server.Router.Use(middlewares.ErrorHandlerMiddleware())
	suite.server.Router.Use(middlewares.AuthMiddleware())
	
	// Setup routes
	prefix := "/api/v1"
	router := suite.server.Router.Group(prefix)
	
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.JobRoutes(router)
	routes.NotificationRoutes(router)
}

func (suite *E2ETestSuite) TearDownSuite() {
	suite.server.Cleanup()
}

func (suite *E2ETestSuite) TestHealthEndpoint() {
	w := utils.MakeRequest(suite.server.Router, "GET", "/api/v1/health", nil)
	
	response := utils.AssertJSONResponse(suite.T(), w, http.StatusOK)
	assert.Equal(suite.T(), "success", response["status"])
}

func (suite *E2ETestSuite) TestAuthFlow() {
	// Test user registration
	registerData := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}
	
	w := utils.MakeRequest(suite.server.Router, "POST", "/api/v1/auth/register", registerData)
	response := utils.AssertJSONResponse(suite.T(), w, http.StatusCreated)
	assert.Equal(suite.T(), "success", response["status"])
	
	// Test user login
	loginData := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}
	
	w = utils.MakeRequest(suite.server.Router, "POST", "/api/v1/auth/login", loginData)
	response = utils.AssertJSONResponse(suite.T(), w, http.StatusOK)
	assert.Equal(suite.T(), "success", response["status"])
	assert.Contains(suite.T(), response, "data")
	
	data := response["data"].(map[string]interface{})
	assert.Contains(suite.T(), data, "token")
}

func (suite *E2ETestSuite) TestUnauthorizedAccess() {
	// Test accessing protected endpoint without token
	w := utils.MakeRequest(suite.server.Router, "GET", "/api/v1/user/profile", nil)
	utils.AssertJSONResponse(suite.T(), w, http.StatusUnauthorized)
}

func (suite *E2ETestSuite) TestInvalidEndpoint() {
	w := utils.MakeRequest(suite.server.Router, "GET", "/api/v1/nonexistent", nil)
	utils.AssertJSONResponse(suite.T(), w, http.StatusNotFound)
}

func (suite *E2ETestSuite) TestJobEndpoints() {
	// First register and login to get token
	registerData := map[string]interface{}{
		"email":    "jobtest@example.com",
		"password": "password123",
		"name":     "Job Test User",
	}
	
	utils.MakeRequest(suite.server.Router, "POST", "/api/v1/auth/register", registerData)
	
	loginData := map[string]interface{}{
		"email":    "jobtest@example.com",
		"password": "password123",
	}
	
	w := utils.MakeRequest(suite.server.Router, "POST", "/api/v1/auth/login", loginData)
	response := utils.AssertJSONResponse(suite.T(), w, http.StatusOK)
	
	data := response["data"].(map[string]interface{})
	token := data["token"].(string)
	
	// Test getting jobs (should be empty initially)
	w = utils.MakeAuthenticatedRequest(suite.server.Router, "GET", "/api/v1/jobs", token, nil)
	utils.AssertJSONResponse(suite.T(), w, http.StatusOK)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}