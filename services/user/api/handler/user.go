package handler

import (
	"github.com/Creative-genius001/Stacklo/services/user/api/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	s service.Service
}

func NewUserHandler(s service.Service) *UserHandler {
	return &UserHandler{s}
}

func (u *UserHandler) GetUser(c *gin.Context) {

}

func (u *UserHandler) UpdateUser(c *gin.Context) {

}
