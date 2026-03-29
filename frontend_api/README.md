# Design2Code AI - Frontend API

## 🚀 技术栈

- **语言**: [Go (Golang)](https://golang.org/)
- **Web 框架**: [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP 路由框架。
- **ORM**: [GORM](https://gorm.io/) - 强大的 Go 对象关系映射库。
- **配置管理**: [Viper](https://github.com/spf13/viper) - 灵活的配置解决方案。
- **数据库**: MySQL - 存储用户数据、设计稿及生成代码。

## 📁 项目结构

frontend_api/
├── config/             # 配置解析逻辑
│   └── config.go       # 映射 config.yaml 到 Go 结构体
├── controllers/        # 控制器层 (处理业务逻辑)
├── middleware/         # 中间件 (JWT, 日志, CORS 等)
├── models/             # 数据模型 (数据库表结构定义)
├── routes/             # 路由定义 (接口路径映射)
├── skills/             # skills层 (处理 AI 生成逻辑)
├── utils/              # 辅助工具类
│   ├── db.go           # 数据库连接与初始化
│   └── response.go     # 统一 API 响应格式
├── config.yaml         # 运行环境配置文件 (DSN, 端口等)
├── main.go             # 项目入口
├── go.mod              # 依赖管理文件
├── go.sum              # 依赖锁定文件
├── init.md             # 初始化说明
└── README.md           # 项目文档

## 🛠️ 安装与运行

### 1. 克隆并进入目录

  cd frontend_api

### 2. 初始化项目

  项目初始化 ：在项目根目录下执行了 go mod init

### 3. 安装依赖

go mod tidy

### 4. 配置数据库

修改 `config.yaml` 文件中的 `mysql.dsn` 字段，确保数据库连接字符串正确：

mysql:
  dsn: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

### 5. 运行服务

go run main.go


