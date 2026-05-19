package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Kihouww/mini-ai-compute-platform/internal/model"
	"github.com/Kihouww/mini-ai-compute-platform/internal/service"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/health", HealthHandler)
	r.POST("/v1/chat", ChatHandler)
}

func HealthHandler(c *gin.Context) {
	Success(c, gin.H{
		"status": "ok",
	})
}

func ChatHandler(c *gin.Context) {
	var req model.ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, 400, "invalid request body")
		return
	}

	resp := service.Chat(req)
	Success(c, resp)
}
