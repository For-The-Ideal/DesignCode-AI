package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"frontend_api/config"
	"io"
	"net/http"
	"time"
)

// ═══════════════════════════════════════════════
//  AI Client — OpenAI 兼容格式的 HTTP 客户端
//  同时支持 DeepSeek / Qwen(通义千问) / OpenAI 等
// ═══════════════════════════════════════════════

// Client AI 客户端接口
type Client interface {
	// Chat 发送对话请求，返回 AI 响应的文本内容
	Chat(ctx context.Context, messages []Message, opts ...Option) (string, error)
}

// Message 对话消息
// Content 可以是 string（纯文本）或 []ContentPart（多模态）
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // string 或 []ContentPart
}

// ContentPart 多模态消息的内容片段
type ContentPart struct {
	Type     string    `json:"type"` // "text" | "image_url"
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL 图片引用
type ImageURL struct {
	URL string `json:"url"`
}

// TextPart 创建文本内容片段
func TextPart(text string) ContentPart {
	return ContentPart{Type: "text", Text: text}
}

// ImagePart 创建图片内容片段
func ImagePart(url string) ContentPart {
	return ContentPart{Type: "image_url", ImageURL: &ImageURL{URL: url}}
}

// chatRequest OpenAI 兼容的请求体
type chatRequest struct {
	Model           string          `json:"model"`
	Messages        []Message       `json:"messages"`
	Stream          bool            `json:"stream"`
	MaxTokens       int             `json:"max_tokens,omitempty"`
	Temperature     float64         `json:"temperature,omitempty"`
	Thinking        *thinkingConfig `json:"thinking,omitempty"`
	ReasoningEffort string          `json:"reasoning_effort,omitempty"`
}

// thinkingConfig DeepSeek 思维链配置
type thinkingConfig struct {
	Type string `json:"type"` // enabled | disabled
}

// chatResponse OpenAI 兼容的响应体
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// chatClient Client 的 HTTP 实现
type chatClient struct {
	httpClient *http.Client
	endpoint   string
	apiKey     string
	model      string
}

// Option 请求选项
type Option func(*chatRequest)

// WithMaxTokens 设置 max_tokens
func WithMaxTokens(tokens int) Option {
	return func(r *chatRequest) {
		r.MaxTokens = tokens
	}
}

// WithTemperature 设置 temperature
func WithTemperature(temp float64) Option {
	return func(r *chatRequest) {
		r.Temperature = temp
	}
}

// WithDeepSeekThinking 启用 DeepSeek 思维链
func WithDeepSeekThinking(effort string) Option {
	return func(r *chatRequest) {
		r.Thinking = &thinkingConfig{Type: "enabled"}
		r.ReasoningEffort = effort
	}
}

// NewClient 创建 AI 客户端
//
//	endpoint: API 地址，如 "https://api.deepseek.com/v1/chat/completions"
//	apiKey:   API Key
//	model:    模型名称，如 "deepseek-chat" / "qwen-vl-plus"
func NewClient(endpoint, apiKey, model string) Client {
	return &chatClient{
		httpClient: &http.Client{
			Timeout: 600 * time.Second, // 10 分钟，长代码生成任务可能较慢
		},
		endpoint: endpoint,
		apiKey:   apiKey,
		model:    model,
	}
}

// NewClientWithTimeout 创建 AI 客户端并指定超时
func NewClientWithTimeout(endpoint, apiKey, model string, timeout time.Duration) Client {
	return &chatClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		endpoint: endpoint,
		apiKey:   apiKey,
		model:    model,
	}
}

// deepSeekError 映射 HTTP 状态码到可读的错误描述
var deepSeekError = map[int]string{
	400: "请求体格式错误",
	401: "API key 认证失败",
	402: "账号余额不足",
	422: "请求体参数错误",
	429: "请求速率达到上限 (TPM/RPM)",
	500: "服务器内部故障",
	503: "服务器负载过高，请稍后重试",
}

// Chat 发送对话请求
func (c *chatClient) Chat(ctx context.Context, messages []Message, opts ...Option) (string, error) {
	reqBody := chatRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
	}
	for _, opt := range opts {
		opt(&reqBody)
	}

	body, _ := json.Marshal(reqBody)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ai: 创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("ai: 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ai: 读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errDesc := deepSeekError[resp.StatusCode]
		if errDesc == "" {
			errDesc = fmt.Sprintf("未知错误 (HTTP %d)", resp.StatusCode)
		}
		return "", fmt.Errorf("ai: %s, body=%s", errDesc, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("ai: 解析响应失败: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("ai: API 错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("ai: 响应中无 choices 字段")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// SelectClient 按能力从配置中选取第一个启用的模型并创建客户端
//
//	capability: "read" | "write" | "both"
func SelectClient(capability string) (Client, error) {
	for _, m := range config.AppConfig.AI.Models {
		if !m.Enabled {
			continue
		}
		if m.Capability == capability || m.Capability == "both" {
			switch m.Provider {
			case "deepseek":
				return NewDeepSeek(m.APIKey, m.Endpoint), nil
			case "qwen":
				return NewQwen(m.APIKey, m.Endpoint), nil
			default:
				return NewClient(m.Endpoint, m.APIKey, m.Provider), nil
			}
		}
	}
	return nil, fmt.Errorf("ai: 未找到 capability=%s 的已启用模型", capability)
}
