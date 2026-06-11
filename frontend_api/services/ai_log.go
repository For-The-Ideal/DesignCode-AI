package services

import (
	"fmt"
	"frontend_api/models"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ═══════════════════════════════════════════════
//  AI 调用日志服务
//
//  日志格式: [日期 时间] [模型名称] [请求状态] [耗时] [详情]
//  示例:    [2026-06-08 14:30:25] [DeepSeek V3] [success] [3.2s] 代码生成成功
//
//  日志文件: logs/ai_call_YYYY-MM-DD.log（按日轮转）
// ═══════════════════════════════════════════════

const (
	logDirName     = "logs"
	logFileNameFmt = "ai_call_%s.log" // ai_call_2026-06-08.log
)

// AILogger AI 模型调用日志记录器
type AILogger struct {
	mu      sync.Mutex
	logDir  string
	logFile *os.File
	logDate string // 记录当前文件日期，用于按日轮转
}

var defaultLogger *AILogger
var loggerOnce sync.Once

// GetAILogger 获取全局 AI 日志记录器实例（单例）
func GetAILogger() *AILogger {
	loggerOnce.Do(func() {
		defaultLogger = &AILogger{}
		defaultLogger.init()
	})
	return defaultLogger
}

// init 初始化日志目录并打开当日文件
func (l *AILogger) init() {
	// 日志目录取项目根目录下的 logs/（可由 config 扩展）
	wd, err := os.Getwd()
	if err != nil {
		// 回退到当前目录
		wd = "."
	}
	l.logDir = filepath.Join(wd, logDirName)

	// 创建目录（如已存在则忽略错误）
	if err := os.MkdirAll(l.logDir, 0755); err != nil {
		log.Printf("[AILogger] 创建日志目录失败: %v", err)
		return
	}

	l.rotateFile()
}

// rotateFile 按日轮转：如果日期变化则创建新文件
func (l *AILogger) rotateFile() {
	today := time.Now().Format("2006-01-02")
	if l.logDate == today && l.logFile != nil {
		return // 同一天，复用已有文件句柄
	}

	// 关闭旧文件
	if l.logFile != nil {
		l.logFile.Close()
	}

	l.logDate = today
	filePath := filepath.Join(l.logDir, fmt.Sprintf(logFileNameFmt, today))

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("[AILogger] 打开日志文件失败 %s: %v", filePath, err)
		return
	}

	l.logFile = f
	log.Printf("[AILogger] 日志文件已创建: %s", filePath)
}

// Log 记录一次 AI 模型调用
//
// 参数:
//   - provider: 模型供应商（如 deepseek, chatgpt）
//   - status:   请求状态（success / fail / retry / failover / max_retries）
//   - duration: 耗时（可传 0）
//   - detail:   补充描述（如 "代码生成成功" / "API 密钥无效"）
func (l *AILogger) Log(provider models.ModelProvider, status models.RequestStatus, duration time.Duration, detail string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 按日轮转
	l.rotateFile()

	// 构建日志行
	now := time.Now()
	ts := now.Format("2006-01-02 15:04:05")
	durStr := ""
	if duration > 0 {
		durStr = fmt.Sprintf("[%s]", duration.Truncate(time.Millisecond*100).String())
	}

	line := fmt.Sprintf("[%s] [%s] [%s] %s %s\n",
		ts, provider, status, durStr, detail)

	// 同时写入标准 log（开发调试用）
	log.Printf("[AILog] %s", line[:len(line)-1])

	// 写入日志文件
	if l.logFile != nil {
		if _, err := l.logFile.WriteString(line); err != nil {
			log.Printf("[AILogger] 写入日志失败: %v", err)
		}
	}
}

// LogResult 通过 AIRequestResult 结构体记录日志（便捷方法）
func (l *AILogger) LogResult(result *models.AIRequestResult, detail string) {
	if detail == "" {
		switch result.Status {
		case models.StatusSuccess:
			detail = "代码生成成功"
		case models.StatusFail:
			detail = result.ErrorMessage
			if detail == "" {
				detail = "AI 请求失败"
			}
		case models.StatusRetry:
			detail = fmt.Sprintf("第 %d 次重试", result.Duration.Milliseconds()/1000+1)
		case models.StatusFailover:
			detail = fmt.Sprintf("已切换至备用模型: %s", result.Provider)
		case models.StatusMaxRetries:
			detail = "已达最大重试次数"
		default:
			detail = string(result.Status)
		}
	}
	l.Log(result.Provider, result.Status, result.Duration, detail)
}

// Close 关闭日志文件（通常在应用退出时调用）
func (l *AILogger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logFile != nil {
		l.logFile.Close()
		l.logFile = nil
	}
}
