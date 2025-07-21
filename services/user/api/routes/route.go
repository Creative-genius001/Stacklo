package routes

import (
	"github.com/Creative-genius001/Stacklo/services/user/api/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, a *handler.AuthHandler, u *handler.UserHandler) {

	user := router.Group("/api/v1/user")
	{
		user.GET("/:id", u.GetUser)
		user.PUT("/:id", u.UpdateUser)
	}

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", a.Login)
		auth.POST("/register", a.Register)
		auth.POST("/verify-otp", a.VerifyOTP)
		auth.POST("/resend-otp", a.ResendOTP)
	}
}
