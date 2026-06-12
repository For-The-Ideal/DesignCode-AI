package ai

// NewDeepSeek 创建 DeepSeek 客户端
// 使用前需在 config.yaml 中配置 deepseek 模型的 api_key 和 endpoint
func NewDeepSeek(apiKey, endpoint string) Client {
	return NewClient(endpoint, apiKey, "deepseek-v4-flash")
}
