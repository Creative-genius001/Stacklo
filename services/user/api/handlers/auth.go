package handlers

import (
	"net/http"

	"github.com/Creative-genius001/Stacklo/services/user/api/services"
	"github.com/Creative-genius001/Stacklo/services/user/types"
	"github.com/Creative-genius001/Stacklo/services/user/utils"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login route"})
}

func Register(c *gin.Context) {
	var RegForm types.RegisterType

	if err := c.ShouldBind(&RegForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	//validate email
	isValid := utils.IsValidEmail(RegForm.Email)
	if !isValid {
		res := utils.NewError(http.StatusBadRequest, "invalid email address")
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
		return
	}

	//validate phoneNumber
	isValid, formattedNum, err := utils.IsValidPhoneNumber(RegForm.Phone, "NG")
	if !isValid || err != nil {
		res := utils.NewError(http.StatusBadRequest, "invalid phone number")
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
		return
	}

	RegForm.Phone = formattedNum

	err = services.RegisterService(RegForm)
	if err != nil {
		res := utils.NewError(http.StatusInternalServerError, err.Error())
		c.AbortWithStatusJSON(res.StatusCode, gin.H{"error": res.Error})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "successfully created"})
}
