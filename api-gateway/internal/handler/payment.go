package handler

import (
	"io"
	"net/http"

	"github.com/Creative-genius001/Stacklo/api-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

var paymentService = client.NewServiceClient("http://localhost:8003/api/v1/payment")

func ProxyGetBankList(c *gin.Context) {
	paymentService.DoRequest(c, http.MethodGet, "/bank-list", nil)
}

func ProxyResolveAccountNumber(c *gin.Context) {
	paymentService.DoRequest(c, http.MethodGet, "/account-details", nil)
}

func ProxyGetOTP(c *gin.Context) {
	var body io.Reader
	body = c.Request.Body
	if body == nil {
		body = nil
	}
	paymentService.DoRequest(c, http.MethodPost, "/transfer/otp", body)
}

func ProxyTransfer(c *gin.Context) {
	var body io.Reader
	body = c.Request.Body
	if body == nil {
		body = nil
	}
	paymentService.DoRequest(c, http.MethodPost, "/transfer", body)
}

func ProxyConvert(c *gin.Context) {
	paymentService.DoRequest(c, http.MethodGet, "/ping", nil)
}

func ProxyPing(c *gin.Context) {
	paymentService.DoRequest(c, http.MethodGet, "/convert", nil)
}
