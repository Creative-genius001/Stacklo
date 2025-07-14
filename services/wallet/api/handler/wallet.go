package handler

import (
	"net/http"
	"strings"

	"github.com/Creative-genius001/Stacklo/pkg/paystack"
	services "github.com/Creative-genius001/Stacklo/services/wallet/api/service"
	"github.com/Creative-genius001/Stacklo/services/wallet/model"
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

func (h *Handler) GetAllWallets(c *gin.Context) {

	walletIDStr := strings.TrimSpace(c.Param("id"))
	if walletIDStr == "" {
		logger.Logger.Warn("Wallet ID is an empty string")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	wallet, err := h.service.GetAllWallets(c.Request.Context(), walletIDStr)
	if err != nil {
		appErr, ok := err.(*errors.CustomError)
		if !ok {
			c.JSON(errors.GetHTTPStatus(errors.TypeForbidden), gin.H{"status": "error", "message": errors.TypeForbidden})
			return
		}
		if appErr.Type == errors.TypeNotFound {
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "error": appErr.Message})
			return
		}
		logger.Logger.Error("Error retrieving wallet", zap.Error(appErr))
		c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "error": errors.TypeInternal})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    wallet,
	})
}

func (h *Handler) GetFiatWallet(c *gin.Context) {

	walletIDStr := strings.TrimSpace(c.Param("id"))
	if walletIDStr == "" {
		logger.Logger.Warn("Wallet ID is an empty string")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	wallet, err := h.service.GetFiatWallet(c.Request.Context(), walletIDStr)
	if err != nil {
		appErr, ok := err.(*errors.CustomError)
		if !ok {
			c.JSON(errors.GetHTTPStatus(errors.TypeForbidden), gin.H{"status": "error", "message": errors.TypeForbidden})
			return
		}
		if appErr.Type == errors.TypeNotFound {
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "error": appErr.Message})
			return
		}
		logger.Logger.Error("Error retrieving wallet", zap.Error(appErr))
		c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "error": errors.TypeInternal})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    wallet,
	})
}

func (h *Handler) CreateFiatWallet(c *gin.Context) {

	var customerReq paystack.CreateCustomerRequest

	if err := c.ShouldBindJSON(&customerReq); err != nil {
		logger.Logger.Warn("Invalid input data")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	wallet, err := h.service.CreateWalletPaystack(customerReq)
	if err != nil {
		appErr, ok := err.(*errors.CustomError)
		if !ok {
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
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

	w, err := h.service.CreateFiatWallet(c.Request.Context(), *wallet)
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

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "fiat wallet successfully created",
		"data": gin.H{
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

func (h *Handler) CreateCryptoWallet(c *gin.Context) {
	var wallet model.Wallet

	if err := c.ShouldBindJSON(&wallet); err != nil {
		logger.Logger.Warn("Invalid input data")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	err := h.service.CreateCryptoWallet(c.Request.Context(), wallet)
	if err != nil {
		appErr, ok := err.(*errors.CustomError)
		if !ok {
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

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Crypto wallet succesfully created",
	})
}
