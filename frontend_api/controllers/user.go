package controllers

import (
	"frontend_api/models"
	"frontend_api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct{}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// GetProfile 获取当前登录用户信息
func (u *UserController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user *models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	utils.Success(c, user, "获取用户信息成功")
}

// ChangePassword 修改密码
func (u *UserController) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	// 查找用户
	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
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
	if err := utils.DB.Save(&user).Error; err != nil {
		utils.InternalError(c, "密码更新失败")
		return
	}

	utils.Success(c, nil, "密码修改成功")
}
