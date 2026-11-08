package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"frontend_api/config"
	"frontend_api/utils"
	"io"
	"log"
	"net/http"
	"strconv"
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
	// ChatWithMeta 与 Chat 相同，但额外返回结束原因与用量，
	// 供生成流程判断输出是否被 max_tokens 截断（finish_reason=length）
	ChatWithMeta(ctx context.Context, messages []Message, opts ...Option) (Result, error)
}

// Result 单次对话的结果
//
//	Content          - 响应正文
//	FinishReason     - stop | length | ...（length 表示被 max_tokens 截断）
//	CompletionTokens - 本次输出消耗的 token 数
type Result struct {
	Content          string
	FinishReason     string
	CompletionTokens int
}

// Truncated 输出是否因 max_tokens 上限被截断
func (r Result) Truncated() bool {
	return r.FinishReason == "length"
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
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage *struct {
		CompletionTokens int `json:"completion_tokens"`
	} `json:"usage,omitempty"`
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
	maxRetries int
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
	return NewClientWithOptions(endpoint, apiKey, model, 600*time.Second, 0)
}

// NewClientWithTimeout 创建 AI 客户端并指定超时
func NewClientWithTimeout(endpoint, apiKey, model string, timeout time.Duration) Client {
	return NewClientWithOptions(endpoint, apiKey, model, timeout, 0)
}

// NewClientWithOptions 创建 AI 客户端，同时指定超时与最大重试次数
// maxRetries 只针对限流（429）与临时故障（5xx/网络错误）生效
func NewClientWithOptions(endpoint, apiKey, model string, timeout time.Duration, maxRetries int) Client {
	if timeout <= 0 {
		timeout = 600 * time.Second
	}
	if maxRetries < 0 {
		maxRetries = 0
	}
	return &chatClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		endpoint:   endpoint,
		apiKey:     apiKey,
		model:      model,
		maxRetries: maxRetries,
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

// Chat 发送对话请求，只返回响应正文
func (c *chatClient) Chat(ctx context.Context, messages []Message, opts ...Option) (string, error) {
	res, err := c.ChatWithMeta(ctx, messages, opts...)
	return res.Content, err
}

// ChatWithMeta 发送对话请求并返回结束原因；限流/临时故障时按 maxRetries 做指数退避重试
func (c *chatClient) ChatWithMeta(ctx context.Context, messages []Message, opts ...Option) (Result, error) {
	reqBody := chatRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
	}
	for _, opt := range opts {
		opt(&reqBody)
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return Result{}, fmt.Errorf("ai: 序列化请求失败: %w", err)
	}

	var lastErr error
	var serverWait time.Duration
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			delay := utils.RetryDelay(attempt)
			// 服务端给了 Retry-After 且更大时，优先听从服务端
			if serverWait > delay && serverWait <= 30*time.Second {
				delay = serverWait
			}
			log.Printf("[AI] %s 第 %d/%d 次重试，等待 %s", c.model, attempt, c.maxRetries, delay)
			select {
			case <-ctx.Done():
				return Result{}, ctx.Err()
			case <-time.After(delay):
			}
		}

		res, retryable, retryAfter, err := c.doChat(ctx, body)
		if err == nil {
			return res, nil
		}
		lastErr, serverWait = err, retryAfter
		if !retryable {
			return Result{}, err
		}
	}
	return Result{}, lastErr
}

// doChat 单次请求；retryable 表示该错误是否值得重试，retryAfter 为服务端建议等待时间
func (c *chatClient) doChat(ctx context.Context, body []byte) (res Result, retryable bool, retryAfter time.Duration, err error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return Result{}, false, 0, fmt.Errorf("ai: 创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		// 上下文取消/超时不再重试
		if ctx.Err() != nil {
			return Result{}, false, 0, ctx.Err()
		}
		return Result{}, true, 0, fmt.Errorf("ai: 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, true, 0, fmt.Errorf("ai: 读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errDesc := deepSeekError[resp.StatusCode]
		if errDesc == "" {
			errDesc = fmt.Sprintf("未知错误 (HTTP %d)", resp.StatusCode)
		}
		retry := resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode == http.StatusInternalServerError ||
			resp.StatusCode == http.StatusBadGateway ||
			resp.StatusCode == http.StatusServiceUnavailable ||
			resp.StatusCode == http.StatusGatewayTimeout
		var wait time.Duration
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			if secs, e := strconv.Atoi(ra); e == nil && secs > 0 {
				wait = time.Duration(secs) * time.Second
			}
		}
		return Result{}, retry, wait, fmt.Errorf("ai: %s, body=%s", errDesc, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return Result{}, false, 0, fmt.Errorf("ai: 解析响应失败: %w", err)
	}

	if chatResp.Error != nil {
		return Result{}, false, 0, fmt.Errorf("ai: API 错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return Result{}, false, 0, fmt.Errorf("ai: 响应中无 choices 字段")
	}

	choice := chatResp.Choices[0]
	res = Result{Content: choice.Message.Content, FinishReason: choice.FinishReason}
	if chatResp.Usage != nil {
		res.CompletionTokens = chatResp.Usage.CompletionTokens
	}
	// 截断是可诊断问题：finish_reason=length 时内容必然不完整，先留痕便于排查
	if res.Truncated() {
		log.Printf("[AI] %s 输出被截断（finish_reason=length, completion_tokens=%d）", c.model, res.CompletionTokens)
	}

	return res, false, 0, nil
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
			// 超时与重试次数改为读取配置（此前写死 600s 且忽略 max_retries）
			timeout := time.Duration(m.Timeout) * time.Second
			if timeout <= 0 {
				timeout = 60 * time.Second
			}
			switch m.Provider {
			case "deepseek":
				return NewDeepSeekWithOptions(m.APIKey, m.Endpoint, timeout, m.MaxRetries), nil
			case "qwen":
				return NewQwenWithOptions(m.APIKey, m.Endpoint, timeout, m.MaxRetries), nil
			default:
				return NewClientWithOptions(m.Endpoint, m.APIKey, m.Provider, timeout, m.MaxRetries), nil
			}
		}
	}
	return nil, fmt.Errorf("ai: 未找到 capability=%s 的已启用模型", capability)
}
