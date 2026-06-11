package ai

import "context"

// OpenAI 实现
// TODO: 实现 OpenAI API 调用
type OpenAI struct{}

func NewOpenAI() *OpenAI {
	return &OpenAI{}
}

func (o *OpenAI) Chat(ctx context.Context, prompt string) (string, error) {
	// TODO: 调用 OpenAI API
	return "", nil
}
