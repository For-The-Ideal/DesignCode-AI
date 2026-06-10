package controllers

import (
	"frontend_api/config"
	"frontend_api/models"
	"frontend_api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB // 注入 DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
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

type CaptchaResponse struct {
	CaptchaID    string `json:"captcha_id"`
	CaptchaImage string `json:"captcha_image"`
}

// Captcha 获取图形验证码
func (a *AuthController) Captcha(c *gin.Context) {
	captchaID, captchaImage, _ := utils.GenerateCaptcha()
	utils.Success(c, CaptchaResponse{
		CaptchaID:    captchaID,
		CaptchaImage: captchaImage,
	}, "验证码获取成功")
}

// Register 用户注册
func (a *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	// 1. 检查用户是否已存在
	var existingUser models.User
	if err := a.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.Error(c, 409, "该用户名已被注册")
		return
	}

	// 2. 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.InternalError(c, "密码处理失败")
		return
	}

	// 3. 创建用户
	newUser := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
	}
	if req.Nickname == "" {
		newUser.Nickname = req.Username
	}

	if err := a.DB.Create(&newUser).Error; err != nil {
		utils.InternalError(c, "用户创建失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id":       newUser.ID,
		"username": newUser.Username,
	}, "注册成功")
}

// Login 用户登录
func (a *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	// 0. 校验验证码（一次性消费）
	if !utils.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
		utils.BadRequest(c, "验证码错误或已过期")
		return
	}

	// 1. 查找用户
	var user models.User
	if err := a.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 2. 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Unauthorized(c, "用户名或密码错误")
		return
	}

	// 3. 生成 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.Server.JWTSecret))
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
