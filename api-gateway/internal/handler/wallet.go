package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Creative-genius001/Stacklo/api-gateway/internal/client"
	"github.com/Creative-genius001/Stacklo/api-gateway/internal/utils/logger"
	"github.com/gin-gonic/gin"
)

var walletService = client.NewServiceClient("http://localhost:8001/api/v1/wallet")

func ProxyGetAllWallets(c *gin.Context) {
	userID := c.Param("id")
	walletService.DoRequest(c, http.MethodGet, "/"+userID, nil)
}

func ProxyGetFiatWallet(c *gin.Context) {
	userID := c.Param("id")
	walletService.DoRequest(c, http.MethodGet, "/fiat/"+userID, nil)
}

func ProxyCreateFiatWallet(c *gin.Context) {

	isVerified, ok := c.Get("isVerified")
	userID, ok := c.Get("id")
	if !ok {
		logger.Logger.Warn("isVerified status unavailable")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Request.Header.Set("X-Is-Verified", fmt.Sprintf("%v", isVerified))
	c.Request.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))

	var body io.Reader
	body = c.Request.Body
	if body == nil {
		body = nil
	}
	walletService.DoRequest(c, http.MethodPost, "/fiat/create", body)
}

func ProxyCreateCryptoWallet(c *gin.Context) {
	isVerified, ok := c.Get("isVerified")
	userID, ok := c.Get("id")
	if !ok {
		logger.Logger.Warn("isVerified status unavailable")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.Request.Header.Set("X-Is-Verified", fmt.Sprintf("%v", isVerified))
	c.Request.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))

	var body io.Reader
	body = c.Request.Body
	if body == nil {
		body = nil
	}
	walletService.DoRequest(c, http.MethodPost, "/crypto/create", body)
}
