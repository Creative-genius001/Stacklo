package routes

import (
	"github.com/Creative-genius001/Stacklo/services/wallet/api/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, h *handler.Handler) {

	walletRouter := router.Group("/api/wallet")
	{
		walletRouter.GET("/:id", h.GetAllWallets)
		walletRouter.GET("/fiat/:id", h.GetFiatWallet)
		walletRouter.POST("/fiat/create", h.CreateFiatWallet)
		walletRouter.POST("/crypto/create", h.CreateCryptoWallet)
		// walletRouter.POST("/deposit", handler.CreateWallet)
		// walletRouter.GET("/balance", handler.CreateWallet)
		// walletRouter.POST("/withdraw", handler.CreateWallet)
	}

}
