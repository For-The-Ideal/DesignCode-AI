package handler

import (
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"frontend_api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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
	Username string `json:"username" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
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

// Captcha 获取图形验证码
func (h *AuthHandler) Captcha(c *gin.Context) {
	captchaID, captchaImage, _ := utils.GenerateCaptcha()
	utils.Success(c, CaptchaResponse{
		CaptchaID:    captchaID,
		CaptchaImage: captchaImage,
	}, "验证码获取成功")
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	db := mysql.GetDB()

	// 检查用户是否已存在
	var existingUser model.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.Error(c, 409, "该用户名已被注册")
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
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
	}
	if req.Nickname == "" {
		newUser.Nickname = req.Username
	}

	if err := db.Create(&newUser).Error; err != nil {
		utils.InternalError(c, "用户创建失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id":       newUser.ID,
		"username": newUser.Username,
	}, "注册成功")
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	// 校验验证码
	if !utils.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
		utils.BadRequest(c, "验证码错误或已过期")
		return
	}

	db := mysql.GetDB()

	// 查找用户
	var user model.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 生成 JWT
	configJWT := mysql.GetConfig() // 通过 mysql 包读取配置
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(configJWT))
	if err != nil {
		utils.InternalError(c, "Token 生成失败")
		return
	}

	utils.Success(c, gin.H{
		"token":    tokenString,
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
	}, "登录成功")
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
