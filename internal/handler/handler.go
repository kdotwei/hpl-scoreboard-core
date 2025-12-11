package handler

import (
	"github.com/kdotwei/hpl-scoreboard/internal/service"
	"github.com/kdotwei/hpl-scoreboard/internal/token"
)

type Handler struct {
	service    service.Service
	tokenMaker token.Maker // 新增依賴
}

// NewHandler 更新建構子，注入 TokenMaker
func NewHandler(s service.Service, tm token.Maker) *Handler {
	return &Handler{
		service:    s,
		tokenMaker: tm,
	}
}
