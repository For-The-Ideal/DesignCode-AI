package handler

import (
	"frontend_api/pkg/logger"
	"frontend_api/utils"

	"github.com/gin-gonic/gin"
)

// ═══════════════════════════════════════════════
//  UploadHandler — 图片上传（占位）
//  职责：接收上传图片 → 保存到腾讯云 COS → 返回 URL
// ═══════════════════════════════════════════════

// UploadHandler 上传处理器
type UploadHandler struct {
	log *logger.Logger
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		log: logger.NewLogger("upload"),
	}
}

// Upload 上传图片（当前 Mock：直接返回空结果）
// POST /api/v1/upload
//
// TODO: 接入腾讯云 COS SDK 进行实际上传
func (h *UploadHandler) Upload(c *gin.Context) {
	// 接收文件
	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequest(c, "请选择要上传的图片")
		return
	}

	h.log.Infof("收到上传文件: %s (%d bytes)", file.Filename, file.Size)

	// TODO: 上传到 COS
	// url, err := cos.Upload(file)
	// if err != nil { ... }

	utils.Success(c, gin.H{
		"url":      "", // TODO: 返回 COS URL
		"filename": file.Filename,
		"size":     file.Size,
	}, "上传成功（Mock）")
}
