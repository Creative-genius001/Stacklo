package routes

import (
	"user/api/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	authRouter := router.Group("/api/auth")
	{
		authRouter.GET("/:id", handlers.Login)
		authRouter.GET("/:id", handlers.Register)
	}
}
