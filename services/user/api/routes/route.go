package routes

import (
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	UserRoutes(router)
	AuthRoutes(router)
}
