package handler

import (
	"net/http"

	services "github.com/Creative-genius001/Stacklo/services/wallet/api/service"
	"github.com/Creative-genius001/Stacklo/services/wallet/types"

	"github.com/Creative-genius001/Stacklo/services/wallet/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Service
}

func NewHandler(s services.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetWallet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		res := utils.NewError(http.StatusBadRequest, "Invalid request data")
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"success": false, "error": res.Error})
		return
	}

	wallet, err := h.service.GetWallet(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "wallet not found" {
			res := utils.NewError(http.StatusInternalServerError, "Wallet not found")
			c.AbortWithStatusJSON(res.StatusCode, gin.H{"success": false, "error": res.Error})
			return
		}
		res := utils.NewError(http.StatusInternalServerError, "Error retrieving wallet")
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"success": false, "error": res.Error})
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
		res := utils.NewError(http.StatusBadRequest, "Invalid input data")
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"success": false, "error": res.Error})
		return
	}

	wallet, err := services.CreateWalletPaystack(customerReq)

	w, err := h.service.CreateWallet(c.Request.Context(), *wallet)
	if err != nil {
		res := utils.NewError(http.StatusInternalServerError, "Error crreating wallet")
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"success": false, "error": res.Error})
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
