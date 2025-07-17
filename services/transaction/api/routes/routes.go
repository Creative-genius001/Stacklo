package routes

import (
	"github.com/Creative-genius001/Stacklo/services/transaction/api/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, h *handler.Handler) {

	transaction := router.Group("/api/v1/transaction")
	{
		transaction.POST("/create", h.CreateTransaction)
		transaction.GET("/", h.GetAllTransactions)
		transaction.GET("/:id", h.GetSingleTransaction)
		transaction.GET("/filter/:user_id", h.GetFilteredTransactions)

	}

}
