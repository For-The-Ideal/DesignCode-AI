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
		if err := mysql.GetDB().AutoMigrate(
			&model.Task{},
			&model.Result{},
			&model.User{},
			&model.ComponentLibrary{},
		); err != nil {
			log.Fatalf("[AutoMigrate] 自动建表失败: %v", err)
		}
		// 初始化组件库种子数据
		var count int64
		mysql.GetDB().Model(&model.ComponentLibrary{}).Count(&count)
		if count == 0 {
			seeds := model.SeedComponentLibraries()
			for i := range seeds {
				if err := mysql.GetDB().Create(&seeds[i]).Error; err != nil {
					log.Printf("[Seed] 插入组件库数据失败: %v", err)
				}
			}
			log.Printf("[Seed] 组件库种子数据已插入 %d 条", len(seeds))
		}
	}

	// 3. 初始化 COS 客户端
	cos.InitClient()

	// 4. 初始化 Logger
	appLog := logger.NewLogger("app")

	// 4. 创建 Repository
	taskRepo := repository.NewTaskRepository()
	resultRepo := repository.NewResultRepository()

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

	// 8. 创建 Worker 并启动
	w := worker.NewGenerateWorker(q, taskRepo, visionSkill, codeToHtmlSkill, resultRepo, sseManager)
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
