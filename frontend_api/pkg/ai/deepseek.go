package ai

import "context"

// DeepSeek 实现
// TODO: 实现 DeepSeek API 调用
type DeepSeek struct{}

func NewDeepSeek() *DeepSeek {
	return &DeepSeek{}
}

func (d *DeepSeek) Chat(ctx context.Context, prompt string) (string, error) {
	// TODO: 调用 DeepSeek API
	return "", nil
}
