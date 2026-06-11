package mysql

import (
	"frontend_api/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
// dsn 从 config.AppConfig.MySQL.DSN 传入
func InitDB(dsn string) {
	if dsn == "" {
		log.Println("[MySQL] DSN 为空，跳过数据库初始化")
		return
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[MySQL] 数据库链接失败: %v", err)
	}
	log.Println("[MySQL] 数据库连接成功")
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

// GetConfig 获取 JWT 密钥等配置项
// 由 auth_handler 调用，保持与旧 controllers/auth.go 兼容
func GetConfig() string {
	if config.AppConfig == nil {
		log.Println("[MySQL] AppConfig 未初始化，返回空密钥")
		return ""
	}
	return config.AppConfig.Server.JWTSecret
}
