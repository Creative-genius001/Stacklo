package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Creative-genius001/Stacklo/pkg/paystack"
	"github.com/Creative-genius001/Stacklo/services/wallet/api/handler"
	"github.com/Creative-genius001/Stacklo/services/wallet/api/routes"
	"github.com/Creative-genius001/Stacklo/services/wallet/api/service"
	"github.com/Creative-genius001/Stacklo/services/wallet/config"
	"github.com/Creative-genius001/Stacklo/services/wallet/middlewares"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"go.uber.org/zap"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tinrab/retry"
)

func main() {
	config.Init()
	c := config.Cfg

	if err := godotenv.Load("../../.env"); err != nil {
		logger.Logger.Fatal(".env file not found", zap.Error(err))
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	payApi := os.Getenv("PAYSTACK_TEST_KEY")
	payUrl := os.Getenv("PAYSTACK_BASE_URL")

	logger.InitLogger(appEnv)
	defer logger.Logger.Sync()

	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	PORT := c.Port

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestLoggerMiddleware())
	r.Use(middlewares.ErrorRecoveryMiddleware())

	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	var re service.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		re, err = service.NewPostgresRepository(c.DBUrl)
		if err != nil {
			logger.Logger.Fatal("Failed to connect to database", zap.Error(err))
		}
		return
	})
	defer re.Close()
	ps := paystack.NewPaystackClient(payApi, payUrl)
	svc := service.NewService(re, ps)
	h := handler.NewHandler(svc)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "404 not found"})
	})
	routes.InitializeRoutes(r, h)

	s := &http.Server{
		Addr:           ":" + PORT,
		Handler:        r,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.Logger.Info("Starting server", zap.String("port", PORT))
	if err := s.ListenAndServe(); err != nil {
		logger.Logger.Fatal("Server failed to start", zap.Error(err))
	}

}
