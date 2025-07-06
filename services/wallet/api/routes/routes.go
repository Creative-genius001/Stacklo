package routes

import (
	"github.com/Creative-genius001/Stacklo/services/wallet/api/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, h *handler.Handler) {

	walletRouter := router.Group("/api/wallet")
	{
		walletRouter.GET("/:id", h.GetWallet)
		walletRouter.POST("/create", h.CreateWallet)
		// walletRouter.POST("/deposit", handler.CreateWallet)
		// walletRouter.GET("/balance", handler.CreateWallet)
		// walletRouter.POST("/withdraw", handler.CreateWallet)
	}

}
