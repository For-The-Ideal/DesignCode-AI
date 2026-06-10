package main

import (
	"frontend_api/config"
	"frontend_api/middleware"
	"frontend_api/routes"
	"frontend_api/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	config.InitConfig()

	// // // 2. 初始化数据库
	utils.InitDB()

	// // 3. 设置 Gin 路由
	r := gin.Default()

	// 信任代理（修复 "don't trust all proxies" 警告）
	r.SetTrustedProxies(nil)

	// 全局 CORS 中间件（必须最先执行，处理 OPTIONS 预检）
	r.Use(middleware.CORSMiddleware())

	// // 4. 初始化路由
	routes.InitRoutes(r)

	// // 5. 启动服务器
	log.Printf("Server starting on %s", config.AppConfig.Server.Local+config.AppConfig.Server.Port)
	r.Run(config.AppConfig.Server.Port)
}
