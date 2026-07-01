package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitUploadRoutes 初始化上传相关路由
func InitUploadRoutes(v1 *gin.RouterGroup, uploadHandler *handler.UploadHandler) {
	// POST /api/v1/upload → 图片上传到 COS
	v1.POST("/upload", uploadHandler.Upload)
}
