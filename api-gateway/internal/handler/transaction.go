package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Creative-genius001/Stacklo/api-gateway/internal/client"
	"github.com/gin-gonic/gin"
)

var transactionService = client.NewServiceClient("http://localhost:8002/api/v1/transaction")

func ProxyTransactionByID(c *gin.Context) {
	transactionID := c.Param("id")
	transactionService.DoRequest(c, http.MethodGet, "/"+transactionID, nil)
}

func ProxyCreateTransaction(c *gin.Context) {
	var body io.Reader
	body = c.Request.Body
	if body == nil {
		body = nil
	}
	transactionService.DoRequest(c, http.MethodPost, "/create", body)
}

func ProxyGetAllTransactions(c *gin.Context) {
	userID := c.Query("id")
	url := fmt.Sprintf("?id=%s", userID)
	transactionService.DoRequest(c, http.MethodGet, url, nil)
}

func ProxyGetFilteredTransactions(c *gin.Context) {
	userID := c.Param("id")
	transactionType := c.Query("transaction_type")
	entryType := c.Query("entry_type")
	status := c.Query("status")
	limitStr := c.Query("limit")
	cursorStr := c.Query("cursor")

	params := url.Values{}
	params.Set("transaction_type", transactionType)
	params.Set("entry_type", entryType)
	params.Set("status", status)
	params.Set("limit", limitStr)
	params.Set("cursor", cursorStr)

	path := "/" + userID + "?" + params.Encode()
	transactionService.DoRequest(c, http.MethodGet, path, nil)
}
