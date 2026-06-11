package handler

import (
	"fmt"
	"frontend_api/internal/sse"
	"frontend_api/pkg/logger"
	"log"

	"github.com/gin-gonic/gin"
)

// ═══════════════════════════════════════════════
//  SSEHandler — GET /api/v1/task/:id/events
//  职责：建立 SSE 长连接，按 taskID 订阅事件并推送给前端
// ═══════════════════════════════════════════════

// SSEHandler SSE 事件流处理器
type SSEHandler struct {
	sseManager *sse.Manager
	log        *logger.Logger
}

// NewSSEHandler 创建 SSE 处理器
func NewSSEHandler(manager *sse.Manager) *SSEHandler {
	return &SSEHandler{
		sseManager: manager,
		log:        logger.NewLogger("sse-handler"),
	}
}

// StreamTask 为指定任务建立 SSE 连接
// GET /api/v1/task/:id/events
// 前端: const es = new EventSource(` + "`" + `/api/task/${taskId}/events` + "`" + `)
func (h *SSEHandler) StreamTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(400, gin.H{"error": "缺少 task_id"})
		return
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(200)

	flusher, ok := c.Writer.(interface{ Flush() })
	if !ok {
		fmt.Fprintf(c.Writer, "event: error\ndata: 服务端不支持流式传输\n\n")
		return
	}

	// 注册到 SSE 管理器
	eventCh, unsubscribe := h.sseManager.Subscribe(taskID)
	defer unsubscribe()

	h.log.Infof("[SSE] client connected: task=%s", taskID)

	// 发送 connected 事件
	writeSSEEvent(c, flusher, "connected", "ok")

	// 阻塞读取事件
	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				log.Printf("[SSE] task %s channel closed", taskID)
				return
			}
			if !writeSSEEvent(c, flusher, event.Event, event.Data) {
				return
			}

		case <-c.Request.Context().Done():
			log.Printf("[SSE] task %s: client disconnected", taskID)
			return
		}
	}
}

// writeSSEEvent 发送标准 SSE 帧
func writeSSEEvent(c *gin.Context, flusher interface{ Flush() }, event, data string) bool {
	buf := make([]byte, 0, 512)
	buf = append(buf, "event: "...)
	buf = append(buf, event...)
	buf = append(buf, '\n')

	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			buf = append(buf, "data: "...)
			buf = append(buf, data[start:i]...)
			buf = append(buf, '\n')
			start = i + 1
		}
	}
	if start < len(data) {
		buf = append(buf, "data: "...)
		buf = append(buf, data[start:]...)
		buf = append(buf, '\n')
	}
	buf = append(buf, '\n')

	_, err := c.Writer.Write(buf)
	if err != nil {
		return false
	}
	if flusher != nil {
		flusher.Flush()
	}
	return true
}
