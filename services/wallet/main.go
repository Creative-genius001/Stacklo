package main

import (
	"net/http"
	"time"

	"github.com/Creative-genius001/Stacklo/services/wallet/api/handler"
	"github.com/Creative-genius001/Stacklo/services/wallet/api/routes"
	"github.com/Creative-genius001/Stacklo/services/wallet/api/service"
	"github.com/Creative-genius001/Stacklo/services/wallet/config"
	"github.com/Creative-genius001/go-logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tinrab/retry"
)

func main() {
	config.Init()

	if err := godotenv.Load("../../.env"); err != nil {
		logger.Fatal("No .env file found or failed to load")
	}

	PORT := config.Cfg.Port

	expectedHost := "localhost:" + config.Cfg.Port

	r := gin.Default()
	r.Use(gin.Recovery())
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
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowOrigins = []string{"*"}
	r.Use(cors.New(corsConfig))
	// router.Use(limit.MaxAllowed(200))

	var re service.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		re, err = service.NewPostgresRepository(config.Cfg.DBUrl)
		if err != nil {
			logger.Error(err)
		}
		return
	})
	defer re.Close()

	svc := service.NewService(re)
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
	logger.Info("Server is starting and running on port: ", PORT)
	if err := s.ListenAndServe(); err != nil {
		logger.Error("Failed to start server ", err, nil)
	}

}
