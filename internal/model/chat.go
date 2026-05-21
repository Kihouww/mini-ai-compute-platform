package model

import "time"

type ChatRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt" binding:"required"`
}

type ChatResponse struct {
	Model     string `json:"model"`
	Answer    string `json:"answer"`
	LatencyMs int64  `json:"latency_ms"`
}

type RequestLog struct {
	ID           int64     `json:"id"`
	UserID       string    `json:"user_id"`
	APIKey       string    `json:"api_key"`
	Model        string    `json:"model"`
	Prompt       string    `json:"prompt"`
	Response     string    `json:"response"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	LatencyMs    int       `json:"latency_ms"`
	Status       string    `json:"status"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
}
