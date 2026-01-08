package e2e

import (
	"net/http"
	"testing"

	"foglio/v2/src/lib"
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
	// Auth is enforced at handler level for protected routes in tests

	// Setup routes
	prefix := "/api/v2"
	router := suite.server.Router.Group(prefix)

	routes.HealthRoutes(router)
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.JobRoutes(router)
	routes.NotificationRoutes(router)

	// Return JSON for unknown endpoints like the main app
	suite.server.Router.NoRoute(lib.GlobalNotFound())
}

func (suite *E2ETestSuite) TearDownSuite() {
	suite.server.Cleanup()
}

func (suite *E2ETestSuite) TestHealthEndpoint() {
	w := utils.MakeRequest(suite.server.Router, "GET", "/api/v2/health", nil)

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

	w := utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signup", registerData)
	response := utils.AssertJSONResponse(suite.T(), w, http.StatusCreated)
	assert.Equal(suite.T(), "success", response["status"])

	// Test user login
	loginData := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}

	w = utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signin", loginData)
	response = utils.AssertJSONResponse(suite.T(), w, http.StatusOK)
	assert.Equal(suite.T(), "success", response["status"])

	// Only check for data if the response is successful
	if response["status"] == "success" {
		assert.Contains(suite.T(), response, "data")
		if data, ok := response["data"].(map[string]interface{}); ok {
			assert.Contains(suite.T(), data, "token")
		}
	}
}

func (suite *E2ETestSuite) TestUnauthorizedAccess() {
	// Test accessing protected endpoint without token
	w := utils.MakeRequest(suite.server.Router, "GET", "/api/v2/user/profile", nil)
	utils.AssertJSONResponse(suite.T(), w, http.StatusUnauthorized)
}

func (suite *E2ETestSuite) TestInvalidEndpoint() {
	w := utils.MakeRequest(suite.server.Router, "GET", "/api/v2/nonexistent", nil)
	utils.AssertJSONResponse(suite.T(), w, http.StatusNotFound)
}

func (suite *E2ETestSuite) TestJobEndpoints() {
	// First register and login to get token
	registerData := map[string]interface{}{
		"email":    "jobtest@example.com",
		"password": "password123",
		"name":     "Job Test User",
	}

	utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signup", registerData)
	// We do not assert signup here because it might fail if user already exists (e.g. from previous run),
	// but we proceed to login which should work either way if the user exists.

	loginData := map[string]interface{}{
		"email":    "jobtest@example.com",
		"password": "password123",
	}

	w := utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signin", loginData)
	response := utils.AssertJSONResponse(suite.T(), w, http.StatusOK)
	assert.Equal(suite.T(), "success", response["status"])

	// Check if data exists before accessing it
	if response["status"] != "success" || response["data"] == nil {
		suite.T().FailNow()
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		suite.T().Fatalf("Response data is not a map: %v", response["data"])
	}
	token := data["token"].(string)

	// Test getting jobs (should be empty initially)
	w = utils.MakeAuthenticatedRequest(suite.server.Router, "GET", "/api/v2/jobs", token, nil)
	utils.AssertJSONResponse(suite.T(), w, http.StatusOK)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
