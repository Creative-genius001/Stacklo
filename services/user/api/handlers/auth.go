package handlers

import (
	"net/http"

	"github.com/Creative-genius001/user/api/services"
	"github.com/Creative-genius001/user/types"
	"github.com/Creative-genius001/user/utils"
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
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	//validate phoneNumber
	isValid, formattedNum, err := utils.IsValidPhoneNumber(RegForm.Phone, "NG")
	if !isValid || err != nil {
		res := utils.NewError(http.StatusBadRequest, "invalid phone number")
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	RegForm.Phone = formattedNum

	//hash password
	passwordHash, err := utils.HashPassword(RegForm.Password)
	if err != nil {
		res := utils.NewError(http.StatusInternalServerError, "sign up failed")
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	RegForm.Password = passwordHash

	err = services.RegisterService(RegForm)
	if err != nil {
		res := utils.NewError(http.StatusInternalServerError, err.Error())
		c.AbortWithStatusJSON(res.StatusCode, res)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "successfully created"})
}
