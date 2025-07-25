package main

import (
	"net/http"
	"os"
	"time"

	//natsclient "github.com/Creative-genius001/Stacklo/pkg/natsClient"
	"github.com/Creative-genius001/Stacklo/services/payment/api/handlers"
	"github.com/Creative-genius001/Stacklo/services/payment/api/routes"
	"github.com/Creative-genius001/Stacklo/services/payment/api/services"
	"github.com/Creative-genius001/Stacklo/services/payment/config"
	"github.com/Creative-genius001/Stacklo/services/payment/middlewares"
	bn "github.com/Creative-genius001/Stacklo/services/payment/pkg/binance"
	"github.com/Creative-genius001/Stacklo/services/payment/pkg/paystack"
	"github.com/Creative-genius001/Stacklo/services/payment/utils/logger"

	// "github.com/adshao/go-binance/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	//init config
	config.Init()
	c := config.Cfg

	if err := godotenv.Load("../../.env"); err != nil {
		logger.Logger.Fatal("No .env file found or failed to load", zap.Error(err))
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

	expectedHost := "localhost:" + c.Port

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.RequestLoggerMiddleware())
	r.Use(middlewares.ErrorRecoveryMiddleware())
	r.Use(func(c *gin.Context) {
		if c.Request.Host != expectedHost {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	})
	// router.Use(limit.MaxAllowed(200))
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"error": "404 not found"})
	})

	//Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))

	// bCli := binance.NewClient(c.BinanceAPIKey, c.BinanceSecretKey)
	binanceCli := bn.NewBinanceClient(c, logger.Logger)
	paystackCli := paystack.NewPaystackClient(payApi, payUrl, logger.Logger)

	paymentSvc := services.NewPaymentService(binanceCli, paystackCli)
	paymentHdlr := handlers.NewpaymentService(paymentSvc)
	//initialize NATS
	// natsclient.InitNATS()
	// defer natsclient.Close()

	//init routes
	routes.InitializeRoutes(r, *paymentHdlr)

	//startup server
	PORT := c.Port
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
