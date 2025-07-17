package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Creative-genius001/Stacklo/api-gateway/internal/middlewares"
	"github.com/Creative-genius001/Stacklo/api-gateway/internal/proxy"
	"github.com/Creative-genius001/Stacklo/api-gateway/internal/utils/logger"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	// "go.uber.org/zap"
)

func main() {

	if err := godotenv.Load("../../.env"); err != nil {
		// logger.Logger.Fatal(".env file not found", zap.Error(err))
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}
	// logger.InitLogger(appEnv)
	defer logger.Logger.Sync()

	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// cfg := config.LoadConfig()
	// if len(cfg.Microservices) == 0 {
	// 	logger.Logger.Fatal("No microservices configured. Please set WALLET_SERVICE_URL, PAYMENT_SERVICE_URL, etc.")
	// }

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
		// logger.Logger.Warn("REDIS_URL not set, using default for rate limiter", zap.String("default_url", redisURL))
	}

	router := gin.New()

	//health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.Use(middlewares.RequestLoggerMiddleware())
	router.Use(middlewares.ErrorRecoveryMiddleware())
	// router.Use(middlewares.IPRateLimiter(redisURL))
	router.Use(middlewares.SecurityHeadersMiddleware(os.Getenv("EXPECTED_HOST")))

	protected := router.Group("/api/v1")
	// protected.Use(middleware.Auth())
	proxy.SetupRoutes(protected)

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// logger.Logger.Info("Starting server", zap.String("port", os.Getenv("PORT")))
	if err := s.ListenAndServe(); err != nil {
		// logger.Logger.Fatal("Server failed to start", zap.Error(err))
	}
}
