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
