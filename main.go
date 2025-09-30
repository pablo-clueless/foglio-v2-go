package main

import (
	"fmt"
	"foglio/v2/src/config"
	"foglio/v2/src/database"
	"foglio/v2/src/docs"
	"foglio/v2/src/handlers"
	"foglio/v2/src/lib"
	"foglio/v2/src/middlewares"
	"foglio/v2/src/routes"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitializeEnvFile()
	config.InitializeConfig()

	err := database.InitializeDatabase()
	defer database.CloseDatabase()
	if err != nil {
		log.Fatal("Database error:", err)
	}

	lib.InitialiseJWT(string(config.AppConfig.JWTTokenSecret))
	corsConfig := cors.Config{
		AllowOrigins:     []string{config.AppConfig.ClientUrl},
		AllowMethods:     config.AppConfig.AllowMethods,
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}

	if config.AppConfig.IsDevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()

	docs.SetupSwagger(app)

	app.Use(gin.Logger())
	app.Use(cors.New(corsConfig))
	app.Use(middlewares.ErrorHandlerMiddleware())
	app.Use(middlewares.AuthMiddleware())
	app.Use(middlewares.RateLimiterMiddleware())
	app.Use(lib.ErrorHandler())

	app.MaxMultipartMemory = 10 << 20 // 10

	hub := lib.NewHub()
	go hub.Run()

	websocket := lib.NewWebSocketHandler(hub)

	prefix := config.AppConfig.Version
	router := app.Group(prefix)

	router.GET("/", func(ctx *gin.Context) {
		lib.Success(ctx, "", map[string]interface{}{})
	})
	router.GET("/favicon.ico", func(ctx *gin.Context) {
		handlers.SendFile(ctx, "/favicon.ico")
	})
	router.GET("/ws", websocket.HandleWebSocket)
	router.GET("/health", func(ctx *gin.Context) {
		lib.Success(ctx, "Foglio API is healthy", map[string]interface{}{
			"version": config.AppConfig.Version,
			"status":  200,
		})
	})
	router.GET("/swagger/*any", func(ctx *gin.Context) {})

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.JobRoutes(router)
	app.NoRoute(lib.GlobalNotFound())

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.AppConfig.Port),
		Handler:        app,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Printf("Server starting on port http://localhost:%s/%s", config.AppConfig.Port, config.AppConfig.Version)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}
