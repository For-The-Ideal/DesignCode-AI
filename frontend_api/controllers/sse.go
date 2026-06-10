package controllers

import (
	"fmt"
	"frontend_api/services"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// SSEController SSE 流式代码生成控制器
type SSEController struct{}

// Handle 建立 SSE 长连接，订阅 Broker 事件并推送给前端
//
// 流程：
//  1. 设置 SSE 响应头
//  2. 向 Broker 注册，获取事件通道
//  3. 阻塞读取 Broker 推送的事件，写入 HTTP 响应
//  4. 客户端断开或超时 → 注销订阅
func (s *SSEController) Handle(c *gin.Context) {
	// 设置 SSE 标准响应头
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

	// 注册到 Broker
	clientID := fmt.Sprintf("sse_%d", time.Now().UnixNano())
	broker := services.GetBroker()
	eventCh, unsubscribe := broker.Subscribe(clientID)
	defer unsubscribe()

	log.Printf("[SSE] client connected: %s", clientID)

	// 发送 connected 事件确认连接
	writeSSEEvent(c, flusher, "connected", "ok")

	// 阻塞读取 Broker 推送的事件
	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				// Broker 通道关闭
				log.Printf("[SSE] client %s: broker channel closed", clientID)
				return
			}
			log.Printf("[SSE] pushing event to client %s: %s", clientID, event.Event)
			if !writeSSEEvent(c, flusher, event.Event, event.Data) {
				return
			}
			// done 后不断开，保持连接等待下一次生成

		case <-c.Request.Context().Done():
			log.Printf("[SSE] client %s: disconnected", clientID)
			return
		}
	}
}

// writeSSEEvent 发送标准 SSE 事件帧: "event: xxx\ndata: yyy\n\n"
func writeSSEEvent(c *gin.Context, flusher interface{ Flush() }, event, data string) bool {
	_, err := c.Writer.Write([]byte(fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)))
	if err != nil {
		return false
	}
	if flusher != nil {
		flusher.Flush()
	}
	return true
}
