package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/Kihouww/mini-ai-compute-platform/internal/api"
	"github.com/Kihouww/mini-ai-compute-platform/internal/config"
	"github.com/Kihouww/mini-ai-compute-platform/internal/middleware"
	"github.com/Kihouww/mini-ai-compute-platform/internal/repository"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load("config.yaml")
	if err != nil {
		logger.Error("load_config_failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db, err := repository.NewMySQLDB(cfg.MySQL)
	if err != nil {
		logger.Error("connect_mysql_failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if err := repository.InitRequestLogsTable(db); err != nil {
		logger.Error("init_request_logs_table_failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	redisClient, err := repository.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Error("connect_redis_failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer redisClient.Close()

	requestLogRepo := repository.NewRequestLogRepository(db)
	authMiddleware := middleware.AuthMiddleware(cfg.Auth.APIKeys)
	rateLimitMiddleware := middleware.RateLimitMiddleware(redisClient, cfg.RateLimit.PerMinute)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger(logger))

	api.RegisterRoutes(r, cfg, requestLogRepo, authMiddleware, rateLimitMiddleware)

	logger.Info("mysql_connected")
	logger.Info("redis_connected")
	logger.Info("rate_limit_enabled", slog.Int("per_minute", cfg.RateLimit.PerMinute))
	logger.Info("server_started", slog.Int("port", cfg.Server.Port))

	if err := r.Run(cfg.Addr()); err != nil {
		logger.Error("server_run_failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
