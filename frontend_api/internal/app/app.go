package app

import (
	"context"
	"frontend_api/config"
	"frontend_api/internal/handler"
	"frontend_api/internal/model"
	"frontend_api/internal/queue"
	"frontend_api/internal/repository"
	"frontend_api/internal/skills/code_to_html"
	"frontend_api/internal/skills/vision"
	"frontend_api/internal/sse"
	"frontend_api/internal/worker"
	"frontend_api/middleware"
	"frontend_api/pkg/ai"
	"frontend_api/pkg/cos"
	"frontend_api/pkg/logger"
	"frontend_api/pkg/mysql"
	"frontend_api/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

const shutdownTimeout = 30 * time.Second

// Run 启动服务器（新架构）
func Run() {
	// 1. 加载配置
	config.InitConfig()

	// 2. 初始化数据库
	mysql.InitDB(config.AppConfig.MySQL.DSN)

	// 自动迁移表结构（Task、Result、User、MembershipPlan、PaymentOrder 等）
	if mysql.GetDB() != nil {
		if err := mysql.GetDB().AutoMigrate(
			&model.Task{},
			&model.Result{},
			&model.User{},
			&model.MembershipPlan{},
			&model.CreditsPackage{},
			&model.PaymentOrder{},
		); err != nil {
			log.Fatalf("[AutoMigrate] 自动建表失败: %v", err)
		}
	}

	// 3. 初始化 COS 客户端
	cos.InitClient()

	// 4. 初始化 Logger
	appLog := logger.NewLogger("app")

	// 4. 创建 Repository
	taskRepo := repository.NewTaskRepository()
	resultRepo := repository.NewResultRepository()

	// 4.1 恢复残留的 running 状态任务（服务重启后清理）
	if n, err := taskRepo.ResetRunningTasks(); err != nil {
		appLog.Errorf("[App] 恢复残留任务失败: %v", err)
	} else if n > 0 {
		appLog.Infof("[App] 已恢复 %d 个残留 running 任务为 failed（服务重启）", n)
	}

	// 5. 创建内存队列
	q := queue.NewInMemoryQueue()

	// 6. 获取 SSE 管理器（单例）
	sseManager := sse.GetManager()

	// 7. 创建 Skills
	readAI, err := ai.SelectClient("read")
	if err != nil {
		appLog.Warn("[App] 未配置视觉AI模型(千问)，视觉分析将不可用")
	}
	visionSkill := vision.NewSkill(readAI)
	writeAI, err := ai.SelectClient("write")
	if err != nil {
		log.Fatalf("初始化AI写客户端失败: %v", err)
	}
	codeToHtmlSkill := code_to_html.NewSkill(writeAI)

	// 8. 创建 Worker 并启动（带可取消 context）
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w := worker.NewGenerateWorker(q, taskRepo, visionSkill, codeToHtmlSkill, resultRepo, sseManager)
	go w.Run(ctx)

	// 9. 创建 Handlers
	taskHandler := handler.NewTaskHandler(taskRepo, resultRepo, q, sseManager)
	sseHandler := handler.NewSSEHandler(sseManager)
	uploadHandler := handler.NewUploadHandler()
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	adminHandler := handler.NewAdminHandler()
	membershipHandler := handler.NewMembershipHandler()

	// 10. 设置 Gin
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(middleware.CORSMiddleware())

	// 11. 初始化路由
	routes.InitV1Routes(r, taskHandler, sseHandler, uploadHandler, authHandler, userHandler, adminHandler, membershipHandler)

	// 12. 创建 HTTP Server（支持优雅关闭）
	// Addr 是监听地址（host:port），不是对外 URL
	// 生产环境由 nginx 反代，后端监听所有网卡即可
	addr := config.AppConfig.Server.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 13. 信号监听（优雅关闭）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动 HTTP 服务（非阻塞）
	go func() {
		appLog.Infof("Server starting on %s", addr)
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server 启动失败: %v", err)
		}
	}()

	// 等待退出信号
	sig := <-quit
	appLog.Infof("收到信号 %v，开始优雅关闭...", sig)
	log.Printf("收到信号 %v，开始优雅关闭...", sig)

	// Step 1: 停止 Worker（不取新任务）
	cancel()
	q.Close() // 唤醒可能阻塞在 Dequeue 的 Worker

	// Step 2: 关闭 HTTP Server（不再接受新请求，等待已有请求完成）
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLog.Errorf("Server 强制关闭: %v", err)
	} else {
		appLog.Info("Server 已安全关闭")
	}

	log.Println("服务已退出")
}
