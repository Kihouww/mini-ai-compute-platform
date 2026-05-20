package model

type ChatRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt" binding:"required"`
}

type ChatResponse struct {
	Model     string `json:"model"`
	Answer    string `json:"answer"`
	LatencyMs int64  `json:"latency_ms"`
}
