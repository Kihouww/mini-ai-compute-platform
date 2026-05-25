# API 文档

## 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 鉴权说明

`/v1` 分组下的接口需要 API Key。

请求头格式：

```http
Authorization: Bearer test-api-key
```

错误响应：

| 场景 | HTTP 状态码 | code | message |
| --- | --- | --- | --- |
| 未提供 Authorization | 401 | 401 | missing authorization header |
| Authorization 格式错误 | 401 | 401 | invalid authorization format |
| API Key 为空 | 401 | 401 | empty api key |
| API Key 错误 | 403 | 403 | invalid api key |
| 超过限流 | 429 | 429 | rate limit exceeded |

---

## 健康检查

### 请求

```http
GET /health
```

### 是否需要鉴权

否

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok"
  }
}
```

---

## Chat 接口

### 请求

```http
POST /v1/chat
```

### 是否需要鉴权

是

### 请求头

```http
Authorization: Bearer test-api-key
Content-Type: application/json
```

### 请求参数

```json
{
  "model": "mock-llm",
  "prompt": "hello"
}
```

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| model | string | 否 | 模型名称，不传则使用默认模型 |
| prompt | string | 是 | 用户输入内容 |

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "model": "mock-llm",
    "answer": "mock response",
    "latency_ms": 12
  }
}
```

### curl 示例

```bash
curl -X POST http://localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-api-key" \
  -d '{"model":"mock-llm","prompt":"hello"}'
```

---

## 查询请求记录

### 请求

```http
GET /v1/requests
```

### 是否需要鉴权

是

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user_id": "anonymous",
      "api_key": "test-api-key",
      "model": "mock-llm",
      "prompt": "hello",
      "response": "mock response",
      "input_tokens": 0,
      "output_tokens": 0,
      "latency_ms": 12,
      "status": "success",
      "error_message": "",
      "created_at": "2026-05-25T10:00:00+08:00"
    }
  ]
}
```

### curl 示例

```bash
curl http://localhost:8080/v1/requests \
  -H "Authorization: Bearer test-api-key"
```

---

## Redis 限流

### 限流规则

每个 API Key 每分钟最多请求 20 次。

Redis key：

```text
rate_limit:{api_key}:{minute}
```

### 超限响应

```json
{
  "code": 429,
  "message": "rate limit exceeded",
  "data": null
}
```

---

## 用户注册

### 请求

```http
POST /v1/users/register
```

### 是否需要鉴权

否

### 请求参数

```json
{
  "username": "testuser",
  "password": "123456",
  "email": "test@example.com"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "username": "testuser"
  }
}
```

---

## 用户登录

### 请求

```http
POST /v1/users/login
```

### 是否需要鉴权

否

### 请求参数

```json
{
  "username": "testuser",
  "password": "123456"
}
```

### 响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "mock-token-user-1-xxxxxxxxxx",
    "user_id": 1
  }
}
```

### 说明

当前 token 是 mock token，用于 Day 8 验证登录流程。后续升级为 JWT。
