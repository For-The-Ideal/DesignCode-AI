package handler

import (
	"frontend_api/internal/skills/vision"
	"frontend_api/pkg/logger"
	"frontend_api/utils"

	"github.com/gin-gonic/gin"
)

// ═══════════════════════════════════════════════
//  UiProjectHandler — UI 图转代码项目相关
//  当前仅提供顶层区域识别（进编辑器时按需调用）
// ═══════════════════════════════════════════════

type regionAnalyzeRequest struct {
	ImageURL string `json:"imageUrl" binding:"required"`
}

// UiProjectHandler 项目处理器
type UiProjectHandler struct {
	log         *logger.Logger
	regionSkill *vision.RegionSkill
}

// NewUiProjectHandler 创建项目处理器
func NewUiProjectHandler(regionSkill *vision.RegionSkill) *UiProjectHandler {
	return &UiProjectHandler{
		log:         logger.NewLogger("ui-project"),
		regionSkill: regionSkill,
	}
}

// AnalyzeRegions 顶层区域识别
// POST /api/v1/ui-project/analyze  body: { imageUrl }
// 返回：{ code:200, data: [{name,summary,styleHints,behavior,childrenHint,layoutHint}] }（data 直接是区域数组，不包 regions 层）
func (h *UiProjectHandler) AnalyzeRegions(c *gin.Context) {
	var req regionAnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请提供图片地址 imageUrl")
		return
	}

	if h.regionSkill == nil {
		h.log.Warn("顶层识别请求到达，但视觉AI未配置")
		utils.InternalError(c, "顶层识别暂不可用：视觉AI未配置")
		return
	}

	out, err := h.regionSkill.Execute(c.Request.Context(), vision.RegionInput{ImageURL: req.ImageURL})
	if err != nil {
		h.log.Errorf("顶层区域识别失败: %v", err)
		utils.InternalError(c, "顶层识别失败: "+err.Error())
		return
	}

	regions, _ := out.(vision.RegionOutput)
	utils.Success(c, regions.Regions, "识别成功")
}
