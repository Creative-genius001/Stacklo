package routes

import (
	"github.com/Creative-genius001/Stacklo/services/user/api/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, a *handler.AuthHandler, u *handler.UserHandler) {

	user := router.Group("/api/user")
	{
		user.GET("/:id", u.GetUser)
		user.PUT("/:id", u.UpdateUser)
	}

	auth := router.Group("/api/auth")
	{
		auth.POST("/login", a.Login)
		auth.POST("/register", a.Register)
	}
}
