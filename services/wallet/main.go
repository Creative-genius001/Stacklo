package wallet

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Creative-genius001/stacklo/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	router := gin.New()
	router.Use(gin.Recovery())
	// router.Use(limit.MaxAllowed(200))

	//initialise DB

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
