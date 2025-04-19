package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Creative-genius001/Stacklo/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	router := gin.New()
	router.Use(gin.Recovery())
	// router.Use(limit.MaxAllowed(200))

	//connect to postgres DB
	DB_URL := os.Getenv("DB_URL")
	dsn := DB_URL
	_, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}

	logger.Info("Connection to database url successful")

	//init routes
	//routes.InitializeRoutes(router)

	// Configure CORS
	//corsConfig := cors.DefaultConfig()
	//corsConfig.AddAllowHeaders("Authorization")
	//corsConfig.AllowOrigins = []string{"*"}
	//router.Use(cors.New(corsConfig))

	//startup server
	PORT := os.Getenv("PORT")
	s := &http.Server{
		Addr:           ":" + PORT,
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.Info("Server is starting and running on port: ", PORT)
	if s.ListenAndServe(); err != nil {
		logger.Error("Failed to start server ", err, nil)
	}

}
