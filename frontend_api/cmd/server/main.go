package main

import "frontend_api/internal/app"

// 后端唯一入口：deploy/deploy.sh 构建的也是 ./cmd/server（根目录不再留 main.go）
func main() {
	app.Run()
}
