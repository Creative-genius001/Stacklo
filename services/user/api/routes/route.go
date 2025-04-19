package routes

import (
	"github.com/Creative-genius001/Stacklo/api/handlers"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {

	userRouter := router.Group("/api/user")
	{
		userRouter.GET("/:id", handlers.GetUserData)
	}

	authRouter := router.Group("/api/auth")
	{
		authRouter.GET("/:id", handlers.Login)
		authRouter.GET("/:id", handlers.Register)
	}
}
