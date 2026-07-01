package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"frontend_api/internal/model"
	"frontend_api/internal/queue"
	"frontend_api/internal/repository"
	"frontend_api/pkg/logger"
	"frontend_api/pkg/mysql"
	"frontend_api/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ═══════════════════════════════════════════════
//  TaskHandler
//  POST /api/v1/task/create  — 创建生成任务
//  GET  /api/v1/task/:id     — 查询任务状态
// ═══════════════════════════════════════════════

func generateTaskID() string {
	// 时间戳前缀（精确到毫秒，可排序）
	ts := time.Now().Format("20060102150405.000")
	// 4字节随机后缀（8 hex chars）
	b := make([]byte, 4)
	if _, err := rand.Read(b); err == nil {
		return ts + "_" + hex.EncodeToString(b)
	}
	// 极少发生的 fallback
	return ts + fmt.Sprintf("_%04d", time.Now().UnixMilli()%10000)
}

// ═══════════════════════════════════════════════
//  TaskHandler
//  POST /api/v1/task/create  — 创建生成任务
//  GET  /api/v1/task/:id     — 查询任务状态
// ═══════════════════════════════════════════════

// CreateTaskRequest 创建任务请求
//
//	target   - flutter | vue3 | react
//	platform - mobile | desktop | tablet（必选）
//	images   - 图片列表，每项含 url / desc / sort_order，最多5张
//	options  - 生成选项 responsive/comment/component
//	advanced - 高级选项 perf/docs
//	component_lib - 使用的组件库
//
//	示例：{"target":"flutter","platform":"mobile","images":[{"url":"https://...","desc":"首页","sort_order":1}],"options":["responsive"],"advanced":[],"component_lib":"material"}
type CreateTaskRequest struct {
	Target       string            `json:"target" binding:"required"`
	Platform     string            `json:"platform" binding:"required"`
	Images       []model.ImageItem `json:"images" binding:"required,min=1,max=5"`
	Options      []string          `json:"options"`
	Advanced     []string          `json:"advanced"`
	ComponentLib string            `json:"component_lib"`
}

// TaskHandler 任务处理器
type TaskHandler struct {
	taskRepo   *repository.TaskRepository
	resultRepo *repository.ResultRepository
	queue      queue.Queue
	log        *logger.Logger
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(taskRepo *repository.TaskRepository, resultRepo *repository.ResultRepository, q queue.Queue) *TaskHandler {
	return &TaskHandler{
		taskRepo:   taskRepo,
		resultRepo: resultRepo,
		queue:      q,
		log:        logger.NewLogger("task-handler"),
	}
}

// ── POST /api/v1/task/create ──────────────────

// CreateTask 创建生成任务
//
//	请求：{"target":"flutter","images":[{"url":"...","desc":"..."}]}
//	流程：校验参数 → 创建 Task(status=pending) → 入队
//	响应：{"task_id":"xxx","status":"pending"}
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	// 校验 target
	validTargets := map[string]bool{"flutter": true, "vue3": true, "react": true}
	if !validTargets[req.Target] {
		utils.BadRequest(c, "target 必须是 flutter/vue3/react 之一")
		return
	}

	// 校验 platform
	validPlatforms := map[string]bool{"mobile": true, "desktop": true, "tablet": true}
	if !validPlatforms[req.Platform] {
		utils.BadRequest(c, "platform 必须是 mobile/desktop/tablet 之一")
		return
	}

	// 获取用户ID（由 AuthMiddleware 注入）
	userID := getUserID(c)
	if userID == 0 {
		utils.Unauthorized(c, "请先登录")
		return
	}

	// 检查积分并扣减
	db := mysql.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		utils.Error(c, 404, "用户不存在")
		return
	}

	imageCount := len(req.Images)
	if user.Credits < imageCount {
		utils.Error(c, 403, "积分不足，需要 "+fmt.Sprintf("%d", imageCount)+" 积分，当前剩余 "+fmt.Sprintf("%d", user.Credits)+" 积分")
		return
	}

	// 创建任务
	task := &model.Task{
		ID:             generateTaskID(),
		Target:         req.Target,
		Platform:       req.Platform,
		Images:         req.Images,
		Options:        req.Options,
		Advanced:       req.Advanced,
		ComponentLib:   req.ComponentLib,
		Quality:        90, // 后端固定90，后续可以通过逻辑来控制
		RequiredPoints: imageCount,
		Status:         model.TaskStatusPending,
		Progress:       0,
		UserID:         userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 只创建任务（积分在任务成功完成后才扣减）
	if err := db.Create(task).Error; err != nil {
		h.log.Errorf("保存任务失败: %v", err)
		utils.InternalError(c, "创建任务失败")
		return
	}

	// 入队
	if err := h.queue.Enqueue(task.ID); err != nil {
		h.log.Errorf("任务入队失败: %v", err)
		utils.InternalError(c, "任务入队失败")
		return
	}

	h.log.Infof("任务创建成功: id=%s, target=%s, images=%d, user=%d, credits_after=%d", task.ID, task.Target, imageCount, userID, user.Credits-imageCount)

	utils.Success(c, gin.H{
		"task_id": task.ID,
		"status":  task.Status,
	}, "任务创建成功")
}

// ── GET /api/v1/task/:id ──────────────────────

// GetTask 查询任务状态（页面刷新恢复用）
//
//	返回：
//	  task_id       - 任务 ID
//	  status        - pending | running | success | failed
//	  progress      - 当前进度 0-100
//	  current_step  - 当前执行步骤（如 VisionAnalyzeSkill / FlutterGenerateSkill）
//	  can_sse       - 是否可以继续连接 SSE
//	  result        - 生成结果（仅 success 时有值）
//
//	页面刷新流程：
//	  1. 从 localStorage 读取 task_id
//	  2. 请求 GET /api/v1/task/{task_id}
//	  3. 如果 can_sse=true → 重新连接 SSE
//	  4. 前端恢复进度显示
func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := h.taskRepo.GetByID(id)
	if err != nil || task == nil {
		utils.Error(c, 404, "任务不存在")
		return
	}

	// 查询关联结果
	var result *model.Result
	if task.Status == model.TaskStatusSuccess {
		r, err := h.resultRepo.GetByTaskID(task.ID)
		if err == nil {
			result = r
		}
	}

	resp := task.ToTaskStatusResponse(result)
	utils.Success(c, resp, "获取成功")
}

// ── GET /api/v1/task/status ───────────────────────

// GetUserTasks 查询当前用户任务状态数量统计
//
//	返回：pending / running / success / failed 四种状态的数量
func (h *TaskHandler) GetUserTaskStatus(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		utils.Unauthorized(c, "请先登录")
		return
	}

	counts, err := h.taskRepo.CountByUserID(userID)
	if err != nil {
		utils.InternalError(c, "查询任务统计失败")
		return
	}

	utils.Success(c, counts, "获取成功")
}

// ── 获取模版 ──
func (h *TaskHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.BadRequest(c, "模版获取失败: "+err.Error())
		return
	}

	db := mysql.GetDB()
	if db == nil {
		utils.InternalError(c, "数据库未连接: "+err.Error())
		return
	}

	var tmpl model.TaskTemplate
	if err := db.Where("id = ?", id).First(&tmpl).Error; err != nil {
		utils.Error(c, 404, "模版获取失败: "+err.Error())
		return
	}

	utils.Success(c, gin.H{
		"id":            tmpl.ID,
		"name":          tmpl.Name,
		"framework":     tmpl.Framework,
		"template_code": tmpl.TemplateCode,
		"preview_code":  tmpl.PreviewCode,
		"thumbnail":     tmpl.Thumbnail,
		"description":   tmpl.Description,
		"created_at":    tmpl.CreatedAt,
		"updated_at":    tmpl.UpdatedAt,
	}, "获取成功")
}
