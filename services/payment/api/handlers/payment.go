package handlers

import (
	"net/http"

	"github.com/Creative-genius001/Stacklo/services/payment/api/services"
	"github.com/Creative-genius001/Stacklo/services/payment/utils"
	"github.com/gin-gonic/gin"
)

func GetBankList(c *gin.Context) {
	banks, err := services.GetBankList()
	if err != nil {
		res := utils.NewError(http.StatusInternalServerError, err.Error())
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": banks})
}

func ResolveAccountNumber(c *gin.Context) {

	accountNumber := c.Query("account_number")
	bankCode := c.Query("bank_code")
	if accountNumber == "" && bankCode == "" {
		res := utils.NewError(http.StatusInternalServerError, "account_number and bank_code parameters are required")
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
		return
	}

	resp, err := services.ResolveAccountNumber(accountNumber, bankCode)
	if err != nil {
		res := utils.NewError(http.StatusInternalServerError, err.Error())
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}
