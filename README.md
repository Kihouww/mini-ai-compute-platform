# Mini AI Compute Platform

一个面向LLM推理服务的轻量级AI Compute Platform

## 当前目标

- 用户通过HTTP API提交推理请求
- 后端完成鉴权、限流、模型路由、调用记录和监控
- 支持Docker Compose一键部署

## 技术栈

Go / Redis / MySQL / Docker / Prometheus / Grafana

## 第一阶段功能

- REST API
- LLM 请求转发
- 调用日志
- API KEY 鉴权
- Redis 限流
- SSE 流式响应
- Prometheus 指标

## 当前进度

- Day 1: 项目初始化、Go 服务框架

### 本地运行

```zsh
go run ./cmd/server
```
访问：

```zsh
curl http://localhost:8080/
```
---
- Day 2: 接口

### 健康检查

```zsh
curl http://localhost:8080/health
```

### Chat接口

```zsh
curl -X POST
http://localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"model":"mock-llm","prompt":"hello"}'
```

### API 文档

详见:

```text
docs/api.md
```
---

- Day 3: 配置文件、日志 + 项目结构规范化

### 配置文件

支持从 config.yaml 读取配置

首次运行:

```zsh
cp config.example.yaml config.yaml
go run ./cmd/server
```
### 请求日志

服务启动后会打印启动日志，每次访问接口也会打印请求日志

### 当前代码分层

```text
cmd/server		程序入口
internal/api		路由、请求解析、响应返回
internal/service	业务逻辑
internal/model		请求和响应数据结构
internal/config		配置文件读取
internal/middleware	HTTP 中间件
```

### 验证接口

```zsh
curl http://localhost:8080/health
```

```zsh
curl -X POST
http://localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"prompt":"hello"}'
```

---

- Day 4: MySQL

### 启动 MySQL

```zsh
docker compose up -d mysql
```

### Chat 接口写入调用日志

```zsh
curl -X POST http://localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"model":"mock-llm","prompt":"hello mysql"}'
```

### 查询最近请求记录

```zsh
curl http://localhost:8080/v1/requests
```

## MySQL 调用日志

调用记录写入 `request_logs` 表。

可通过以下命令查看：

```zsh
docker exec -it mini-ai-mysql mysql -uroot -ppassword ai_compute -e "SELECT id, model, prompt, status, created_at FROM request_logs ORDER BY id DESC LIMIT 5;"
```

---

- Day 5: Redis

### 启动 Redis

```bash
docker compose up -d redis
```

### Redis 连接测试

服务启动时会 ping Redis。成功后日志中会输出：

```text
redis_connected
```

### Chat 接口鉴权

`/v1/chat` 需要携带 API Key：

```zsh
curl -X POST http://localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-api-key" \
  -d '{"model":"mock-llm","prompt":"hello auth"}'
```

### 查询请求记录

`/v1/requests` 同样需要携带 API Key：

```zsh
curl http://localhost:8080/v1/requests \
  -H "Authorization: Bearer test-api-key"
```
