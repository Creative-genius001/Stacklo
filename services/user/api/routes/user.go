package routes

import (
	"user/api/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userRouter := router.Group("/api/user")
	{
		userRouter.GET("/:id", handlers.GetUserData)
	}
}
