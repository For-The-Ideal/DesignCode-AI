package services

import (
	"log"
	"sync"
)

// SSEEvent SSE 推送事件
type SSEEvent struct {
	Event string
	Data  string
}

// SSEBroker SSE 事件广播器，连接 /generate/send 和 /sse
type SSEBroker struct {
	mu          sync.RWMutex
	subscribers map[string]chan SSEEvent
}

var defaultBroker = &SSEBroker{
	subscribers: make(map[string]chan SSEEvent),
}

// GetBroker 获取全局 SSE 事件代理实例
func GetBroker() *SSEBroker {
	return defaultBroker
}

// Subscribe 注册 SSE 客户端，返回事件通道和取消函数
func (b *SSEBroker) Subscribe(clientID string) (chan SSEEvent, func()) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan SSEEvent, 64)
	b.subscribers[clientID] = ch
	log.Printf("[SSE Broker] client subscribed: %s (total: %d)", clientID, len(b.subscribers))

	unsubscribe := func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if ch, ok := b.subscribers[clientID]; ok {
			close(ch)
			delete(b.subscribers, clientID)
			log.Printf("[SSE Broker] client unsubscribed: %s (total: %d)", clientID, len(b.subscribers))
		}
	}

	return ch, unsubscribe
}

// Publish 向所有已注册的 SSE 客户端广播事件
func (b *SSEBroker) Publish(event SSEEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	log.Printf("[SSE Broker] publishing event: %s (subscribers: %d)", event.Event, len(b.subscribers))
	for id, ch := range b.subscribers {
		select {
		case ch <- event:
		default:
			log.Printf("[SSE Broker] client %s buffer full, dropping event", id)
		}
	}
}
