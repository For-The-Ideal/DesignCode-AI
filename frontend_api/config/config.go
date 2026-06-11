package config

import (
	"frontend_api/models"
	"log"
	"time"

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
	AI AIConfig `mapstructure:"ai"`
}

// AIConfig AI 相关配置
type AIConfig struct {
	Logging AILoggingConfig   `mapstructure:"logging"`
	Models  []AIModelYAMLItem `mapstructure:"models"`
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
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	log.Printf("[Config] 配置加载完成，已注册 %d 个 AI 模型", len(AppConfig.AI.Models))
}

// LoadModelsIntoManager 将 config.yaml 中的 AI 模型配置加载到 AIModelManager
func LoadModelsIntoManager(mgr *models.AIModelManager) {
	if AppConfig == nil {
		log.Println("[Config] AppConfig 未初始化，跳过模型加载")
		return
	}

	loaded := 0
	for _, item := range AppConfig.AI.Models {
		cfg := &models.AIModelConfig{
			Provider:   models.ModelProvider(item.Provider),
			Name:       item.Name,
			APIKey:     item.APIKey,
			Endpoint:   item.Endpoint,
			Enabled:    item.Enabled,
			MaxRetries: item.MaxRetries,
			Timeout:    time.Duration(item.Timeout) * time.Second,
			Priority:   item.Priority,
		}
		mgr.Register(cfg.Provider, cfg)
		loaded++
	}

	log.Printf("[Config] 已加载 %d 个 AI 模型配置到管理器", loaded)
}
