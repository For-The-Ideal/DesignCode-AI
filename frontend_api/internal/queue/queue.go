package queue

import (
	"log"
	"sync"
)

// ═══════════════════════════════════════════════
//  任务队列 - 内存实现
//  未来可替换为 Redis Queue (如: github.com/go-redsync/redsync)
// ═══════════════════════════════════════════════

// Queue 任务队列接口
type Queue interface {
	// Enqueue 将任务 ID 加入队列
	Enqueue(taskID string) error
	// Dequeue 阻塞获取一个任务 ID（返回空串表示队列关闭）
	Dequeue() (string, error)
	// Len 返回队列长度
	Len() int
	// Close 关闭队列
	Close()
}

// InMemoryQueue 内存队列（当前实现，后续可替换为 RedisQueue）
type InMemoryQueue struct {
	mu     sync.Mutex
	cond   *sync.Cond
	items  []string
	closed bool
}

// NewInMemoryQueue 创建内存队列
func NewInMemoryQueue() *InMemoryQueue {
	q := &InMemoryQueue{
		items: make([]string, 0, 64),
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Enqueue 入队
func (q *InMemoryQueue) Enqueue(taskID string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil
	}

	q.items = append(q.items, taskID)
	log.Printf("[Queue] enqueued: task=%s (length=%d)", taskID, len(q.items))
	q.cond.Signal()
	return nil
}

// Dequeue 出队（阻塞直到有任务）
func (q *InMemoryQueue) Dequeue() (string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.items) == 0 && !q.closed {
		q.cond.Wait()
	}

	if q.closed {
		return "", nil
	}

	taskID := q.items[0]
	q.items = q.items[1:]
	log.Printf("[Queue] dequeued: task=%s (length=%d)", taskID, len(q.items))
	return taskID, nil
}

// Len 队列长度
func (q *InMemoryQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Close 关闭队列，唤醒所有等待的 Dequeue
func (q *InMemoryQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cond.Broadcast()
}
