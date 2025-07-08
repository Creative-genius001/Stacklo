package handler

import (
	"net/http"
	"strings"

	services "github.com/Creative-genius001/Stacklo/services/wallet/api/service"
	"github.com/Creative-genius001/Stacklo/services/wallet/types"
	errors "github.com/Creative-genius001/Stacklo/services/wallet/utils/error"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service services.Service
}

func NewHandler(s services.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetWallet(c *gin.Context) {

	walletIDStr := strings.TrimSpace(c.Param("id"))
	if walletIDStr == "" {
		logger.Logger.Warn("Invalid request data")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	wallet, err := h.service.GetWallet(c.Request.Context(), walletIDStr)
	if err != nil {
		appErr, _ := err.(*errors.CustomError)
		if appErr.Type == errors.TypeNotFound {
			logger.Logger.Info("Wallet not found")
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "error": appErr.Message})
			return
		}
		logger.Logger.Error("Error retrieving wallet", zap.Error(err))
		c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "error": errors.TypeInternal})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    wallet,
	})
}

func (h *Handler) CreateWallet(c *gin.Context) {

	var customerReq types.CreateCustomerRequest

	if err := c.ShouldBindJSON(&customerReq); err != nil {
		logger.Logger.Warn("Invalid input data")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	wallet, err := services.CreateWalletPaystack(customerReq)
	if err != nil {
		appErr, ok := err.(*errors.CustomError)
		if !ok {
			logger.Logger.Error("Fuck it didnt assert")
			c.JSON(errors.GetHTTPStatus(errors.TypeForbidden), gin.H{"status": "error", "message": errors.TypeForbidden})
			return
		}

		if appErr.Type == errors.TypeInternal || appErr.Type == errors.TypeExternal {
			logger.Logger.Error("Service error during wallet fetch", zap.Error(appErr))
		} else {
			logger.Logger.Info("Client-facing error during wallet fetch", zap.Error(appErr))
		}
		c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Message})
		return
	}

	w, err := h.service.CreateWallet(c.Request.Context(), *wallet)
	if err != nil {
		appErr, ok := err.(*errors.CustomError)
		if !ok {
			logger.Logger.Error("Fuck it didnt assert")
			c.JSON(errors.GetHTTPStatus(errors.TypeForbidden), gin.H{"status": "error", "message": errors.TypeForbidden})
			return
		}
		if appErr.Type == errors.TypeInternal {
			logger.Logger.Error("Service error: Could not create wallet", zap.Error(appErr))
		} else {
			logger.Logger.Info("Client-facing error during wallet fetch", zap.Error(appErr))
		}
		c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"wallet": gin.H{
			"id":                     w.ID,
			"active":                 w.Active,
			"balance":                w.Balance,
			"currency":               w.Currency,
			"virtual_account_name":   w.VirtualAccountName,
			"virtual_account_number": w.VirtualAccountNumber,
			"virtual_bank_name":      w.VirtualBankName,
		},
	})
}
