package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Kihouww/mini-ai-compute-platform/internal/config"
	"github.com/Kihouww/mini-ai-compute-platform/internal/model"
	"github.com/Kihouww/mini-ai-compute-platform/internal/repository"
	"github.com/Kihouww/mini-ai-compute-platform/internal/service"
)

type Handler struct {
	Config         *config.Config
	RequestLogRepo *repository.RequestLogRepository
}

func RegisterRoutes(
	r *gin.Engine,
	cfg *config.Config,
	requestLogRepo *repository.RequestLogRepository,
	authMiddleware gin.HandlerFunc,
	rateLimitMiddleware gin.HandlerFunc,
) {
	h := &Handler{
		Config:         cfg,
		RequestLogRepo: requestLogRepo,
	}

	r.GET("/health", h.HealthHandler)

	v1 := r.Group("/v1")
	v1.Use(authMiddleware)
	v1.Use(rateLimitMiddleware)

	v1.POST("/chat", h.ChatHandler)
	v1.GET("/requests", h.ListRequestsHandler)
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

	apiKey := ""
	if value, exists := c.Get("api_key"); exists {
		if key, ok := value.(string); ok {
			apiKey = key
		}
	}

	requestLog := &model.RequestLog{
		UserID:       "anonymous",
		APIKey:       apiKey,
		Model:        resp.Model,
		Prompt:       req.Prompt,
		Response:     resp.Answer,
		InputTokens:  0,
		OutputTokens: 0,
		LatencyMs:    int(resp.LatencyMs),
		Status:       "success",
		ErrorMessage: "",
	}

	if err := h.RequestLogRepo.Create(c.Request.Context(), requestLog); err != nil {
		Error(c, http.StatusInternalServerError, 500, "save request log failed")
		return
	}

	Success(c, resp)
}

func (h *Handler) ListRequestsHandler(c *gin.Context) {
	items, err := h.RequestLogRepo.ListRecent(c.Request.Context(), 20)
	if err != nil {
		Error(c, http.StatusInternalServerError, 500, "query request logs failed")
		return
	}

	Success(c, items)
}
