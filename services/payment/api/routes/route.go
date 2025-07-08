package routes

import (
	"github.com/Creative-genius001/Stacklo/services/payment/api/handlers"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {

	paymentRouter := router.Group("/api/payment")
	{
		paymentRouter.GET("/bank-list", handlers.GetBankList)
		// paymentRouter.GET("/account-details", handlers.ResolveAccountNumber)
		// paymentRouter.POST("/otp/get", handlers.RequestOTP)
		// paymentRouter.POST("/otp/retry", handlers.RetryOtp)
		// paymentRouter.POST("/transfer", handlers.Transfer)
	}

}
