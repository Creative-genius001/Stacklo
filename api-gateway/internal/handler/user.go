package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/Creative-genius001/Stacklo/api-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

var userService = client.NewServiceClient("http://localhost:8000/api/v1/user")
var authService = client.NewServiceClient("http://localhost:8000/api/v1/auth")

func ProxyGetUser(c *gin.Context) {
	userID := c.Param("id")

	userService.DoRequest(c, http.MethodGet, "/"+userID, nil)
}

func ProxyRegister(c *gin.Context) {
	var body io.Reader
	body = c.Request.Body
	if body == nil {
		body = nil
	}
	authService.DoRequest(c, http.MethodPost, "/register", body)
}

func ProxyLogin(c *gin.Context) {
	var body io.Reader
	body = c.Request.Body
	if body == nil {
		body = nil
	}
	authService.DoRequest(c, http.MethodPost, "/login", body)
}

func ProxyVerifyOTP(c *gin.Context) {
	email := c.Query("email")
	otp := c.Query("otp")

	params := url.Values{}
	params.Set("email", email)
	params.Set("otp", otp)

	path := "/verify-otp" + "?" + params.Encode()
	authService.DoRequest(c, http.MethodPost, path, nil)
}

func ProxyUpdateUser(c *gin.Context) {
	userID := c.Param("id")

	userService.DoRequest(c, http.MethodPut, "/"+userID, nil)
}

func ProxyResendOTP(c *gin.Context) {
	email := c.Query("email")

	params := url.Values{}
	params.Set("email", email)

	path := "/resend-otp" + "?" + params.Encode()
	authService.DoRequest(c, http.MethodPost, path, nil)
}
