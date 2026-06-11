package sse

import (
	"log"
	"sync"
)

// ═══════════════════════════════════════════════
//  SSE 连接管理器
//  职责：按 taskID 管理事件订阅/推送
//
//  推送协议：
//  进度：
//    event: progress
//    data: {"progress":20,"step":"VisionAnalyzeSkill"}
//
//  完成：
//    event: done
//    data: {"progress":100,"step":"Done"}
//
//  失败：
//    event: error
//    data: {"message":"AI接口超时"}
// ═══════════════════════════════════════════════

// SSEEvent SSE 推送事件
type SSEEvent struct {
	Event string `json:"event"` // progress | done | error
	Data  string `json:"data"`  // JSON 字符串
}

// Manager SSE 连接管理器
type Manager struct {
	mu       sync.RWMutex
	channels map[string]chan SSEEvent // key: taskID
}

var defaultManager *Manager
var managerOnce sync.Once

// GetManager 获取全局 SSE 管理器（单例）
func GetManager() *Manager {
	managerOnce.Do(func() {
		defaultManager = &Manager{
			channels: make(map[string]chan SSEEvent),
		}
	})
	return defaultManager
}

// Subscribe 注册 taskID 监听，返回事件通道和取消函数
func (m *Manager) Subscribe(taskID string) (<-chan SSEEvent, func()) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch := make(chan SSEEvent, 128)
	m.channels[taskID] = ch
	log.Printf("[SSE Manager] subscribed: task=%s (total: %d)", taskID, len(m.channels))

	unsubscribe := func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if c, ok := m.channels[taskID]; ok {
			close(c)
			delete(m.channels, taskID)
			log.Printf("[SSE Manager] unsubscribed: task=%s (total: %d)", taskID, len(m.channels))
		}
	}

	return ch, unsubscribe
}

// Push 向指定 taskID 推送事件
func (m *Manager) Push(taskID string, event SSEEvent) bool {
	m.mu.RLock()
	ch, ok := m.channels[taskID]
	m.mu.RUnlock()

	if !ok {
		log.Printf("[SSE Manager] task %s has no subscriber, dropping event: %s", taskID, event.Event)
		return false
	}

	select {
	case ch <- event:
		return true
	default:
		log.Printf("[SSE Manager] task %s channel full, dropping event: %s", taskID, event.Event)
		return false
	}
}

// NumSubscribers 返回当前订阅数
func (m *Manager) NumSubscribers() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.channels)
}
