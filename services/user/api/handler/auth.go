package handler

import (
	er "errors"
	"net/http"

	"github.com/Creative-genius001/Stacklo/services/user/api/service"
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
	otp  service.OTPServ
}

func NewAuthHandler(a auth.Auth, o service.OTPServ) *AuthHandler {
	return &AuthHandler{a, o}
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
		case errors.TypeConflict:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Message})
			return
		case errors.TypeNotFound:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Message})
			return
		default:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "login successful", "data": user})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var form types.RegisterType

	if err := c.ShouldBind(&form); err != nil {
		logger.Logger.Info("invalid input data")
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	payload := model.User{
		Email:        form.Email,
		PasswordHash: form.Password,
		Country:      form.Country,
		FirstName:    form.FirstName,
		LastName:     form.LastName,
		PhoneNumber:  form.Phone,
		IsVerified:   false,
		KycStatus:    "not_started",
	}

	err := h.auth.Register(c, payload)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error", zap.Error(err))
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		switch appErr.Type {
		case errors.TypeConflict:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Message})
			return
		case errors.TypeInvalidInput:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Message})
			return
		default:
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "OTP verification code sent"})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	email := c.Query("email")
	otp := c.Query("otp")

	if otp == "" || email == "" {
		c.JSON(errors.GetHTTPStatus(errors.TypeInvalidInput), gin.H{"status": "error", "message": errors.TypeInvalidInput})
		return
	}

	err := h.auth.SignupOTPVerification(c.Request.Context(), email, otp)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error", zap.Error(err))
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		if appErr.Type == errors.TypeInternal {
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		} else {
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Message})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "verification successful"})
}

func (h *AuthHandler) ResendOTP(c *gin.Context) {
	email := c.Query("email")

	err := h.otp.SendOTP(email)
	if err != nil {
		var appErr *errors.CustomError
		if !er.As(err, &appErr) {
			logger.Logger.Error("Unexpected error", zap.Error(err))
			c.JSON(errors.GetHTTPStatus(errors.TypeInternal), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		}
		if appErr.Type == errors.TypeInternal {
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": errors.TypeInternal})
			return
		} else {
			c.JSON(errors.GetHTTPStatus(appErr.Type), gin.H{"status": "error", "message": appErr.Message})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "verification code sent"})
}
