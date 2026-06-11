package app

import (
	"context"
	"frontend_api/config"
	"frontend_api/internal/handler"
	"frontend_api/internal/model"
	"frontend_api/internal/queue"
	"frontend_api/internal/repository"
	"frontend_api/internal/skills/vision"
	"frontend_api/internal/sse"
	"frontend_api/internal/worker"
	"frontend_api/middleware"
	"frontend_api/pkg/logger"
	"frontend_api/pkg/mysql"
	"frontend_api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

// Run 启动服务器（新架构）
func Run() {
	// 1. 加载配置
	config.InitConfig()

	// 2. 初始化数据库
	mysql.InitDB(config.AppConfig.MySQL.DSN)

	// 自动迁移表结构（Task、Result、User）
	if mysql.GetDB() != nil {
		if err := mysql.GetDB().AutoMigrate(&model.Task{}, &model.Result{}, &model.User{}); err != nil {
			log.Fatalf("[AutoMigrate] 自动建表失败: %v", err)
		}
		log.Println("[AutoMigrate] 表结构同步完成")
	}

	// 3. 初始化 Logger
	appLog := logger.NewLogger("app")

	// 4. 创建 Repository
	taskRepo := repository.NewTaskRepository()
	resultRepo := repository.NewResultRepository()

	// 5. 创建内存队列
	q := queue.NewInMemoryQueue()

	// 6. 获取 SSE 管理器（单例）
	sseManager := sse.GetManager()

	// 7. 创建 Vision Skill
	visionSkill := vision.NewSkill()

	// 8. 创建 Worker 并启动（Worker 内部创建 Generator Skill + Workflow）
	w := worker.NewGenerateWorker(q, taskRepo, visionSkill, resultRepo, sseManager)
	go w.Run(context.Background())

	// 9. 创建 Handlers
	taskHandler := handler.NewTaskHandler(taskRepo, resultRepo, q)
	sseHandler := handler.NewSSEHandler(sseManager)
	uploadHandler := handler.NewUploadHandler()
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()

	// 10. 设置 Gin
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(middleware.CORSMiddleware())

	// 11. 初始化路由（新架构 + 保持旧接口兼容）
	routes.InitV1Routes(r, taskHandler, sseHandler, uploadHandler, authHandler, userHandler)

	// 12. 启动服务器
	appLog.Infof("Server starting on %s", config.AppConfig.Server.Local+config.AppConfig.Server.Port)
	log.Printf("Server starting on %s", config.AppConfig.Server.Local+config.AppConfig.Server.Port)
	r.Run(config.AppConfig.Server.Port)
}
