package client

import (
	"io"
	"net/http"
	"time"

	"github.com/Creative-genius001/Stacklo/api-gateway/internal/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ServiceClient struct {
	BaseURL string
	Client  *http.Client
}

func NewServiceClient(baseURL string) *ServiceClient {
	return &ServiceClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *ServiceClient) DoRequest(c *gin.Context, method string, path string, body io.Reader) {
	url := s.BaseURL + path

	logger.Logger.Debug("URL-PATH", zap.String("url", url))

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	req.Header = c.Request.Header.Clone()

	resp, err := s.Client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			c.Writer.Header().Add(k, vv)
		}
	}
	c.Writer.WriteHeader(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
