package handler

import (
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"frontend_api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ═══════════════════════════════════════════════
//  UserHandler — 用户相关（从旧 controllers/user.go 迁移）
//  保持原有接口兼容
// ═══════════════════════════════════════════════

// UserHandler 用户处理器
type UserHandler struct{}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// GetUserInfo 获取当前登录用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user *model.User
	if err := mysql.GetDB().First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "用户不存在")
		return
	}
	utils.Success(c, user, "获取用户信息成功")
}

// UpdateUserPassword 修改密码
func (h *UserHandler) UpdateUserPassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	db := mysql.GetDB()

	// 查找用户
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 校验旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		utils.Error(c, http.StatusBadRequest, "当前密码错误")
		return
	}

	// 加密并保存新密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "密码处理失败")
		return
	}
	user.Password = string(hashed)
	if err := db.Save(&user).Error; err != nil {
		utils.InternalError(c, "密码更新失败")
		return
	}

	utils.Success(c, nil, "密码修改成功")
}
