package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Kihouww/mini-ai-compute-platform/internal/config"
	"github.com/Kihouww/mini-ai-compute-platform/internal/model"
	"github.com/Kihouww/mini-ai-compute-platform/internal/service"
)

type Handler struct {
	Config *config.Config
}

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	h := &Handler{
		Config: cfg,
	}

	r.GET("/health", h.HealthHandler)
	r.POST("/v1/chat", h.ChatHandler)
}

func (h *Handler) HealthHandler(c *gin.Context) {
	Success(c, gin.H{
		"status": "ok",
	})
}

func (h *Handler) ChatHandler(c *gin.Context) {
	var req model.ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	resp := service.Chat(req, h.Config.LLM.DefaultModel)
	Success(c, resp)
}
