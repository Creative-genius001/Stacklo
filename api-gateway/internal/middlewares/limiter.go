package middlewares

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
)

type ClientLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	burst    int
}

func NewClientLimiter(r rate.Limit, burst int) *ClientLimiter {
	return &ClientLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		burst:    burst,
	}
}

func (cl *ClientLimiter) getLimiter(ip string) *rate.Limiter {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	limiter, exists := cl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(cl.r, cl.burst)
		cl.limiters[ip] = limiter
	}

	return limiter
}

func (cl *ClientLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := cl.getLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
