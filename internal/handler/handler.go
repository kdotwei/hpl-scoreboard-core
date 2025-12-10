package handler

import (
	"github.com/kdotwei/hpl-scoreboard/internal/service"
)

type Handler struct {
	service service.Service
}

// NewHandler 建構子，注入 Service 依賴
func NewHandler(s service.Service) *Handler {
	return &Handler{service: s}
}
