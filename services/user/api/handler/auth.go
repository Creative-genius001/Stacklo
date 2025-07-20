package handler

import (
	er "errors"
	"net/http"

	"github.com/Creative-genius001/Stacklo/services/user/api/service/auth"
	"github.com/Creative-genius001/Stacklo/services/user/model"
	"github.com/Creative-genius001/Stacklo/services/user/types"
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

	var form types.LoginType

	if err := c.ShouldBind(&form); err != nil {
		logger.Logger.Info("invalid input data")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "error": errors.TypeInvalidInput})
		return
	}

	user, err := h.auth.Login(c, form.Email, form.Password)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error", zap.Error(err))
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		switch appErr.Type {
		case errors.TypeInternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		case errors.TypeConflict:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Err})
			return
		default:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
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

func (h *AuthHandler) Register(c *gin.Context) {
	var form types.RegisterType

	if err := c.ShouldBind(&form); err != nil {
		logger.Logger.Info("invalid input data")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "error": errors.TypeInvalidInput})
		return
	}

	payload := model.User{
		Email:        form.Email,
		PasswordHash: form.Password,
		Country:      form.Country,
		FirstName:    form.FirstName,
		LastName:     form.LastName,
		PhoneNumber:  form.Phone,
		KycStatus:    "not_started",
	}

	user, err := h.auth.Register(c, payload)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error", zap.Error(err))
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		switch appErr.Type {
		case errors.TypeInternal:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		case errors.TypeConflict:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Err})
			return
		default:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
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

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "user successfully registered", "data": data})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {

}
