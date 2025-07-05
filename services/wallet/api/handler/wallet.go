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
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
		return
	}

	// wallet, err := h.service.GetWallet(c.Request.Context(), id)
	// if err != nil {
	// 	res := utils.NewError(http.StatusInternalServerError, "Error retrieving wallet")
	// 	c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    "wallet",
	})
}

func CreateWallet(c *gin.Context) {
	var customerReq types.CreateCustomerRequest

	if err := c.ShouldBindJSON(&customerReq); err != nil {
		res := utils.NewError(http.StatusBadRequest, "Invalid input data")
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
		return
	}

	customer, err := services.CreateCustomer(customerReq)
	if err != nil {
		res := utils.NewError(http.StatusInternalServerError, "Error creating wallet")
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
		return
	}

	createWalletReq := types.CreateDVAWalletRequest{
		FirstName:     customerReq.FirstName,
		LastName:      customerReq.LastName,
		Email:         customerReq.Email,
		Phone:         customerReq.Phone,
		Customer:      customer.Data.ID,
		Country:       "NG",
		PreferredBank: "wema-bank",
	}

	wallet, err := services.CreateDVAWallet(&createWalletReq)
	if err != nil {
		res := utils.NewError(http.StatusInternalServerError, "Error creating wallet")
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status_code": http.StatusCreated,
		"message":     "wallet created successfully",
		"wallet":      wallet,
	})
}
