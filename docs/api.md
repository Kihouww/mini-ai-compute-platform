# API 文档

## 统一响应格式

```json
{
	"code": 0,
	"message": "success",
	"data": {}
}
```

---

## 健康检查接口

### 接口名称

健康检查

### 请求路径

```http
GET /health
```

### 请求方法

```http
GET
```

### 请求参数

无

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

### 错误码

暂无

---

## Chat 接口

### 接口名称

提交 LLM 推理请求

### 提交路径

```http
POST /v1/chat
```

### 请求方法

```http
POST
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
| model | string | 是 | 模型名称，当前使用mock-llm |
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

### 错误码

| code | message | 说明 
| --- | --- | --- |
| 400 | invalid request body | 请求体错误或缺少必填字段 |

---

## 请求记录接口

### 接口名称

查询最近请求记录

### 请求路径

```http
GET /v1/requests
```

### 请求参数

无

### 响应示例

```json
{
	"code": 0,
	"message": "success",
	"data": [
		{
			"id": 1,
			"user_id": "anonymous",
			"api_key": "",
			"model": "mock-llm",
			"prompt": "hello mysql",
			"response": "mock response",
			"input_tokens": 0,
			"output_tokens": 0,
			"latency_ms": 12,
			"status": "success",
			"error_message": "",
			"created_at": "2026-05-21T10:00:00+08:00"
		}
	]
}
```

### 错误码

| code | message | 说明 |
| --- | --- | --- |
| 500 | query request logs failed | 查询请求日志失败 |

---

## API Key 鉴权

### 请求头格式

```http
Authorization: Bearer test-api-key
```

### 鉴权规则

| 场景 | HTTP 状态码 | code | message |
| --- | --- | --- | --- |
| 未提供 Authorization | 401 | 401 | missing authorization header |
| Authorization 格式错误 | 401 | 401 | invalid authorization format |
| API Key 为空 | 401 | 401 | empty api key |
| API Key 错误 | 403 | 403 | invalid api key |

### 需要鉴权的接口

当前 `/v1` 分组下的接口都需要鉴权：

- POST /v1/chat
- GET /v1/requests

### Chat 接口示例

```zsh
curl -X POST http://localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-api-key" \
  -d '{"model":"mock-llm","prompt":"hello auth"}'
```

### 查询请求记录示例

```zsh
curl http://localhost:8080/v1/requests \
  -H "Authorization: Bearer test-api-key"
```

---

## Redis 限流

### 限流规则

当前基于 API Key 做简单固定窗口限流。

限流 key：

```text
rate_limit:{api_key}:{minute}
```

示例：

```text
rate_limit:test-api-key:202605251430
```

### 当前限制

每个 API Key 每分钟最多请求 20 次。

配置项：

```yaml
rate_limit:
  per_minute: 20
```

### 超限响应

HTTP 状态码：

```text
429 Too Many Requests
```

响应示例：

```json
{
  "code": 429,
  "message": "rate limit exceeded",
  "data": null
}
```

### 测试命令

```bash
for i in $(seq 1 25); do
  curl -s -o /dev/null -w "%{http_code}\n" -X POST http://localhost:8080/v1/chat \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer test-api-key" \
    -d "{\"model\":\"mock-llm\",\"prompt\":\"hello $i\"}"
done
```
