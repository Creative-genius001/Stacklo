package handler

import (
	er "errors"
	"net/http"

	"github.com/Creative-genius001/Stacklo/services/user/api/service/auth"
	"github.com/Creative-genius001/Stacklo/services/user/model"
	"github.com/Creative-genius001/Stacklo/services/user/types"
	"github.com/Creative-genius001/Stacklo/services/user/utils"
	errors "github.com/Creative-genius001/Stacklo/services/user/utils/error"
	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	auth auth.Auth
}

func NewAuthHandler(a auth.Auth) *AuthHandler {
	return &AuthHandler{a}
}

func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login route"})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var form types.RegisterType

	if err := c.ShouldBind(&form); err != nil {
		logger.Logger.Info("invalid input data")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"error": errors.TypeInvalidInput})
		return
	}

	//validate email
	isValid := utils.IsValidEmail(form.Email)
	if !isValid {
		logger.Logger.Info("invalid email address")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"error": errors.TypeInvalidInput})
		return
	}

	//validate phoneNumber
	isValid, formatPhone, err := utils.IsValidPhoneNumber(form.Phone, "NG")
	if !isValid || err != nil {
		logger.Logger.Info("invalid phone number")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"error": errors.TypeInvalidInput})
		return
	}

	form.Phone = formatPhone

	var PasswordHash string
	form.Password = PasswordHash

	payload := model.User{
		Email:        form.Email,
		PasswordHash: form.Password,
		Country:      form.Country,
		FirstName:    form.FirstName,
		LastName:     form.LastName,
		PhoneNumber:  form.Phone,
	}

	user, err := h.auth.CreateUser(c, payload)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error from PaystackAPIWrapper", zap.Error(err))
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		switch appErr.Type {
		case errors.TypeInternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		case errors.TypeConflict:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeConflict})
			return
		default:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "user successfully registered", "data": user})
}
