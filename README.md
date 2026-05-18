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

- Day 1: 项目初始化、Go 服务框架、README 初版

## 本地运行

```zsh
go run ./cmd/server
```
访问：

```zsh
curl http://localhost:8080/
```

