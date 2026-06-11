package services

import (
	"fmt"
	"frontend_api/config"
	"frontend_api/mockdata"
	"frontend_api/models"
	"log"
	"time"
)

// ═══════════════════════════════════════════════
//  AI 模型调用服务
//  职责：配置管理 + 请求执行 + 重试/故障转移 + 日志记录
// ═══════════════════════════════════════════════

// AIService AI 模型调用服务
type AIService struct {
	Manager *models.AIModelManager
	Logger  *AILogger
}

// NewAIService 创建 AI 服务实例（自动从 config.yaml 加载模型配置）
func NewAIService() *AIService {
	mgr := models.NewAIModelManager()
	// 从 config.yaml 加载模型配置（覆盖默认 Mock）
	config.LoadModelsIntoManager(mgr)
	return &AIService{
		Manager: mgr,
		Logger:  GetAILogger(),
	}
}

// Call 执行一次 AI 模型调用（含重试 + 故障转移）
//
// 参数:
//   - provider:  指定模型供应商（传空则使用最高优先级启用的模型）
//   - request:   统一请求参数
//   - useStream: 是否启用流式输出（为 true 时通过 Broker 推送）
//
// 返回:
//   - result: 请求结果（含状态、耗时、数据）
func (s *AIService) Call(provider models.ModelProvider, request *models.AIGenerateRequest, useStream bool) *models.AIRequestResult {
	// 选择模型
	cfg := s.selectModel(provider)
	if cfg == nil {
		return &models.AIRequestResult{
			Provider:     models.ModelMock,
			Status:       models.StatusFail,
			ErrorMessage: "没有可用的 AI 模型配置",
			StartedAt:    time.Now(),
		}
	}

	result := &models.AIRequestResult{
		Provider:  cfg.Provider,
		Status:    models.StatusPending,
		StartedAt: time.Now(),
	}

	// 执行调用（含重试）
	result = s.executeWithRetry(cfg, request, useStream, result)

	return result
}

// selectModel 选择要使用的模型（指定 / 默认 / Mock 兜底）
func (s *AIService) selectModel(preferred models.ModelProvider) *models.AIModelConfig {
	if preferred != "" && preferred != models.ModelMock {
		if cfg := s.Manager.GetConfig(preferred); cfg != nil && cfg.Enabled {
			return cfg
		}
	}

	// 按优先级选择
	avail := s.Manager.GetAvailableModels()
	for _, cfg := range avail {
		if cfg.Provider != models.ModelMock {
			return cfg
		}
	}

	// 兜底使用 Mock
	return s.Manager.GetConfig(models.ModelMock)
}

// executeWithRetry 执行并支持重试
func (s *AIService) executeWithRetry(cfg *models.AIModelConfig, req *models.AIGenerateRequest, useStream bool, prevResult *models.AIRequestResult) *models.AIRequestResult {
	result := prevResult
	if result == nil {
		result = &models.AIRequestResult{
			Provider:  cfg.Provider,
			Status:    models.StatusPending,
			StartedAt: time.Now(),
		}
	}

	maxRetries := cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 1 // 至少执行一次
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			result.Status = models.StatusRetry
			s.Logger.Log(cfg.Provider, models.StatusRetry, 0,
				fmt.Sprintf("第 %d 次重试 (共 %d 次)", attempt, maxRetries))
		}

		result.Status = models.StatusRunning
		r, err := s.invokeModel(cfg, req, useStream)
		result.Duration = time.Since(result.StartedAt)

		if err == nil {
			result.Status = models.StatusSuccess
			result.ModelCode = r.Code
			result.PreviewHTML = r.Preview
			result.Score = r.Score
			result.Dimensions = r.Dimensions
			now := time.Now()
			result.DoneAt = &now

			s.Logger.LogResult(result, "代码生成成功")
			return result
		}

		// 失败
		result.ErrorCode = err.Code
		result.ErrorMessage = err.Message
		result.Status = models.StatusFail

		// 还有重试次数
		if attempt < maxRetries-1 {
			s.Logger.Log(cfg.Provider, models.StatusRetry, result.Duration,
				fmt.Sprintf("%s — 即将重试", err.Message))
			continue
		}
	}

	// 全部重试失败 → 故障转移
	failoverCfg := s.Manager.Failover(cfg.Provider)
	if failoverCfg != nil {
		result.Status = models.StatusFailover
		result.Provider = failoverCfg.Provider
		s.Logger.Log(cfg.Provider, models.StatusFailover, 0,
			fmt.Sprintf("已切换至 %s", failoverCfg.Name))

		// 递归调用备用模型（不限制重试次数以免无限递归）
		return s.executeWithoutRetry(failoverCfg, req, useStream, result)
	}

	result.Status = models.StatusMaxRetries
	now := time.Now()
	result.DoneAt = &now
	s.Logger.Log(cfg.Provider, models.StatusMaxRetries, result.Duration, "所有模型均失败")
	return result
}

// executeWithoutRetry 无重试的调用（故障转移用，只执行一次）
func (s *AIService) executeWithoutRetry(cfg *models.AIModelConfig, req *models.AIGenerateRequest, useStream bool, fallbackFrom *models.AIRequestResult) *models.AIRequestResult {
	result := &models.AIRequestResult{
		Provider:  cfg.Provider,
		Status:    models.StatusRunning,
		StartedAt: time.Now(),
	}

	r, err := s.invokeModel(cfg, req, useStream)
	result.Duration = time.Since(result.StartedAt)

	if err == nil {
		result.Status = models.StatusSuccess
		result.ModelCode = r.Code
		result.PreviewHTML = r.Preview
		result.Score = r.Score
		result.Dimensions = r.Dimensions
		now := time.Now()
		result.DoneAt = &now
		s.Logger.LogResult(result, "故障转移成功")
		return result
	}

	result.Status = models.StatusFail
	result.ErrorCode = err.Code
	result.ErrorMessage = err.Message
	now := time.Now()
	result.DoneAt = &now
	s.Logger.Log(cfg.Provider, models.StatusFail, result.Duration,
		fmt.Sprintf("故障转移也失败: %s", err.Message))
	return result
}

// ═══════════════════════════════════════════════
//  模型调用接口（当前使用 Mock，后续替换为真实 API）
// ═══════════════════════════════════════════════

// ModelError 模型调用错误
type ModelError struct {
	Code    string
	Message string
}

func (e *ModelError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// invokeModel 执行模型调用（当前为 Mock 实现）
// TODO: 根据 cfg.Provider 路由到不同的 AI 模型 SDK
func (s *AIService) invokeModel(cfg *models.AIModelConfig, req *models.AIGenerateRequest, useStream bool) (*models.AIGenerateResponse, *ModelError) {
	log.Printf("[AIService] 调用模型: %s (%s), framework=%s, quality=%d, useStream=%v",
		cfg.Name, cfg.Provider, req.Framework, req.Quality, useStream)

	switch cfg.Provider {
	case models.ModelMock:
		return s.mockInvoke(req, useStream)

	// ═══ 后续接入真实 AI 模型 ═══
	// case models.ModelDeepSeek:
	//   return s.deepSeekInvoke(req, useStream)
	// case models.ModelChatGPT:
	//   return s.chatGPTInvoke(req, useStream)
	// case models.ModelCodex:
	//   return s.codexInvoke(req, useStream)
	// case models.ModelGemini:
	//   return s.geminiInvoke(req, useStream)

	default:
		// 未知模型 → 回退 Mock
		log.Printf("[AIService] 未知模型 %s，回退到 Mock", cfg.Provider)
		return s.mockInvoke(req, useStream)
	}
}

// mockInvoke Mock 模式调用
func (s *AIService) mockInvoke(req *models.AIGenerateRequest, useStream bool) (*models.AIGenerateResponse, *ModelError) {
	// 评分数据
	score := mockdata.GetMockScore(req.Quality)
	dims := []models.ScoreDimension{
		{Name: "视觉还原度", Score: mockdata.CalcDimScore(req.Quality, 0), Icon: "fas fa-palette"},
		{Name: "代码质量", Score: mockdata.CalcDimScore(req.Quality, -3), Icon: "fas fa-code"},
		{Name: "响应式设计", Score: mockdata.CalcDimScore(req.Quality, -5), Icon: "fas fa-mobile-alt"},
		{Name: "性能优化", Score: mockdata.CalcDimScore(req.Quality, -2), Icon: "fas fa-tachometer-alt"},
	}

	// 如果是流式模式，通过 Broker 推送
	if useStream {
		broker := GetBroker()
		for chunk := range mockdata.StreamTemplateCode() {
			broker.Publish(SSEEvent{Event: "message", Data: chunk})
		}
		previewHTML := mockdata.GetTemplatePreviewHTML()
		broker.Publish(SSEEvent{Event: "preview", Data: previewHTML})

		scoreJSON := fmt.Sprintf(
			`{"score":%d,"dimensions":[{"name":"%s","score":%d,"icon":"%s"},{"name":"%s","score":%d,"icon":"%s"},{"name":"%s","score":%d,"icon":"%s"},{"name":"%s","score":%d,"icon":"%s"}]}`,
			score,
			dims[0].Name, dims[0].Score, dims[0].Icon,
			dims[1].Name, dims[1].Score, dims[1].Icon,
			dims[2].Name, dims[2].Score, dims[2].Icon,
			dims[3].Name, dims[3].Score, dims[3].Icon,
		)
		broker.Publish(SSEEvent{Event: "score", Data: scoreJSON})
		broker.Publish(SSEEvent{Event: "done", Data: `{"id":1}`})

		log.Printf("[AIService][Mock] 流式推送完成")
	}

	return &models.AIGenerateResponse{
		Code:       mockdata.GetTemplatePreviewHTML(),
		Preview:    mockdata.GetTemplatePreviewHTML(),
		Score:      score,
		Dimensions: dims,
		ID:         1,
	}, nil
}
