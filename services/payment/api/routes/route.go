package routes

import (
	"github.com/Creative-genius001/Stacklo/services/payment/api/handlers"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, h handlers.PaymentService) {

	payment := router.Group("/api/v1/payment")
	{
		payment.GET("/bank-list", h.GetBankList)
		payment.GET("/account-details", h.ResolveAccountNumber)
		payment.POST("/transfer/otp", h.GetOTP)
		// payment.POST("/otp/retry", handlers.RetryOtp)
		payment.POST("/transfer", h.Transfer)
		payment.GET("/ping", h.Ping)
		payment.GET("/convert", h.Convert)
	}

}
