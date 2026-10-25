package handler

import (
	"frontend_api/internal/repository"
	"frontend_api/utils"

	"github.com/gin-gonic/gin"
)

// ═══════════════════════════════════════════════
//  ConfigHandler — 生成配置字典（框架/平台/选项/高级/组件库）
// ═══════════════════════════════════════════════

// ConfigHandler 配置处理器
type ConfigHandler struct {
	genConfigRepo *repository.GenConfigRepository
}

// NewConfigHandler 创建配置处理器
func NewConfigHandler(repo *repository.GenConfigRepository) *ConfigHandler {
	return &ConfigHandler{genConfigRepo: repo}
}

// GetOptions 返回分组后的生成配置字典（公开接口，登录前即可加载）
//
//	GET /api/v1/config/options
func (h *ConfigHandler) GetOptions(c *gin.Context) {
	opts, err := h.genConfigRepo.GetOptions()
	if err != nil {
		utils.InternalError(c, "获取生成配置失败: "+err.Error())
		return
	}
	utils.Success(c, opts, "success")
}
