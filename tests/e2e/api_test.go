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
	suite.server.Router.Use(middlewares.ErrorHandlerMiddleware())

	prefix := "/api/v2"
	router := suite.server.Router.Group(prefix)
	routes.HealthRoutes(router)
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.JobRoutes(router)
	routes.NotificationRoutes(router)

	suite.server.Router.NoRoute(lib.GlobalNotFound())
}

func (suite *E2ETestSuite) TearDownSuite() {
	suite.server.Cleanup()
}

func (suite *E2ETestSuite) TestRootEndpoint() {
	w := utils.MakeRequest(suite.server.Router, "GET", "/", nil)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *E2ETestSuite) TestHealthEndpoint() {
	w := utils.MakeRequest(suite.server.Router, "GET", "/api/v2/health", nil)
	response := utils.AssertJSONResponse(suite.T(), w, http.StatusOK)
	assert.Equal(suite.T(), "success", response["status"])
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
