package ai

import "time"

// QwenModel 通义千问默认模型名
const QwenModel = "qwen3-vl-plus"

// NewQwen 创建通义千问客户端（OpenAI 兼容模式）
// 使用前需在 config.yaml 中配置 qwen 模型的 api_key 和 endpoint
func NewQwen(apiKey, endpoint string) Client {
	return NewClient(endpoint, apiKey, QwenModel)
}

// NewQwenWithOptions 创建带超时与重试次数的通义千问客户端
func NewQwenWithOptions(apiKey, endpoint string, timeout time.Duration, maxRetries int) Client {
	return NewClientWithOptions(endpoint, apiKey, QwenModel, timeout, maxRetries)
}
