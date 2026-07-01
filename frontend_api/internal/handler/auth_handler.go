package handler

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"frontend_api/internal/model"
	"frontend_api/pkg/email"
	"frontend_api/pkg/mysql"
	"frontend_api/utils"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════
//  AuthHandler — 认证相关（从旧 controllers/auth.go 迁移）
//  保持原有接口兼容
// ═══════════════════════════════════════════════

// AuthHandler 认证处理器
type AuthHandler struct{}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" binding:"required"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaID   string `json:"captcha_id" binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaID    string `json:"captcha_id"`
	CaptchaImage string `json:"captcha_image"`
}

// TemplateRequest 模板请求参数
type TemplateRequest struct {
	Template int `form:"template" json:"template"`
}

// ForgotPasswordRequest 忘记密码请求参数
type ForgotPasswordRequest struct {
	Email         string `json:"email" binding:"required"`
	EmailProvider string `json:"email_provider"` // 前端可选透传，后端优先从邮箱域名自动识别
}

// ResetPasswordRequest 重置密码请求参数
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// Captcha 获取图形验证码
func (h *AuthHandler) Captcha(c *gin.Context) {
	captchaID, captchaImage, _ := utils.GenerateCaptcha()
	utils.Success(c, CaptchaResponse{
		CaptchaID:    captchaID,
		CaptchaImage: captchaImage,
	}, "验证码获取成功")
}

// Logout 退出登录
func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT 无状态，无需后端处理
	utils.Success(c, gin.H{}, "已退出登录")
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	// 统一转小写
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	db := mysql.GetDB()

	// 检查邮箱是否已注册
	var existingUser model.User
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.Error(c, 409, "该邮箱已被注册")
		return
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "密码处理失败")
		return
	}

	// 创建用户
	newUser := model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
	}
	if req.Nickname == "" {
		// 随机生成 3~5 位昵称
		chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
		length := 3 + rand.Intn(3) // 3~5
		name := make([]byte, length)
		for i := range name {
			name[i] = chars[rand.Intn(len(chars))]
		}
		newUser.Nickname = string(name)
	}

	// 生成随机头像（DiceBear，基于 email seed 保持一致）
	avatarStyle := "avataaars"
	newUser.Avatar = fmt.Sprintf("https://api.dicebear.com/7.x/%s/svg?seed=%s", avatarStyle, req.Email)

	if err := db.Create(&newUser).Error; err != nil {
		utils.InternalError(c, "用户创建失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id":       newUser.ID,
		"email":    newUser.Email,
		"nickname": newUser.Nickname,
		"avatar":   newUser.Avatar,
	}, "注册成功")
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	// 统一转小写
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// 校验验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
		utils.BadRequest(c, "验证码错误或已过期")
		return
	}

	db := mysql.GetDB()

	// 查找用户（通过 email）
	var user model.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Unauthorized(c, "该邮箱未注册")
		} else {
			utils.InternalError(c, "查询用户失败")
		}
		return
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Unauthorized(c, "密码错误")
		return
	}

	// 生成 JWT
	configJWT := mysql.GetConfig() // 通过 mysql 包读取配置
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(configJWT))
	if err != nil {
		utils.InternalError(c, "Token 生成失败")
		return
	}

	utils.Success(c, gin.H{
		"token":        tokenString,
		"id":           user.ID,
		"email":        user.Email,
		"nickname":     user.Nickname,
		"avatar":       user.Avatar,
		"credits":      user.Credits,
		"credits_used": user.CreditsUsed,
	}, "登录成功")
}

// ForgotPassword 忘记密码 - 生成重置令牌并发送邮件
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// 从邮箱域名自动识别服务商
	provider, err := email.ValidateEmail(req.Email)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	db := mysql.GetDB()

	var user model.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.Error(c, 404, "该邮箱未注册")
		return
	}
	if user.Status == "disabled" {
		utils.Error(c, 403, "账号已被禁用")
		return
	}

	tokenBytes := make([]byte, 32)
	if _, err := cryptorand.Read(tokenBytes); err != nil {
		utils.InternalError(c, "令牌生成失败")
		return
	}
	token := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(30 * time.Minute)

	if err := db.Model(&user).Updates(map[string]interface{}{
		"reset_token":            token,
		"reset_token_expires_at": expiresAt,
		"reset_token_used_at":    nil,
	}).Error; err != nil {
		utils.InternalError(c, "令牌存储失败: "+err.Error())
		return
	}

	sender, err := email.GetSender(provider)
	if err != nil {
		utils.InternalError(c, "邮件服务配置错误: "+err.Error())
		return
	}
	if err := sender.SendResetEmail(req.Email, token); err != nil {
		utils.InternalError(c, "邮件发送失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{}, "重置邮件已发送，请检查邮箱")
}

// ResetPassword 重置密码
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	db := mysql.GetDB()

	var user model.User
	if err := db.Where("reset_token = ?", req.Token).First(&user).Error; err != nil {
		utils.Error(c, 400, "重置链接无效或已使用")
		return
	}

	if user.ResetTokenUsedAt != nil {
		utils.Error(c, 400, "重置链接已被使用")
		return
	}
	if user.ResetTokenExpiresAt == nil || time.Now().After(*user.ResetTokenExpiresAt) {
		utils.Error(c, 400, "重置链接已过期，请重新申请")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "密码处理失败")
		return
	}

	now := time.Now()
	if err := db.Model(&user).Updates(map[string]interface{}{
		"password":            string(hashedPassword),
		"reset_token_used_at": &now,
	}).Error; err != nil {
		utils.InternalError(c, "密码更新失败")
		return
	}

	utils.Success(c, gin.H{}, "密码重置成功，请重新登录")
}

// Template 获取模板数据
// GET /api/template?template=1
func (h *AuthHandler) Template(c *gin.Context) {
	var req TemplateRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, "参数解析失败: "+err.Error())
		return
	}

	if req.Template == 0 {
		utils.BadRequest(c, "缺少必要参数: template")
		return
	}

	// 从 mockdata 获取
	templateID := req.Template

	// 直接引用 mockdata（保持兼容）
	_ = templateID

	utils.Error(c, 404, "模板获取功能已迁移，请使用 POST /api/v1/generate-ui")
}

// Helper: 从请求中提取 user_id（由 AuthMiddleware 注入）
func getUserID(c *gin.Context) uint {
	id, _ := c.Get("user_id")
	if uid, ok := id.(float64); ok {
		return uint(uid)
	}
	if uid, ok := id.(uint); ok {
		return uid
	}
	return 0
}
