# 项目架构

## 当前架构

```text
Client
  ↓
API Server
  ↓
Auth Middleware
  ↓
Rate Limit Middleware
  ↓
Chat Handler
  ↓
Chat Service
  ↓
Mock LLM

Chat Handler
  ↓
RequestLogRepository
  ↓
MySQL request_logs

Rate Limit Middleware
  ↓
Redis
```

## 模块说明

### cmd/server

程序入口，负责初始化配置、日志、MySQL、Redis、Gin 路由和中间件。

### internal/api

负责 HTTP 路由注册、请求参数解析、响应返回。

### internal/service

负责业务逻辑，比如模型选择、模拟 LLM 响应、耗时统计。

### internal/model

负责请求、响应、数据库记录等数据结构定义。

### internal/config

负责读取 config.yaml，并提供默认配置。

### internal/middleware

负责通用 HTTP 中间件，目前包括请求日志、API Key 鉴权和 Redis 限流。

### internal/repository

负责数据库和缓存访问，目前包括 MySQL 连接、请求日志读写、Redis 连接。

## 当前请求链路

### /v1/chat

```text
1. Client 发送 POST /v1/chat
2. AuthMiddleware 校验 Authorization Header
3. RateLimitMiddleware 基于 Redis 检查 API Key 请求次数
4. ChatHandler 解析 JSON 请求体
5. ChatService 生成 mock response
6. RequestLogRepository 写入 MySQL request_logs
7. API 返回统一 JSON 响应
```

## 限流设计

当前使用固定窗口限流。

Redis key：

```text
rate_limit:{api_key}:{minute}
```

示例：

```text
rate_limit:test-api-key:202605251430
```

每次请求执行 INCR。

第一次创建 key 时设置 60 秒过期。

超过配置限制后返回 429。
