package utils

import (
	"log"

	"frontend_api/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.AppConfig.MySQL.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库链接失败: %v", err)
		return
	}

	// 自动迁移模型
	// DB.AutoMigrate(models.AllModels...)
}
