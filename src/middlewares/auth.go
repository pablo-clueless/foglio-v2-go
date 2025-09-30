package middlewares

import (
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/lib"
	"foglio/v2/src/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	bearerPrefix = "Bearer "
)

func isOpenRoute(path, method string) bool {
	for _, openRoute := range config.AppConfig.NonAuthRoutes {
		if matchRoute(openRoute.Endpoint, path) && openRoute.Method == method {
			return true
		}
	}
	return false
}

func matchRoute(pattern, path string) bool {
	if pattern == path {
		return true
	}

	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i := range patternParts {
		if strings.HasPrefix(patternParts[i], ":") {
			continue
		}
		if patternParts[i] != pathParts[i] {
			return false
		}
	}

	return true
}

func extractBearerToken(authHeader string) (string, bool) {
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", false
	}
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	trimmedToken := strings.TrimSpace(token)
	return trimmedToken, trimmedToken != ""
}

func AuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService(database.GetDatabase())

	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		log.Println("here is your request:", path, method)

		if isOpenRoute(path, method) {
			ctx.Next()
			return
		}

		authHeader := ctx.Request.Header.Get("Authorization")
		token, ok := extractBearerToken(authHeader)
		if !ok {
			ctx.Error(lib.NewApiErrror("No auth token found", http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		claims, err := lib.ValidateToken(token)
		if err != nil {
			ctx.Error(lib.NewApiErrror("Invalid auth token", http.StatusUnauthorized))
			ctx.Abort()
			return
		}

		user, err := authService.FindUserById(claims.UserId.String())
		if err != nil {
			ctx.Error(lib.NewApiErrror("User not found", http.StatusNotFound))
			ctx.Abort()
			return
		}

		ctx.Set(config.AppConfig.CurrentUser, user)
		ctx.Set(config.AppConfig.CurrentUserId, user.ID.String())
		ctx.Next()
	}
}
