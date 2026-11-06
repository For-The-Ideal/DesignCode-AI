package sse

import (
	"log"
	"sync"
)

// ═══════════════════════════════════════════════
//  SSE 连接管理器
//  职责：按 taskID / userID 管理事件订阅与推送
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
	Event string `json:"event"` // progress | done | error | task_status
	Data  string `json:"data"`  // JSON 字符串
}

// subscriber 单个订阅连接。
//
// ch 永远不会被 close：取消订阅只是从 map 中摘除，channel 交给 GC 回收。
// 旧实现会在 unsubscribe 里 close(channel)，而 Push 是解锁后才发送，
// 两者并发就会 "send on closed channel" 直接打崩进程（worker goroutine 没有 recover）。
type subscriber struct {
	ch chan SSEEvent
}

// Manager SSE 连接管理器
//
// 并发模型：
//   - 同一 taskID / userID 允许多个订阅者：多标签页、断线重连叠加时互不踢出
//   - Push 持读锁遍历发送（非阻塞 select），unsubscribe 持写锁摘除，二者互斥
//   - channel 永不 close，因此不存在关闭与发送的竞态
type Manager struct {
	mu           sync.RWMutex
	channels     map[string]map[*subscriber]struct{} // key: taskID
	userChannels map[string]map[*subscriber]struct{} // key: userID (fmt.Sprintf("%d", userID))
}

var defaultManager *Manager
var managerOnce sync.Once

// GetManager 获取全局 SSE 管理器（单例）
func GetManager() *Manager {
	managerOnce.Do(func() {
		defaultManager = &Manager{
			channels:     make(map[string]map[*subscriber]struct{}),
			userChannels: make(map[string]map[*subscriber]struct{}),
		}
	})
	return defaultManager
}

// addSubscriber 向分组集合添加订阅者，返回该组当前订阅数（调用方需持写锁）
func addSubscriber(groups map[string]map[*subscriber]struct{}, key string, sub *subscriber) int {
	set := groups[key]
	if set == nil {
		set = make(map[*subscriber]struct{}, 1)
		groups[key] = set
	}
	set[sub] = struct{}{}
	return len(set)
}

// removeSubscriber 从分组集合移除订阅者（幂等），空集合自动清理（调用方需持写锁）
func removeSubscriber(groups map[string]map[*subscriber]struct{}, key string, sub *subscriber) {
	set, ok := groups[key]
	if !ok {
		return
	}
	delete(set, sub)
	if len(set) == 0 {
		delete(groups, key)
	}
}

// Subscribe 注册 taskID 监听，返回事件通道和取消函数
func (m *Manager) Subscribe(taskID string) (<-chan SSEEvent, func()) {
	sub := &subscriber{ch: make(chan SSEEvent, 128)}

	m.mu.Lock()
	n := addSubscriber(m.channels, taskID, sub)
	m.mu.Unlock()
	log.Printf("[SSE Manager] subscribed: task=%s (subscribers: %d)", taskID, n)

	unsubscribe := func() {
		m.mu.Lock()
		removeSubscriber(m.channels, taskID, sub)
		remaining := len(m.channels)
		m.mu.Unlock()
		log.Printf("[SSE Manager] unsubscribed: task=%s (total tasks: %d)", taskID, remaining)
	}

	return sub.ch, unsubscribe
}

// SubscribeUser 注册用户级 SSE 监听，返回事件通道和取消函数
func (m *Manager) SubscribeUser(userID string) (<-chan SSEEvent, func()) {
	sub := &subscriber{ch: make(chan SSEEvent, 64)}

	m.mu.Lock()
	n := addSubscriber(m.userChannels, userID, sub)
	m.mu.Unlock()
	log.Printf("[SSE Manager] user subscribed: user=%s (connections: %d)", userID, n)

	unsubscribe := func() {
		m.mu.Lock()
		removeSubscriber(m.userChannels, userID, sub)
		remaining := len(m.userChannels)
		m.mu.Unlock()
		log.Printf("[SSE Manager] user unsubscribed: user=%s (total users: %d)", userID, remaining)
	}

	return sub.ch, unsubscribe
}

// Push 向指定 taskID 的所有订阅者推送事件；至少送达一个订阅者时返回 true
func (m *Manager) Push(taskID string, event SSEEvent) bool {
	m.mu.RLock()
	set := m.channels[taskID]
	subs := len(set)
	sent, dropped := false, 0
	for sub := range set {
		select {
		case sub.ch <- event:
			sent = true
		default:
			dropped++
		}
	}
	m.mu.RUnlock()

	switch {
	case subs == 0:
		log.Printf("[SSE Manager] task %s has no subscriber, dropping event: %s", taskID, event.Event)
	case dropped > 0:
		log.Printf("[SSE Manager] task %s buffer full, dropped event for %d/%d subscriber(s): %s", taskID, dropped, subs, event.Event)
	default:
		log.Printf("[SSE Manager] task %s pushed event: %s | data=%s", taskID, event.Event, truncateStr(event.Data, 60))
	}
	return sent
}

// PushUser 向指定 userID 的所有连接推送事件（任务状态变更通知）
func (m *Manager) PushUser(userID string, event SSEEvent) bool {
	m.mu.RLock()
	set := m.userChannels[userID]
	subs := len(set)
	sent, dropped := false, 0
	for sub := range set {
		select {
		case sub.ch <- event:
			sent = true
		default:
			dropped++
		}
	}
	m.mu.RUnlock()

	if subs == 0 {
		return false // 用户未连接 SSE，静默丢弃
	}
	if dropped > 0 {
		log.Printf("[SSE Manager] user %s buffer full, dropped event for %d/%d connection(s): %s", userID, dropped, subs, event.Event)
	} else {
		log.Printf("[SSE Manager] user %s pushed event: %s | data=%s", userID, event.Event, truncateStr(event.Data, 60))
	}
	return sent
}

// truncateStr 截断过长字符串用于日志
func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// NumSubscribers 返回当前任务级订阅连接总数
func (m *Manager) NumSubscribers() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	total := 0
	for _, set := range m.channels {
		total += len(set)
	}
	return total
}
