package ai

import "context"

// Client AI 客户端接口
// TODO: 后续实现具体的 AI 调用逻辑
type Client interface {
	Chat(ctx context.Context, prompt string) (string, error)
}
