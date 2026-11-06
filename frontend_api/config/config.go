package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

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
		SiteURL   string `mapstructure:"site_url"`
		JWTSecret string `mapstructure:"jwt_secret"`

		// AdminEmails 管理员邮箱白名单：命中才能调用 /api/v1/admin/*；为空时全部拒绝
		AdminEmails []string `mapstructure:"admin_emails"`

		// AllowedImageHosts 生成任务允许下载的图片域名（除 COS 桶域名外的额外放行项）
		// 支持精确域名（cdn.example.com）与通配（*.example.com）
		AllowedImageHosts []string `mapstructure:"allowed_image_hosts"`
	} `mapstructure:"server"`
	AI    AIConfig    `mapstructure:"ai"`
	Queue QueueConfig `mapstructure:"queue"`
	COS   COSConfig   `mapstructure:"cos"`
	Email EmailConfig `mapstructure:"email"`
}

// AIConfig AI 相关配置
type AIConfig struct {
	Logging AILoggingConfig   `mapstructure:"logging"`
	Models  []AIModelYAMLItem `mapstructure:"models"`
}

// QueueConfig 队列配置
type QueueConfig struct {
	Type string `mapstructure:"type"` // memory | redis

	// MaxActivePerUser 单用户允许同时存在的活跃任务数（pending + running）
	// 达到上限后拒绝提交新任务；<=0 时使用默认值 2
	MaxActivePerUser int `mapstructure:"max_active_per_user"`

	// WorkerCount 并发消费任务的 worker 数量（>=1）
	// 1 表示串行执行；调大后可提升吞吐，但需关注内存与 AI 限流；<=0 时使用默认值 1
	WorkerCount int `mapstructure:"worker_count"`
}

// COSConfig 腾讯云对象存储配置
type COSConfig struct {
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	SecretID  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
}

// EmailConfig 邮件服务配置
type EmailConfig struct {
	Providers map[string]SMTPConfig `mapstructure:"providers"`
}

// SMTPConfig 单个 SMTP 服务商配置
type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
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
//
// 配置文件不依赖启动目录：当前目录找不到时会向上逐级查找，
// 因此从 cmd/server 子目录（air / go run）启动也能正常工作。
func InitConfig() {
	envName, envLabel := "dev", "开发"
	if os.Getenv("APP_ENV") == "prod" {
		envName, envLabel = "prod", "生产"
	}
	configFile := "config." + envName + ".yaml"

	configPath := resolveConfigPath(configFile)
	if configPath == "" {
		log.Fatalf("[Config] 找不到 %s（可在 frontend_api 根目录启动，或用 APP_CONFIG_PATH 指定文件路径）", configFile)
	}
	log.Printf("[Config] %s环境，加载配置文件: %s", envLabel, configPath)

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// 启动安全校验：JWT 密钥为空等于任何人都能伪造登录 token
	if AppConfig.Server.JWTSecret == "" {
		log.Fatalf("[Config] 缺少 server.jwt_secret，拒绝启动（否则任何人都能伪造登录 token）")
	}
	if len(AppConfig.Server.JWTSecret) < 16 {
		log.Printf("[Config] [WARN] server.jwt_secret 过短，建议使用 32 位以上随机字符串")
	}

	// 队列默认值：单用户活跃任务上限（未配置时回退 2）、worker 数（未配置时回退 1）
	if AppConfig.Queue.MaxActivePerUser <= 0 {
		AppConfig.Queue.MaxActivePerUser = 2
	}
	if AppConfig.Queue.WorkerCount <= 0 {
		AppConfig.Queue.WorkerCount = 1
	}

	log.Printf("[Config] 配置加载完成，已注册 %d 个 AI 模型（单用户活跃任务上限=%d，worker 数=%d）",
		len(AppConfig.AI.Models), AppConfig.Queue.MaxActivePerUser, AppConfig.Queue.WorkerCount)
}

// resolveConfigPath 查找配置文件，返回可用路径（未找到返回空串）。
//
// 查找顺序：
//  1. APP_CONFIG_PATH 环境变量（支持绝对路径或相对当前目录）
//  2. 当前工作目录
//  3. 向上逐级（最多 3 层，覆盖 cmd/server 子目录启动的场景）
func resolveConfigPath(name string) string {
	if p := strings.TrimSpace(os.Getenv("APP_CONFIG_PATH")); p != "" {
		if abs, err := filepath.Abs(p); err == nil {
			if _, statErr := os.Stat(abs); statErr == nil {
				return abs
			}
		}
		log.Printf("[Config] [WARN] APP_CONFIG_PATH 指向的文件不存在: %s", p)
	}

	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for depth := 0; depth < 4; depth++ {
		candidate := filepath.Join(dir, name)
		if _, statErr := os.Stat(candidate); statErr == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // 已到盘符根，避免死循环
		}
		dir = parent
	}
	return ""
}
