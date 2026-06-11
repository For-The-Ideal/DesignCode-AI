package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ═══════════════════════════════════════════════
//  结构化日志服务
//
//  日志格式: [日期 时间] [级别] [模块] [详情]
//  示例:    [2026-06-08 14:30:25] [INFO] [AI] 代码生成成功
//
//  日志文件: logs/{module}_YYYY-MM-DD.log（按日轮转）
// ═══════════════════════════════════════════════

const (
	logDirName     = "logs"
	logFileNameFmt = "%s_%s.log" // {module}_2026-06-08.log
)

// Level 日志级别
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelDebug Level = "DEBUG"
)

// Logger 结构日志记录器
type Logger struct {
	mu      sync.Mutex
	logDir  string
	logFile *os.File
	logDate string
	module  string
}

var (
	defaultLogger *Logger
	loggerOnce    sync.Once
)

// GetLogger 获取全局日志记录器（单例）
func GetLogger() *Logger {
	loggerOnce.Do(func() {
		defaultLogger = NewLogger("app")
	})
	return defaultLogger
}

// SetModule 设置日志模块名（并发安全）
func (l *Logger) SetModule(module string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.module = module
}

// NewLogger 创建指定模块的日志记录器
func NewLogger(module string) *Logger {
	l := &Logger{module: module}
	l.init()
	return l
}

func (l *Logger) init() {
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	l.logDir = filepath.Join(wd, logDirName)
	if err := os.MkdirAll(l.logDir, 0755); err != nil {
		log.Printf("[Logger] 创建日志目录失败: %v", err)
		return
	}
	l.rotateFile()
}

func (l *Logger) rotateFile() {
	today := time.Now().Format("2006-01-02")
	if l.logDate == today && l.logFile != nil {
		return
	}
	if l.logFile != nil {
		l.logFile.Close()
	}
	l.logDate = today
	filePath := filepath.Join(l.logDir, fmt.Sprintf(logFileNameFmt, l.module, today))
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("[Logger] 打开日志文件失败 %s: %v", filePath, err)
		return
	}
	l.logFile = f
}

// Log 写入日志
func (l *Logger) Log(level Level, detail string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.rotateFile()
	now := time.Now()
	ts := now.Format("2006-01-02 15:04:05")
	line := fmt.Sprintf("[%s] [%s] [%s] %s\n", ts, level, l.module, detail)
	log.Printf("[Logger] %s", line[:len(line)-1])
	if l.logFile != nil {
		l.logFile.WriteString(line)
	}
}

// Info 快捷方法
func (l *Logger) Info(detail string) { l.Log(LevelInfo, detail) }

// Warn 快捷方法
func (l *Logger) Warn(detail string) { l.Log(LevelWarn, detail) }

// Error 快捷方法
func (l *Logger) Error(detail string) { l.Log(LevelError, detail) }

// Debug 快捷方法
func (l *Logger) Debug(detail string) { l.Log(LevelDebug, detail) }

// Infof 格式化 Info
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Log(LevelInfo, fmt.Sprintf(format, args...))
}

// Errorf 格式化 Error
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Log(LevelError, fmt.Sprintf(format, args...))
}

// Close 关闭日志文件
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.logFile != nil {
		l.logFile.Close()
		l.logFile = nil
	}
}
