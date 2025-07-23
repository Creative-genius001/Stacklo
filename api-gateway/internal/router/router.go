package router

import (
	"github.com/Creative-genius001/Stacklo/api-gateway/internal/handler"
	"github.com/Creative-genius001/Stacklo/api-gateway/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	transaction := router.Group("/api/v1/transaction")
	transaction.Use(middlewares.JWTAuth())
	{
		transaction.GET("/:id", handler.ProxyTransactionByID)
		transaction.POST("/create", handler.ProxyCreateTransaction)
		transaction.GET("/", handler.ProxyGetAllTransactions)
		transaction.GET("/filter/:user_id", handler.ProxyGetFilteredTransactions)
	}

	wallet := router.Group("/api/v1/wallet")
	wallet.Use(middlewares.JWTAuth())
	{
		wallet.GET("/:id", handler.ProxyGetAllWallets)
		wallet.GET("/fiat/:id", handler.ProxyGetFiatWallet)
		wallet.POST("/fiat/create", handler.ProxyCreateFiatWallet)
		wallet.POST("/crypto/create", handler.ProxyCreateCryptoWallet)
	}

	payment := router.Group("/api/v1/payment")
	payment.Use(middlewares.JWTAuth())
	{
		payment.GET("/bank-list", handler.ProxyGetBankList)
		payment.GET("/account-details", handler.ProxyResolveAccountNumber)
		payment.POST("/transfer/otp", handler.ProxyGetOTP)
		payment.POST("/transfer", handler.ProxyTransfer)
		payment.GET("/ping", handler.ProxyPing)
		payment.GET("/convert", handler.ProxyConvert)
	}

	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", handler.ProxyLogin)
		auth.POST("/register", handler.ProxyRegister)
		auth.POST("/verify-otp", handler.ProxyVerifyOTP)
		auth.POST("/resend-otp", handler.ProxyResendOTP)
	}

	user := router.Group("/api/v1/user")
	user.Use(middlewares.JWTAuth())
	{
		user.GET("/:id", handler.ProxyGetUser)
		user.PUT("/:id", handler.ProxyUpdateUser)
	}
}
