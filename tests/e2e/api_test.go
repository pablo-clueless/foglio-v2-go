package e2e

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"foglio/v2/src/database"
	"foglio/v2/src/lib"
	"foglio/v2/src/middlewares"
	"foglio/v2/src/models"
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

func (suite *E2ETestSuite) TestHealthEndpoint() {
	w := utils.MakeRequest(suite.server.Router, "GET", "/api/v2/health", nil)
	response := utils.AssertJSONResponse(suite.T(), w, http.StatusOK)
	assert.Equal(suite.T(), "success", response["status"])
}

func (suite *E2ETestSuite) TestCompleteAuthFlow() {
	t := suite.T()
	db := database.GetDatabase()

	testEmail := fmt.Sprintf("authtest_%d@example.com", time.Now().Unix())
	testUsername := fmt.Sprintf("authuser_%d", time.Now().Unix())
	testPassword := "SecurePassword123!"

	var userID string
	var userOTP string

	defer func() {
		if userID != "" {
			db.Exec("DELETE FROM users WHERE id = ?", userID)
		}
	}()

	// Step 1: User Signup
	t.Run("1. User Signup", func(t *testing.T) {
		registerData := map[string]interface{}{
			"email":    testEmail,
			"password": testPassword,
			"name":     "Auth Test User",
			"username": testUsername,
		}

		w := utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signup", registerData)
		response := utils.AssertJSONResponse(t, w, http.StatusOK)

		assert.Equal(t, "success", response["status"])
		assert.Contains(t, response, "data")

		data := response["data"].(map[string]interface{})
		userID = data["id"].(string)

		assert.NotEmpty(t, userID)
		assert.Equal(t, testEmail, data["email"])
	})

	// Step 2: Get OTP from database (simulating email)
	t.Run("2. Retrieve OTP", func(t *testing.T) {
		var user models.User
		err := db.Where("email = ?", testEmail).First(&user).Error
		assert.NoError(t, err, "User should exist in database")

		userOTP = user.Otp
		assert.NotEmpty(t, userOTP, "OTP should be generated")
		assert.False(t, user.Verified, "User should not be verified yet")
	})

	// Step 3: Verify OTP
	t.Run("3. Verify User with OTP", func(t *testing.T) {
		verifyData := map[string]interface{}{
			"otp": userOTP,
		}

		w := utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/verification", verifyData)
		response := utils.AssertJSONResponse(t, w, http.StatusOK)

		assert.Equal(t, "success", response["status"])
		assert.Contains(t, response, "data")

		data := response["data"].(map[string]interface{})
		assert.True(t, data["verified"].(bool), "User should be verified after OTP verification")
	})

	// Step 4: Verify user is marked as verified in database
	t.Run("4. Confirm User Verified in Database", func(t *testing.T) {
		var user models.User
		err := db.Where("email = ?", testEmail).First(&user).Error
		assert.NoError(t, err)
		assert.True(t, user.Verified, "User should be verified in database")
	})

	// Step 5: Sign in with verified user
	t.Run("5. User Signin After Verification", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    testEmail,
			"password": testPassword,
		}

		w := utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signin", loginData)
		response := utils.AssertJSONResponse(t, w, http.StatusOK)

		assert.Equal(t, "success", response["status"])
		assert.Contains(t, response, "data")

		data := response["data"].(map[string]interface{})
		assert.Contains(t, data, "token")
		assert.NotEmpty(t, data["token"], "Token should be present")
	})

	// Step 6: Test signin fails for unverified user
	t.Run("6. Unverified User Cannot Signin", func(t *testing.T) {
		unverifiedEmail := fmt.Sprintf("unverified_%d@example.com", time.Now().Unix())
		unverifiedUsername := fmt.Sprintf("unverified_%d", time.Now().Unix())

		// Create unverified user
		registerData := map[string]interface{}{
			"email":    unverifiedEmail,
			"password": testPassword,
			"name":     "Unverified User",
			"username": unverifiedUsername,
		}

		utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signup", registerData)

		// Try to sign in without verification
		loginData := map[string]interface{}{
			"email":    unverifiedEmail,
			"password": testPassword,
		}

		w := utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signin", loginData)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Unverified user should not be able to sign in")

		// Cleanup
		db.Exec("DELETE FROM users WHERE email = ?", unverifiedEmail)
	})
}

func (suite *E2ETestSuite) TestAuthFlowInvalidCredentials() {
	t := suite.T()

	// Test with non-existent user
	t.Run("Signin with non-existent user", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "password123",
		}

		w := utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signin", loginData)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// Test with wrong password
	t.Run("Signin with wrong password", func(t *testing.T) {
		db := database.GetDatabase()
		testEmail := fmt.Sprintf("wrongpass_%d@example.com", time.Now().Unix())
		testUsername := fmt.Sprintf("wrongpass_%d", time.Now().Unix())

		// Create and verify user
		registerData := map[string]interface{}{
			"email":    testEmail,
			"password": "CorrectPassword123!",
			"name":     "Wrong Pass Test",
			"username": testUsername,
		}

		w := utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signup", registerData)
		utils.AssertJSONResponse(t, w, http.StatusOK)

		// Get OTP and verify
		var user models.User
		db.Where("email = ?", testEmail).First(&user)

		verifyData := map[string]interface{}{
			"otp": user.Otp,
		}
		utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/verification", verifyData)

		// Try to sign in with wrong password
		loginData := map[string]interface{}{
			"email":    testEmail,
			"password": "WrongPassword123!",
		}

		w = utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signin", loginData)
		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should fail with wrong password")

		// Cleanup
		db.Exec("DELETE FROM users WHERE email = ?", testEmail)
	})
}

func (suite *E2ETestSuite) TestUnauthorizedAccess() {
	w := utils.MakeRequest(suite.server.Router, "GET", "/api/v2/user/profile", nil)
	utils.AssertJSONResponse(suite.T(), w, http.StatusUnauthorized)
}

func (suite *E2ETestSuite) TestInvalidEndpoint() {
	w := utils.MakeRequest(suite.server.Router, "GET", "/api/v2/nonexistent", nil)
	utils.AssertJSONResponse(suite.T(), w, http.StatusNotFound)
}

func (suite *E2ETestSuite) TestJobEndpointsWithAuth() {
	t := suite.T()
	db := database.GetDatabase()

	testEmail := fmt.Sprintf("jobtest_%d@example.com", time.Now().Unix())
	testUsername := fmt.Sprintf("jobtest_%d", time.Now().Unix())
	testPassword := "JobTestPassword123!"

	// Register user
	registerData := map[string]interface{}{
		"email":    testEmail,
		"password": testPassword,
		"name":     "Job Test User",
		"username": testUsername,
	}
	utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signup", registerData)

	// Get OTP and verify
	var user models.User
	db.Where("email = ?", testEmail).First(&user)

	verifyData := map[string]interface{}{
		"otp": user.Otp,
	}
	utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/verification", verifyData)

	// Login to get token
	loginData := map[string]interface{}{
		"email":    testEmail,
		"password": testPassword,
	}

	w := utils.MakeRequest(suite.server.Router, "POST", "/api/v2/auth/signin", loginData)
	response := utils.AssertJSONResponse(t, w, http.StatusOK)

	data := response["data"].(map[string]interface{})
	token := data["token"].(string)

	// Test getting jobs with valid token
	w = utils.MakeAuthenticatedRequest(suite.server.Router, "GET", "/api/v2/jobs", token, nil)
	utils.AssertJSONResponse(t, w, http.StatusOK)

	// Cleanup
	db.Exec("DELETE FROM users WHERE email = ?", testEmail)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
