package routes

import (
	"github.com/Creative-genius001/Stacklo/services/transaction/api/handler"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, h *handler.Handler) {

	transactionR := router.Group("/api/transaction")
	{
		transactionR.POST("/create", h.CreateTransaction)
		transactionR.GET("/", h.GetAllTransactions)
		transactionR.GET("/:id", h.GetSingleTransaction)

	}

}
