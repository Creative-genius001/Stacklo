package handler

import (
	"io"
	"net/http"

	"github.com/Creative-genius001/Stacklo/api-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

var walletService = client.NewServiceClient("http://localhost:8001/api/v1/wallet")

func ProxyGetAllWallets(c *gin.Context) {
	userID := c.Param("id")
	walletService.DoRequest(c, http.MethodGet, "/"+userID, nil)
}

func ProxyGetFiatWallet(c *gin.Context) {
	userID := c.Param("id")
	walletService.DoRequest(c, http.MethodGet, "/fiat"+userID, nil)
}

func ProxyCreateFiatWallet(c *gin.Context) {
	var body io.Reader
	body = c.Request.Body
	if body == nil {
		body = nil
	}
	walletService.DoRequest(c, http.MethodPost, "/fiat/create", body)
}

func ProxyCreateCryptoWallet(c *gin.Context) {
	var body io.Reader
	body = c.Request.Body
	if body == nil {
		body = nil
	}
	walletService.DoRequest(c, http.MethodPost, "/crypto/create", body)
}
