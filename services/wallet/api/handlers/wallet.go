package handlers

import (
	"net/http"

	"github.com/Creative-genius001/Stacklo/services/wallet/api/services"
	"github.com/Creative-genius001/Stacklo/services/wallet/types"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils"
	"github.com/gin-gonic/gin"
)

func GetWalletDetails(c *gin.Context) {

}

func CreateWallet(c *gin.Context) {
	var customerReq types.CreateCustomerRequest

	if err := c.ShouldBindJSON(&customerReq); err != nil {
		err := utils.NewError(http.StatusBadRequest, "Invalid input data")
		c.AbortWithStatusJSON(err.StatusCode, err)
		return
	}

	customer, err := services.CreateCustomer(customerReq)
	if err != nil {
		utils.JSONErrorResp(c, http.StatusInternalServerError, "Error creating wallet")
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
		utils.JSONErrorResp(c, http.StatusInternalServerError, "Error creating wallet")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status_code": http.StatusCreated,
		"message":     "wallet created successfully",
		"wallet":      wallet,
	})
}
