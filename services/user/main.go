package main

import (
	"net/http"

	"time"

	"github.com/Creative-genius001/Stacklo/services/user/api/routes"
	"github.com/Creative-genius001/Stacklo/services/user/config"
	"github.com/Creative-genius001/Stacklo/services/user/db"
	"github.com/Creative-genius001/Stacklo/utils/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(gin.Recovery())
	// router.Use(limit.MaxAllowed(200))

	//connect to postgres DB
	db.InitDB()

	logger.Info("Connection to database url successful")

	//init routes
	routes.InitializeRoutes(router)

	//init config
	config.Init()

	//Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowOrigins = []string{"*"}
	router.Use(cors.New(corsConfig))

	//startup server
	PORT := config.Cfg.Port
	s := &http.Server{
		Addr:           ":" + PORT,
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.Info("Server is starting and running on port: ", PORT)
	if err := s.ListenAndServe(); err != nil {
		logger.Fatal("Failed to start server ", err)
	}

}
