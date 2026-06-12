package ai

// NewQwen 创建通义千问客户端（OpenAI 兼容模式）
// 使用前需在 config.yaml 中配置 qwen 模型的 api_key 和 endpoint
func NewQwen(apiKey, endpoint string) Client {
	return NewClient(endpoint, apiKey, "qwen3-vl-plus")
}
