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
	if err := config.InitializeEnvFile(); err != nil {
		log.Fatal("Failed to initialize env file:", err)
	}
	config.InitializeConfig()

	err := database.InitializeDatabase()
	defer func() {
		if err = database.CloseDatabase(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
	if err != nil {
		log.Fatal("Database error:", err)
	}

	lib.InitialiseJWT(string(config.AppConfig.JWTTokenSecret))

	app := gin.Default()

	docs.SetupSwagger(app)

	app.Use(gin.Logger())

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	app.Use(cors.New(corsConfig))

	app.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset, Content-Length, Accept-Encoding, X-CSRF-Token")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "43200") // 12 hours
		c.Status(http.StatusNoContent)
	})

	app.Use(middlewares.ErrorHandlerMiddleware())
	app.Use(middlewares.AuthMiddleware())
	app.Use(middlewares.RateLimiterMiddleware())
	app.Use(lib.ErrorHandler())

	app.MaxMultipartMemory = 10 << 20 // 10 MB

	hub := lib.NewHub()
	go hub.Run()

	websocket := lib.NewWebSocketHandler(hub)

	app.GET("/", func(ctx *gin.Context) {
		lib.Success(ctx, "Welcome to Foglio API", map[string]interface{}{})
	})

	app.GET("/favicon.ico", func(ctx *gin.Context) {
		handlers.SendFile(ctx, "/favicon.ico")
	})

	prefix := config.AppConfig.Version
	router := app.Group(prefix)

	router.GET("/", func(ctx *gin.Context) {
		lib.Success(ctx, "Foglio API v"+config.AppConfig.Version+" is running", map[string]interface{}{})
	})
	router.GET("/ws", websocket.HandleWebSocket)
	router.GET("/ws/stats", websocket.GetStats)
	router.POST("/ws/send-notification", websocket.SendNotification)
	router.POST("/ws/broadcast", websocket.Broadcast)
	router.GET("/health", func(ctx *gin.Context) {
		lib.Success(ctx, "Foglio API is healthy", map[string]interface{}{
			"version": config.AppConfig.Version,
			"status":  200,
		})
	})
	router.GET("/swagger/*any", func(ctx *gin.Context) {
		lib.Success(ctx, "Swagger documentation endpoint", map[string]interface{}{})
	})

	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	routes.JobRoutes(router)
	routes.SelfRoutes(router)
	routes.TestingRoutes(router)
	routes.NotificationRoutes(router)
	routes.SubscriptionRoutes(router)
	routes.PaystackRoutes(router)
	routes.DomainRoutes(router)
	routes.PortfolioRoutes(router)
	routes.AnalyticsRoutes(router)
	app.NoRoute(lib.GlobalNotFound())

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.AppConfig.Port),
		Handler:        app,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	if config.AppConfig.IsDevMode {
		log.Printf("Server starting on port http://localhost:%s/%s", config.AppConfig.Port, config.AppConfig.Version)
		log.Printf("Swagger docs at http://localhost:%s/swagger/index.html", config.AppConfig.Port)
		log.Printf("CORS enabled for origins: %v", []string{"*"})
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}
