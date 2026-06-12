package handler

import (
	"frontend_api/pkg/cos"
	"frontend_api/pkg/logger"
	"frontend_api/utils"

	"github.com/gin-gonic/gin"
)

type uploadRequest struct {
	Image    string `json:"image" binding:"required"`
	Filename string `json:"filename"`
}

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

// Upload 上传图片到腾讯云 COS
// POST /api/v1/upload
func (h *UploadHandler) Upload(c *gin.Context) {
	var req uploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请提供图片数据")
		return
	}

	// 解码 base64
	data, err := cos.DecodeBase64Image(req.Image)
	if err != nil {
		utils.BadRequest(c, "图片数据格式错误")
		return
	}

	filename := req.Filename
	if filename == "" {
		filename = "image.png"
	}

	h.log.Infof("收到上传图片: %s (%d bytes)", filename, len(data))

	// 上传到 COS
	url, err := cos.UploadBytes(data, filename)
	if err != nil {
		h.log.Errorf("上传到 COS 失败: %v", err)
		utils.InternalError(c, "上传失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"url":      url,
		"filename": filename,
		"size":     len(data),
	}, "上传成功")
}
