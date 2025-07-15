package handler

import (
	"net/http"
	"strings"

	services "github.com/Creative-genius001/Stacklo/services/transaction/api/service"
	"github.com/Creative-genius001/Stacklo/services/transaction/model"
	"github.com/Creative-genius001/Stacklo/services/transaction/utils"
	errors "github.com/Creative-genius001/Stacklo/services/transaction/utils/error"
	"github.com/Creative-genius001/Stacklo/services/transaction/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service services.Service
}

func NewHandler(s services.Service) *Handler {
	return &Handler{service: s}
}

func (s *Handler) GetSingleTransaction(c *gin.Context) {
	transactionID := strings.TrimSpace(c.Param("id"))
	if transactionID == "" {
		logger.Logger.Warn("Wallet ID is an empty string")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	isValid := utils.IsValidUUID(transactionID)
	if isValid == false {
		logger.Logger.Warn("Wallet ID is not a valid ID")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	transaction, err := s.service.GetSingleTransaction(c.Request.Context(), transactionID)
	if err != nil {
		logger.Logger.Error("unable to retrive transaction", zap.Error(err))
		c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "error": errors.TypeInternal})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transaction,
	})
}

func (s *Handler) GetAllTransactions(c *gin.Context) {
	userID := strings.TrimSpace(c.Query("id"))
	if userID == "" {
		logger.Logger.Warn("Wallet ID is an empty string")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	isValid := utils.IsValidUUID(userID)
	if isValid == false {
		logger.Logger.Warn("Wallet ID is not a valid ID")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	transactions, err := s.service.GetAllTransactions(c.Request.Context(), userID)
	if err != nil {
		logger.Logger.Error("unable to retrive transaction", zap.Error(err))
		c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "error": errors.TypeInternal})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transactions,
	})
}

func (s *Handler) CreateTransaction(c *gin.Context) {
	var transaction model.Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		logger.Logger.Warn("Invalid input data")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	err := s.service.CreateTransaction(c.Request.Context(), transaction)
	if err != nil {
		logger.Logger.Error("unable to insert transaction in db", zap.Error(err))
		c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "error": errors.TypeInternal})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "successfully created",
	})
}
