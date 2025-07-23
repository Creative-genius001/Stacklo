package routes

import (
	"github.com/Creative-genius001/Stacklo/services/wallet/api/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, h *handler.Handler) {

	wallet := router.Group("/api/v1/wallet")
	{
		wallet.GET("/:id", h.GetAllWallets)
		wallet.GET("/fiat/:id", h.GetFiatWallet)
		wallet.POST("/fiat/create", h.CreateFiatWallet)
		wallet.POST("/crypto/create", h.CreateCryptoWallet)
		// walletRouter.POST("/deposit", handler.CreateWallet)
		// walletRouter.GET("/balance", handler.CreateWallet)
		// walletRouter.POST("/withdraw", handler.CreateWallet)
	}

}
