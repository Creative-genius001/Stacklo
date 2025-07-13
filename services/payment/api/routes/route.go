package routes

import (
	"github.com/Creative-genius001/Stacklo/services/payment/api/handlers"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, h handlers.PaymentService) {

	paymentRouter := router.Group("/api/payment")
	{
		paymentRouter.GET("/bank-list", h.GetBankList)
		paymentRouter.GET("/account-details", h.ResolveAccountNumber)
		paymentRouter.POST("/transfer/otp", h.GetOTP)
		// paymentRouter.POST("/otp/retry", handlers.RetryOtp)
		paymentRouter.POST("/transfer", h.Transfer)
		paymentRouter.GET("/ping", h.Ping)
		paymentRouter.GET("/convert", h.Convert)
	}

}
