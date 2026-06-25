package handler

import (
	"fmt"
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"frontend_api/utils"
	"log"
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ═══════════════════════════════════════════════
//  AdminHandler — 后台管理相关接口
//  需要管理权限方可调用
// ═══════════════════════════════════════════════

// AdminHandler 管理员处理器
type AdminHandler struct{}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// CreateAccountRequest 创建账号请求
type CreateAccountRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Level    int    `json:"level" binding:"required,oneof=1 2 3"`
}

// 等级对应的积分映射
var levelCreditsMap = map[int]int{
	1: 100, // vip1
	2: 150, // vip2
	3: 200, // vip3
}

// CreateAccount 管理员创建用户账号
// POST /api/v1/admin/create-account
func (h *AdminHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
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

	// 根据等级确定初始积分
	initialCredits := levelCreditsMap[req.Level]

	// 随机生成昵称
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	length := 3 + rand.Intn(3) // 3~5
	name := make([]byte, length)
	for i := range name {
		name[i] = chars[rand.Intn(len(chars))]
	}

	// 生成随机头像（DiceBear）
	avatarStyle := "avataaars"
	avatar := fmt.Sprintf("https://api.dicebear.com/7.x/%s/svg?seed=%s", avatarStyle, req.Email)

	// 创建用户
	newUser := model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: string(name),
		Avatar:   avatar,
		Status:   "active",
		Credits:  initialCredits,
		Level:    req.Level,
	}

	if err := db.Create(&newUser).Error; err != nil {
		utils.InternalError(c, "用户创建失败: "+err.Error())
		return
	}

	log.Printf("[Admin] 创建账号成功: email=%s, level=%d, credits=%d", req.Email, req.Level, initialCredits)

	utils.Success(c, gin.H{
		"id":       newUser.ID,
		"email":    newUser.Email,
		"nickname": newUser.Nickname,
		"level":    newUser.Level,
		"credits":  newUser.Credits,
		"avatar":   newUser.Avatar,
	}, "账号创建成功")
}
