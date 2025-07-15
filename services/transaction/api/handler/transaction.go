package handler

import (
	services "github.com/Creative-genius001/Stacklo/services/transaction/api/service"
)

type Handler struct {
	service services.Service
}

func NewHandler(s services.Service) *Handler {
	return &Handler{service: s}
}
