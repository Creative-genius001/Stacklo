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

	c.JSON(http.StatusCreated, gin.H{"data": banks})
}
