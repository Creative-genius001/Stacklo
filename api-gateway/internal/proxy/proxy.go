package proxy

import (
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var serviceMap = map[string]string{
	"user":        "http://localhost:8000",
	"wallet":      "http://localhost:8001",
	"transaction": "http://localhost:8002",
	"payment":     "http://localhost:8003",
}

func SetupRoutes(r *gin.RouterGroup) {
	for name, target := range serviceMap {
		u, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(u)

		route := "/" + name + "/*path"

		r.Any(route, func(proxy *httputil.ReverseProxy) gin.HandlerFunc {
			return func(c *gin.Context) {
				c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/"+name)
				proxy.ServeHTTP(c.Writer, c.Request)
			}
		}(proxy))
	}

}
