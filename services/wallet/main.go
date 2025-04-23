package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Creative-genius001/Stacklo/services/wallet/api/routes"
	"github.com/Creative-genius001/Stacklo/services/wallet/config"
	"github.com/Creative-genius001/Stacklo/services/wallet/db"
	"github.com/Creative-genius001/Stacklo/utils/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	PORT := os.Getenv("PORT")

	router := gin.New()
	router.Use(gin.Recovery())
	// router.Use(limit.MaxAllowed(200))

	//initiallize postgres DB
	db.InitDB()

	//initiallize config
	config.Init()

	logger.Info("Connection to database url successful")

	//init routes
	routes.InitializeRoutes(router)

	//Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowOrigins = []string{"*"}
	router.Use(cors.New(corsConfig))

	//startup server
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
