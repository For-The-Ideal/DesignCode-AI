package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// ═══════════════════════════════════════════════
//  应用全局配置
// ═══════════════════════════════════════════════

type Config struct {
	MySQL struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"mysql"`
	Server struct {
		Port      string `mapstructure:"port"`
		Local     string `mapstructure:"local"`
		JWTSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"server"`
	AI    AIConfig    `mapstructure:"ai"`
	Queue QueueConfig `mapstructure:"queue"`
	COS   COSConfig   `mapstructure:"cos"`
}

// AIConfig AI 相关配置
type AIConfig struct {
	Logging AILoggingConfig   `mapstructure:"logging"`
	Models  []AIModelYAMLItem `mapstructure:"models"`
}

// QueueConfig 队列配置
type QueueConfig struct {
	Type string `mapstructure:"type"` // memory | redis
}

// COSConfig 腾讯云对象存储配置
type COSConfig struct {
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	SecretID  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
}

// AILoggingConfig AI 日志配置
type AILoggingConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Directory string `mapstructure:"directory"`
}

// AIModelYAMLItem config.yaml 中 AI 模型的单条配置
type AIModelYAMLItem struct {
	Provider   string `mapstructure:"provider"`
	Name       string `mapstructure:"name"`
	APIKey     string `mapstructure:"api_key"`
	Endpoint   string `mapstructure:"endpoint"`
	Enabled    bool   `mapstructure:"enabled"`
	MaxRetries int    `mapstructure:"max_retries"`
	Timeout    int    `mapstructure:"timeout"` // 秒
	Priority   int    `mapstructure:"priority"`
	Capability string `mapstructure:"capability"` // read(vision分析) / write(代码生成) / both(两者)
}

var AppConfig *Config

// InitConfig 加载配置文件
// 通过 APP_ENV 环境变量区分环境：
//   - APP_ENV=prod → 加载 config.prod.yaml
//   - 其他/未设置 → 加载 config.dev.yaml（本地开发）
func InitConfig() {
	configFile := "config.dev.yaml"
	if env := os.Getenv("APP_ENV"); env == "prod" {
		configFile = "config.prod.yaml"
		log.Printf("[Config] 生产环境，加载 %s", configFile)
	} else {
		log.Printf("[Config] 开发环境，加载 %s", configFile)
	}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	log.Printf("[Config] 配置加载完成，已注册 %d 个 AI 模型", len(AppConfig.AI.Models))
}
