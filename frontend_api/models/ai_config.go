package models

import "time"

// ═══════════════════════════════════════════════
//  AI 模型供应商定义与配置
// ═══════════════════════════════════════════════

// ModelProvider AI 模型供应商枚举
type ModelProvider string

const (
	ModelDeepSeek ModelProvider = "deepseek"
	ModelChatGPT  ModelProvider = "chatgpt"
	ModelCodex    ModelProvider = "codex"
	ModelGemini   ModelProvider = "gemini"
	ModelClaude   ModelProvider = "claude"
	ModelMock     ModelProvider = "mock" // 本地 Mock 模式
)

// AllProviders 全部已注册的模型供应商列表
var AllProviders = []ModelProvider{
	ModelMock,
	ModelDeepSeek,
	ModelChatGPT,
	ModelCodex,
	ModelGemini,
	ModelClaude,
}

// ═══════════════════════════════════════════════
//  模型实例配置
// ═══════════════════════════════════════════════

// AIModelConfig 单个 AI 模型的配置
type AIModelConfig struct {
	Provider   ModelProvider `json:"provider" yaml:"provider"`       // 供应商标识
	Name       string        `json:"name" yaml:"name"`               // 展示名称（如 "DeepSeek V3"）
	APIKey     string        `json:"api_key" yaml:"api_key"`         // API 密钥
	Endpoint   string        `json:"endpoint" yaml:"endpoint"`       // 请求地址
	Enabled    bool          `json:"enabled" yaml:"enabled"`         // 是否启用
	MaxRetries int           `json:"max_retries" yaml:"max_retries"` // 最大重试次数
	Timeout    time.Duration `json:"timeout" yaml:"timeout"`         // 请求超时
	Priority   int           `json:"priority" yaml:"priority"`       // 优先级（数字越小越优先）
}

// ═══════════════════════════════════════════════
//  请求状态跟踪
// ═══════════════════════════════════════════════

// RequestStatus 单次请求的状态
type RequestStatus string

const (
	StatusPending    RequestStatus = "pending"     // 等待执行
	StatusRunning    RequestStatus = "running"     // 执行中
	StatusSuccess    RequestStatus = "success"     // 成功
	StatusFail       RequestStatus = "fail"        // 失败
	StatusRetry      RequestStatus = "retry"       // 重试中
	StatusFailover   RequestStatus = "failover"    // 已切换至备用模型
	StatusMaxRetries RequestStatus = "max_retries" // 超过最大重试
)

// AIRequestResult 单次 AI 请求的结果记录
type AIRequestResult struct {
	Provider ModelProvider `json:"provider"` // 实际执行的模型供应商
	Status   RequestStatus `json:"status"`   // 请求状态
	Duration time.Duration `json:"duration"` // 耗时

	// 生成结果
	ModelCode   string `json:"model_code"`   // 生成的模板代码
	PreviewHTML string `json:"preview_html"` // 预览 HTML

	// 评分
	Score      int              `json:"score"`
	Dimensions []ScoreDimension `json:"dimensions"`

	// 错误信息
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`

	// 时间戳
	StartedAt time.Time  `json:"started_at"`
	DoneAt    *time.Time `json:"done_at,omitempty"`
}

// ═══════════════════════════════════════════════
//  模型管理器 — 统一管理所有 AI 模型
// ═══════════════════════════════════════════════

// AIModelManager 全局 AI 模型管理器
type AIModelManager struct {
	Models     map[ModelProvider]*AIModelConfig `json:"models"`
	ActiveID   string                           `json:"active_id"`   // 当前激活的 session ID
	RetryCount int                              `json:"retry_count"` // 全局重试计数
}

// NewAIModelManager 创建默认模型管理器（预填充 Mock 配置）
func NewAIModelManager() *AIModelManager {
	mgr := &AIModelManager{
		Models: make(map[ModelProvider]*AIModelConfig),
	}

	// 注册 Mock 模型（本地调试用）
	mgr.Register(ModelMock, &AIModelConfig{
		Provider:   ModelMock,
		Name:       "Local Mock",
		APIKey:     "",
		Endpoint:   "",
		Enabled:    true,
		MaxRetries: 0,
		Timeout:    0,
		Priority:   0,
	})

	return mgr
}

// Register 注册或更新一个模型配置
func (m *AIModelManager) Register(provider ModelProvider, cfg *AIModelConfig) {
	m.Models[provider] = cfg
}

// GetConfig 获取指定模型的配置
func (m *AIModelManager) GetConfig(provider ModelProvider) *AIModelConfig {
	return m.Models[provider]
}

// GetAvailableModels 获取所有已启用的模型，按优先级排序
func (m *AIModelManager) GetAvailableModels() []*AIModelConfig {
	// 简单返回所有启用的
	var result []*AIModelConfig
	for _, cfg := range m.Models {
		if cfg.Enabled {
			result = append(result, cfg)
		}
	}
	return result
}

// Failover 失败转移：返回下一个可用的模型配置
func (m *AIModelManager) Failover(failedProvider ModelProvider) *AIModelConfig {
	for _, cfg := range m.GetAvailableModels() {
		if cfg.Provider != failedProvider {
			return cfg
		}
	}
	return nil
}

// ═══════════════════════════════════════════════
//  AI 请求的输入 / 输出结构（对齐 generate 接口）
// ═══════════════════════════════════════════════

// AIGenerateRequest 传递给 AI 模型的统一请求结构
type AIGenerateRequest struct {
	Designs   []DesignItem `json:"designs"`
	Framework string       `json:"framework"`
	Quality   int          `json:"quality"`
}

// AIGenerateResponse AI 模型生成的统一响应
type AIGenerateResponse struct {
	Code       string           `json:"code"`       // Flutter / React / Vue 代码
	Preview    string           `json:"preview"`    // 预览 HTML
	Score      int              `json:"score"`      // 综合评分
	Dimensions []ScoreDimension `json:"dimensions"` // 各维度评分
	ID         int              `json:"id"`         // 模板 ID
}
