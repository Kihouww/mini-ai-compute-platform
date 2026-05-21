package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kihouww/mini-ai-compute-platform/internal/model"
)

type RequestLogRepository struct {
	db *sql.DB
}

func NewRequestLogRepository(db *sql.DB) *RequestLogRepository {
	return &RequestLogRepository{
		db: db,
	}
}

func (r *RequestLogRepository) Create(ctx context.Context, item *model.RequestLog) error {
	query := `
INSERT INTO request_logs (
  user_id,
  api_key,
  model,
  prompt,
  response,
  input_tokens,
  output_tokens,
  latency_ms,
  status,
  error_message
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	result, err := r.db.ExecContext(
		ctx,
		query,
		item.UserID,
		item.APIKey,
		item.Model,
		item.Prompt,
		item.Response,
		item.InputTokens,
		item.OutputTokens,
		item.LatencyMs,
		item.Status,
		item.ErrorMessage,
	)
	if err != nil {
		return fmt.Errorf("insert request log failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		item.ID = id
	}

	return nil
}

func (r *RequestLogRepository) ListRecent(ctx context.Context, limit int) ([]model.RequestLog, error) {
	query := `
SELECT
  id,
  user_id,
  api_key,
  model,
  prompt,
  response,
  input_tokens,
  output_tokens,
  latency_ms,
  status,
  COALESCE(error_message, ''),
  created_at
FROM request_logs
ORDER BY id DESC
LIMIT ?
`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("query request logs failed: %w", err)
	}
	defer rows.Close()

	items := make([]model.RequestLog, 0)

	for rows.Next() {
		var item model.RequestLog

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.APIKey,
			&item.Model,
			&item.Prompt,
			&item.Response,
			&item.InputTokens,
			&item.OutputTokens,
			&item.LatencyMs,
			&item.Status,
			&item.ErrorMessage,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan request log failed: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate request logs failed: %w", err)
	}

	return items, nil
}
