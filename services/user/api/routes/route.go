package routes

import (
	"github.com/Creative-genius001/Stacklo/services/user/api/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, a *handler.AuthHandler, u *handler.UserHandler) {

	userRouter := router.Group("/api/user")
	{
		userRouter.GET("/:id", u.GetUser)
	}

	authRouter := router.Group("/api/auth")
	{
		authRouter.POST("/login", a.Login)
		authRouter.POST("/register", a.Register)
	}
}
