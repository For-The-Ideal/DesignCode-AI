package controllers

import (
	"frontend_api/models"
	"frontend_api/services"
	"frontend_api/utils"
	"log"

	"github.com/gin-gonic/gin"
)

// GenerateController AI 代码生成控制器
type GenerateController struct{}

// Generate 处理代码生成请求（阻塞式）
//
// 流程：
//  1. 参数校验 & 结构化解析
//  2. 通过 AIService 调用 AI 模型（当前为 Mock，后续可切换 DeepSeek / ChatGPT / Codex / Gemini / Claude）
//  3. 通过 Broker 将数据分批次推送给已连接的 SSE 客户端
//  4. 返回生成状态
func (g *GenerateController) Generate(c *gin.Context) {
	var req models.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Generate] 参数校验失败: %v", err)
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	// 参数合法性校验
	if len(req.Designs) == 0 {
		utils.BadRequest(c, "至少需要上传一张设计稿")
		return
	}
	if req.Framework == "" {
		utils.BadRequest(c, "请选择目标框架（flutter/react/vue）")
		return
	}
	if req.Quality < 1 || req.Quality > 100 {
		utils.BadRequest(c, "质量参数需在 1-100 之间")
		return
	}

	log.Printf("[Generate] 收到请求: framework=%s, quality=%d, designs=%d",
		req.Framework, req.Quality, len(req.Designs))

	// ═══ 通过 AIService 调用 AI 模型（含重试 / 故障转移 / 日志） ═══
	aiService := services.NewAIService()

	// 转换为统一的 AI 请求结构
	aiReq := &models.AIGenerateRequest{
		Designs:   req.Designs,
		Framework: req.Framework,
		Quality:   req.Quality,
	}

	// useStream=true → AIService 内部通过 Broker 逐片推送
	result := aiService.Call("", aiReq, true)

	if result.Status != models.StatusSuccess {
		log.Printf("[Generate] AI 调用失败: [%s] %s", result.ErrorCode, result.ErrorMessage)
		utils.Fail(c, "AI 生成失败: "+result.ErrorMessage)
		return
	}

	log.Printf("[Generate] AI 调用成功，模型=%s，耗时=%v，评分=%d",
		result.Provider, result.Duration, result.Score)

	// 仅返回生成状态，业务数据通过 SSE 实时推送
	utils.Success(c, gin.H{"status": "generating"}, "代码生成中")
}
