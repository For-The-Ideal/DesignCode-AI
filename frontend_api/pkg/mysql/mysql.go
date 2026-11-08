package mysql

import (
	"frontend_api/config"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
// dsn 从 config.AppConfig.MySQL.DSN 传入
//
// MySQL 可能正在重启/恢复（部署脚本刚起的库、开机时与 api 同时启动），
// 单次失败就 Fatal 会让 systemd 进入 auto-restart，自检时表现为一直 activating。
// 这里给 60s 重试窗口，期间进程保持存活，连上后再继续建表。
func InitDB(dsn string) {
	if dsn == "" {
		log.Println("[MySQL] DSN 为空，跳过数据库初始化")
		return
	}

	var err error
	for i := 1; i <= 30; i++ {
		db, openErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if openErr == nil {
			DB = db
			log.Println("[MySQL] 数据库连接成功")
			return
		}
		err = openErr
		log.Printf("[MySQL] 连接失败（第 %d/30 次），2s 后重试: %v", i, err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("[MySQL] 数据库链接失败: %v", err)
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
