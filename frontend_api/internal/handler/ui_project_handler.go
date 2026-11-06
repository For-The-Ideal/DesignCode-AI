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

// regionDescribeRequest AI 描述请求体（与前端 api/uiProject.js describeRegion 一致）
type regionDescribeRequest struct {
	ImageURL string              `json:"imageUrl" binding:"required"`
	Tree     []vision.RegionNode `json:"tree"` // 整页区域树（所有顶层区域及全部子级）
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

// DescribeRegion 单区域 AI 描述（右侧面板「AI 描述」按钮）
// POST /api/v1/ui-project/describe
// body: { imageUrl, tree }
//   - tree：整页区域树（所有顶层区域及全部子级，含 nodeKey/name/summary/styleHints/layout/behavior/bbox/asset/isFocus，children 递归任意层级）
//
// 返回：{ code:200, data: { tree: [...] } }
//   - tree 与请求 tree 完全同构（nodeKey/name/styleHints/layout/behavior/bbox/asset/isFocus/children 原样），
//     仅 summary 被替换为 AI 生成的描述；前端按 nodeKey 全部回填
func (h *UiProjectHandler) DescribeRegion(c *gin.Context) {
	var req regionDescribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请提供图片地址 imageUrl")
		return
	}
	// 旧版请求（无 tree）不再可用：没有区域结构时 AI 只能凭空看图，生成结果不可控
	if len(req.Tree) == 0 {
		utils.BadRequest(c, "缺少区域树结构 tree，请刷新页面重试")
		return
	}

	if h.regionSkill == nil {
		h.log.Warn("区域描述请求到达，但视觉AI未配置")
		utils.InternalError(c, "AI 描述暂不可用：视觉AI未配置")
		return
	}

	out, err := h.regionSkill.Describe(c.Request.Context(), vision.RegionDescribeInput{
		ImageURL: req.ImageURL,
		Tree:     req.Tree,
	})
	if err != nil {
		h.log.Errorf("区域描述生成失败: %v", err)
		utils.InternalError(c, "AI 描述生成失败: "+err.Error())
		return
	}

	utils.Success(c, out, "生成成功")
}
