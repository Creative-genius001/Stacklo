package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Creative-genius001/Stacklo/api-gateway/internal/middlewares"
	r "github.com/Creative-genius001/Stacklo/api-gateway/internal/router"
	"github.com/Creative-genius001/Stacklo/api-gateway/internal/utils/logger"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(".env file not found", zap.Error(err))
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.InitLogger(appEnv)
	defer logger.Logger.Sync()

	// cfg := config.LoadConfig()
	// if len(cfg.Microservices) == 0 {
	// 	logger.Logger.Fatal("No microservices configured. Please set WALLET_SERVICE_URL, PAYMENT_SERVICE_URL, etc.")
	// }

	router := gin.New()

	//health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	limiter := middlewares.NewClientLimiter(rate.Every(time.Minute/5), 10)
	router.Use(limiter.Middleware())

	router.Use(middlewares.RequestLoggerMiddleware())
	router.Use(middlewares.ErrorRecoveryMiddleware())
	router.Use(middlewares.SecurityHeadersMiddleware(os.Getenv("EXPECTED_HOST")))
	r.SetupRoutes(router)

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.Logger.Info("Starting server", zap.String("port", os.Getenv("PORT")))
	if err := s.ListenAndServe(); err != nil {
		logger.Logger.Fatal("Server failed to start", zap.Error(err))
	}
}
