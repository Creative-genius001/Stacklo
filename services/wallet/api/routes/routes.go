package routes

import (
	"github.com/Creative-genius001/Stacklo/services/wallet/api/handlers"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {

	walletRouter := router.Group("/api/wallet")
	{
		walletRouter.GET("/:id", handlers.GetWalletDetails)
		walletRouter.POST("/create", handlers.CreateWallet)
		walletRouter.POST("/deposit", handlers.CreateWallet)
		walletRouter.GET("/balance", handlers.CreateWallet)
		walletRouter.POST("/withdraw", handlers.CreateWallet)
	}

}
