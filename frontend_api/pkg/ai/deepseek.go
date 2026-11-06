package ai

import "time"

// DeepSeekModel DeepSeek 默认模型名
const DeepSeekModel = "deepseek-v4-flash"

// NewDeepSeek 创建 DeepSeek 客户端
// 使用前需在 config.yaml 中配置 deepseek 模型的 api_key 和 endpoint
func NewDeepSeek(apiKey, endpoint string) Client {
	return NewClient(endpoint, apiKey, DeepSeekModel)
}

// NewDeepSeekWithOptions 创建带超时与重试次数的 DeepSeek 客户端
func NewDeepSeekWithOptions(apiKey, endpoint string, timeout time.Duration, maxRetries int) Client {
	return NewClientWithOptions(endpoint, apiKey, DeepSeekModel, timeout, maxRetries)
}
