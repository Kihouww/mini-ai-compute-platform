package service

import (
	"time"

	"github.com/Kihouww/mini-ai-compute-platform/internal/model"
)

func Chat(req model.ChatRequest) model.ChatResponse {
	start := time.Now()

	time.Sleep(12 * time.Millisecond)

	return model.ChatResponse{
		Model:     req.Model,
		Answer:    "mock response",
		LatencyMs: time.Since(start).Milliseconds(),
	}
}
