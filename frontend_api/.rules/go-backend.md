---
trigger: always_on
alwaysApply: true
scope: workspace
---

# DesignCode-AI Go 后端编码规范

> **执行规则**：每次生成代码前，必须逐章自查本规范。

## 技术栈

- **框架**: Gin (github.com/gin-gonic/gin)
- **配置**: Viper (github.com/spf13/viper) 读取 config.yaml
- **数据库**: MySQL + GORM (AutoMigrate)
- **文件存储**: 腾讯云 COS
- **认证**: JWT
- **端口**: :8888

## 一、目录结构

```
frontend_api/
├── main.go                    # 入口
├── config.yaml                # 配置文件
├── cmd/server/main.go         # 启动命令
├── config/config.go           # 配置结构体 + viper 加载
├── middleware/                 # 中间件（auth.go / cors.go）
├── routes/                    # 路由注册
│   ├── router.go              # 总入口 InitV1Routes()
│   ├── auth_route.go          # /api/v1/auth/*
│   ├── user_route.go          # /api/v1/user/*
│   └── generate_route.go      # /api/v1/task/* /template/* /upload
├── internal/
│   ├── app/app.go             # Run() 启动流程
│   ├── handler/               # HTTP Handler
│   │   ├── task_handler.go
│   │   ├── sse_handler.go
│   │   ├── upload_handler.go
│   │   ├── auth_handler.go
│   │   └── user_handler.go
│   ├── model/                 # GORM 模型 + DTO
│   │   ├── task.go
│   │   ├── result.go
│   │   ├── template.go
│   │   ├── user.go
│   │   └── dsl.go
│   ├── repository/            # 数据访问层 (DB 封装)
│   ├── queue/                 # 任务队列 (memory/redis)
│   ├── sse/                   # SSE 连接管理器
│   ├── worker/                # 后台消费队列
│   ├── workflow/              # 生成工作流编排
│   └── skills/                # AI 技能
│       ├── generator/         # 代码生成 (Vue3/React/Flutter)
│       ├── vision/analyze_skill.go  # 视觉分析
│       └── code_to_html/      # 代码→HTML 转换
├── pkg/                       # 基础设施
│   ├── ai/client.go           # AI 客户端（多模型分发）
│   ├── cos/cos.go             # COS 上传
│   ├── mysql/mysql.go         # DB 连接
│   └── logger/logger.go       # 日志
└── utils/                     # 工具函数
    ├── response.go            # 统一响应封装
    └── captcha.go             # 验证码
```

## 二、分层架构

```
routes  → handler → repository → MySQL
            ↓
        workflow/skills → AI client
            ↓
          queue/worker → SSE
```

1. **routes**: 只做路由注册，不写业务逻辑
2. **handler**: 参数校验 + 调用下层 + 返回响应
3. **repository**: 纯 DB 操作，封装 GORM 查询
4. **workflow/skills**: 业务编排和 AI 调用
5. **pkg**: 基础设施，无业务依赖

## 三、Handler 规范

1. 所有 Handler 使用依赖注入（通过构造函数传入 repository/client）
2. 参数用 `c.Param("id")`（路径）或 `c.Query("key")`（查询）或 `c.ShouldBindJSON(&req)`（body）
3. 响应统一使用 `utils.Success/Error/BadRequest/InternalError/Unauthorized`

```go
// ✅ 正确
func (h *TaskHandler) GetTask(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil { utils.BadRequest(c, "id 无效"); return }
    task, err := h.repo.FindByID(uint(id))
    if err != nil { utils.InternalError(c, err.Error()); return }
    utils.Success(c, task, "ok")
}
```

## 四、路由规范

1. 所有接口以 `/api/v1/` 开头
2. 路由按功能模块拆分文件：`auth_route.go` / `user_route.go` / `generate_route.go`
3. 使用 `r.Group("/api/v1")` 创建路由组
4. 路径参数用 `:id` 格式

## 五、Model 规范

1. GORM 模型放在 `internal/model/` 下
2. 使用 GORM 标签定义字段映射
3. 表结构通过 `AutoMigrate` 自动同步
4. JSON 标签使用小驼峰命名

```go
type Task struct {
    ID        uint       `gorm:"primaryKey" json:"id"`
    Status    TaskStatus `gorm:"type:varchar(20);default:pending" json:"status"`
    CreatedAt time.Time  `json:"created_at"`
}
```

## 六、AI 技能规范

1. 所有技能实现 `Skill` 接口（`internal/skills/skill.go`）
2. 技能按能力分类：`vision`（视觉分析）/ `generator`（代码生成）/ `code_to_html`（后处理）
3. 每个技能一个文件，放在对应子目录下
4. AI 调用统一走 `pkg/ai/client.go` 的 `AIClient`

## 七、配置规范

1. 配置定义在 `config.yaml`，结构体在 `config/config.go`
2. 使用 viper 读取，禁止硬编码
3. `server.port` 为 `:8888`
4. `server.local` 为 `http://localhost`（生产覆盖为 IP）
5. `ai.models` 数组配置多模型，`enabled: true` 的才启用

## 八、错误处理

1. 所有 error 必须处理，禁止 `_` 忽略
2. Handler 返回前必须调用一个 `utils.xxx()` 响应函数
3. Repository 返回原始 error，由 Handler 决定响应码

## 九、命名规范

1. **文件名**: 小写下划线，如 `task_handler.go`、`auth_route.go`
2. **包名**: 小写单数，如 `handler`、`model`、`repository`
3. **结构体**: PascalCase，如 `TaskHandler`、`AIClient`
4. **函数**: camelCase（小写开头私有，大写开头公开）
5. **常量**: PascalCase 或 UPPER_SNAKE_CASE，如 `TaskStatusPending`
