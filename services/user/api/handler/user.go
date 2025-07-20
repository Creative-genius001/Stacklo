package handler

import (
	"net/http"

	"github.com/Creative-genius001/Stacklo/services/user/api/service"
	"github.com/Creative-genius001/Stacklo/services/user/model"
	"github.com/Creative-genius001/Stacklo/services/user/types"
	errors "github.com/Creative-genius001/Stacklo/services/user/utils/error"
	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.Service
}

func NewUserHandler(s service.Service) *UserHandler {
	return &UserHandler{s}
}

func (u *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := u.service.GetUser(c, userID)
	if err != nil {
		appErr, _ := err.(*errors.CustomError)
		if appErr.Type == errors.TypeInvalidInput {
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "error": appErr.Message})
			return
		} else {
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "error": err})
			return
		}
	}

	data := model.User{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Country:     user.Country,
		KycStatus:   user.KycStatus,
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "login successful", "data": data})

}

func (u *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.GetString("id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}
	var data types.UpdateUser
	if err := c.ShouldBind(&data); err != nil {
		logger.Logger.Info("invalid input data")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "error": errors.TypeInvalidInput})
		return
	}
	// userID := c.Param("id")
	// isValid := utils.IsValidUUID(userID)
	// if !isValid {
	// 	logger.Logger.Error("user id is invalid and could not be parsed")
	// 	c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status0": "error", "message": "user id could not be parsed"})
	// 	return
	// }
	// data.ID = userID

	err := u.service.UpdateUser(c, userID, data)
	if err != nil {
		appErr, _ := err.(*errors.CustomError)
		if appErr.Type == errors.TypeInvalidInput {
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Message})
			return
		} else {
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": "failed to update user"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "update successful",
	})
}
