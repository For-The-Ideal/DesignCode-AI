package controllers

import (
	"fmt"
	"frontend_api/mockdata"
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
//  2. 调用模型生成目标业务数据（当前为 mock）
//  3. 通过 Broker 将数据分批次推送给已连接的 SSE 客户端
//  4. 按原协议返回最终结果
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

	// TODO: 接入 AI 服务 (OpenAI / Claude / 自定义模型)
	// aiService := services.NewAIService()
	// result, err := aiService.GenerateCode(req.Designs, req.Framework, req.Quality)

	// ═══ 评分数据生成 ═══
	score := mockdata.GetMockScore(req.Quality)
	dims := []models.ScoreDimension{
		{Name: "视觉还原度", Score: mockdata.CalcDimScore(req.Quality, 0), Icon: "fas fa-palette"},
		{Name: "代码质量", Score: mockdata.CalcDimScore(req.Quality, -3), Icon: "fas fa-code"},
		{Name: "响应式设计", Score: mockdata.CalcDimScore(req.Quality, -5), Icon: "fas fa-mobile-alt"},
		{Name: "性能优化", Score: mockdata.CalcDimScore(req.Quality, -2), Icon: "fas fa-tachometer-alt"},
	}

	// ═══ 通过 Broker → SSE 逐片推送 ═══
	broker := services.GetBroker()

	// 流式推送代码（template_1.go 通过 goroutine channel 逐块产出）
	for chunk := range mockdata.StreamTemplateCode() {
		broker.Publish(services.SSEEvent{Event: "message", Data: chunk})
	}

	// 推送预览（template_1.go 的真实商品卡片 HTML）
	previewHTML := mockdata.GetTemplatePreviewHTML()
	broker.Publish(services.SSEEvent{Event: "preview", Data: previewHTML})

	// 推送评分
	scoreJSON := fmt.Sprintf(
		`{"score":%d,"dimensions":[{"name":"%s","score":%d,"icon":"%s"},{"name":"%s","score":%d,"icon":"%s"},{"name":"%s","score":%d,"icon":"%s"},{"name":"%s","score":%d,"icon":"%s"}]}`,
		score,
		dims[0].Name, dims[0].Score, dims[0].Icon,
		dims[1].Name, dims[1].Score, dims[1].Icon,
		dims[2].Name, dims[2].Score, dims[2].Icon,
		dims[3].Name, dims[3].Score, dims[3].Icon,
	)
	broker.Publish(services.SSEEvent{Event: "score", Data: scoreJSON})

	// 推送完成信号
	broker.Publish(services.SSEEvent{Event: "done", Data: "ok"})

	log.Printf("[Generate] 已通过 Broker 推送所有事件")

	// 仅返回生成状态，不做业务数据响应（数据走 SSE）
	utils.Success(c, gin.H{"status": "generating"}, "代码生成中，请通过 SSE 接收实时数据")
}
