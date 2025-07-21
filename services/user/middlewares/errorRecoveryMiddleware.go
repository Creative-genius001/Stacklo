package middlewares

import (
	"net/http"

	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Log the panic with stack trace
				logger.Logger.Error("Panic recovered in HTTP handler",
					zap.Any("panic_value", r),
					zap.Stack("stacktrace"),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "An unexpected server error occurred.",
				})
			}
		}()
		c.Next()
	}
}
