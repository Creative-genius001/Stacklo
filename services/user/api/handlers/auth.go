package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login route"})
}

func Register(c *gin.Context) {

}
