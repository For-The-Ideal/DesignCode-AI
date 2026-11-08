package handler

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"frontend_api/config"
	"frontend_api/internal/model"
	"frontend_api/internal/queue"
	"frontend_api/internal/repository"
	"frontend_api/internal/sse"
	"frontend_api/pkg/logger"
	"frontend_api/pkg/mysql"
	"frontend_api/utils"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════
//  TaskHandler
//  POST /api/v1/task/create  — 创建生成任务
//  GET  /api/v1/task/:id     — 查询任务状态
// ═══════════════════════════════════════════════

func generateTaskID() string {
	// 标准 ULID：26 位 Crockford Base32。
	// 前 48bit = 毫秒时间戳（字典序即时间序，天然可排序），后 80bit = 加密随机。
	var id [16]byte
	ms := uint64(time.Now().UnixMilli())
	id[0] = byte(ms >> 40)
	id[1] = byte(ms >> 32)
	id[2] = byte(ms >> 24)
	id[3] = byte(ms >> 16)
	id[4] = byte(ms >> 8)
	id[5] = byte(ms)
	if _, err := rand.Read(id[6:]); err != nil {
		// 极端情况兜底：用纳秒时间滚动填充随机段（不影响可排序性）
		n := uint64(time.Now().UnixNano())
		for i := 6; i < 16; i++ {
			id[i] = byte(n)
			n = n>>8 | n<<56
		}
	}
	return crockfordBase32(id)
}

// crockfordBase32 将 16 字节（128bit）编码为 26 位 Crockford Base32（大写，不含 I/L/O/U）。
func crockfordBase32(id [16]byte) string {
	const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	n := new(big.Int).SetBytes(id[:])
	mod := big.NewInt(32)
	rem := new(big.Int)
	out := make([]byte, 26)
	for i := 25; i >= 0; i-- {
		n.DivMod(n, mod, rem)
		out[i] = alphabet[rem.Int64()]
	}
	return string(out)
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
//	images   - 页面列表（一图一页），每项含 url / tree / width / height / sort_order
//	           以及多图模式的页面身份：page_id / page_name / page_desc / route / is_entry，最多8页
//	page_links - 页面跳转边（多图模式），如 [{from_page_id:"p1",to_page_id:"p2",trigger:"node_click"}]
//	options  - 生成选项 responsive/comment/component
//	advanced - 高级选项 perf/docs
//	component_lib - 使用的组件库
//
//	示例：{"target":"flutter","platform":"mobile","images":[{"url":"https://...","tree":[...],"width":2732,"height":1534,"sort_order":1,"page_id":"p1","page_name":"商品列表","is_entry":true}],"page_links":[{"from_page_id":"p1","to_page_id":"p2","trigger":"node_click"}]}
type CreateTaskRequest struct {
	Target       string            `json:"target" binding:"required"`
	Platform     string            `json:"platform" binding:"required"`
	Images       []model.ImageItem `json:"images" binding:"required,min=1,max=8"`
	PageLinks    []model.PageLink  `json:"page_links"`
	Options      []string          `json:"options"`
	Advanced     []string          `json:"advanced"`
	ComponentLib string            `json:"component_lib"`
}

// validatePageLinks 校验页面边表：页面引用存在、无自环、触发方式合法、边数受限。
// 空边表合法（单页/多页均可无显式跳转，交由生成模型推断）。
func validatePageLinks(links []model.PageLink, images []model.ImageItem) error {
	if len(links) == 0 {
		return nil
	}
	if len(links) > 64 {
		return fmt.Errorf("页面关系过多（%d 条，上限 64 条）", len(links))
	}
	nameOf := make(map[string]string, len(images))
	for _, img := range images {
		nameOf[img.PageID] = img.PageName
	}
	validTrigger := map[string]bool{"node_click": true, "bottom_tab": true, "back": true, "auto": true}
	for i, l := range links {
		fromName, okFrom := nameOf[l.FromPageID]
		if !okFrom {
			return fmt.Errorf("第 %d 条页面关系：起始页 %q 不存在", i+1, l.FromPageID)
		}
		if _, okTo := nameOf[l.ToPageID]; !okTo {
			return fmt.Errorf("第 %d 条页面关系：目标页 %q 不存在", i+1, l.ToPageID)
		}
		if l.FromPageID == l.ToPageID {
			return fmt.Errorf("第 %d 条页面关系：不能跳转到自己（%s）", i+1, fromName)
		}
		if l.Trigger != "" && !validTrigger[l.Trigger] {
			return fmt.Errorf("第 %d 条页面关系：trigger 必须是 node_click/bottom_tab/back/auto 之一", i+1)
		}
	}
	return nil
}

// 创建任务事务内使用的哨兵错误
var errInsufficientCredits = errors.New("积分不足")

// activeTaskLimitError 单用户活跃任务超限（携带实际数量，便于返回友好提示）
type activeTaskLimitError struct {
	active int64
	limit  int
}

func (e *activeTaskLimitError) Error() string {
	return fmt.Sprintf("活跃任务超限: %d/%d", e.active, e.limit)
}

// TaskHandler 任务处理器
type TaskHandler struct {
	taskRepo   *repository.TaskRepository
	resultRepo *repository.ResultRepository
	queue      queue.Queue
	sseManager *sse.Manager
	log        *logger.Logger
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(taskRepo *repository.TaskRepository, resultRepo *repository.ResultRepository, q queue.Queue, sseManager *sse.Manager) *TaskHandler {
	return &TaskHandler{
		taskRepo:   taskRepo,
		resultRepo: resultRepo,
		queue:      q,
		sseManager: sseManager,
		log:        logger.NewLogger("task-handler"),
	}
}

// ── POST /api/v1/task/create ──────────────────

// CreateTask 创建生成任务
//
//	请求：{"target":"flutter","images":[{"url":"...","tree":[...],"width":2732,"height":1534}]}
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

	// 图片地址白名单校验：只允许 https + COS/配置域名，防止生成任务被用作 SSRF 探测内网
	for i, img := range req.Images {
		if err := utils.ValidateImageURL(img.URL); err != nil {
			utils.BadRequest(c, fmt.Sprintf("第 %d 张图片地址不合法: %v", i+1, err))
			return
		}
	}

	// 检查积分
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
		PageLinks:      req.PageLinks,
		Options:        req.Options,
		Advanced:       req.Advanced,
		ComponentLib:   req.ComponentLib,
		RequiredPoints: imageCount,
		Status:         model.TaskStatusPending,
		Progress:       0,
		UserID:         userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 页面身份兜底（page_id/page_name/入口页/路由）+ 页面跳转边校验。
	// 单图任务同样走一遍，保证新旧请求都能拿到完整的页面元信息。
	task.NormalizePages()
	if err := validatePageLinks(task.PageLinks, task.Images); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// 事务内完成：活跃上限校验 → 建任务 → 原子扣积分（任一步失败整体回滚）
	err := db.Transaction(func(tx *gorm.DB) error {
		// 单用户活跃任务上限（pending + running），防止单用户占满全局串行队列
		if limit := config.AppConfig.Queue.MaxActivePerUser; limit > 0 {
			var active int64
			if err := tx.Model(&model.Task{}).
				Where("user_id = ? AND status IN ?", userID, []model.TaskStatus{model.TaskStatusPending, model.TaskStatusRunning}).
				Count(&active).Error; err != nil {
				return fmt.Errorf("统计活跃任务失败: %w", err)
			}
			if active >= int64(limit) {
				return &activeTaskLimitError{active: active, limit: limit}
			}
		}

		if err := tx.Create(task).Error; err != nil {
			return fmt.Errorf("保存任务失败: %w", err)
		}

		// 预扣积分：WHERE credits >= ? 保证并发下不会扣成负数（防双花）
		if imageCount > 0 {
			res := tx.Model(&model.User{}).
				Where("id = ? AND credits >= ?", userID, imageCount).
				Updates(map[string]interface{}{
					"credits":      gorm.Expr("credits - ?", imageCount),
					"credits_used": gorm.Expr("credits_used + ?", imageCount),
				})
			if res.Error != nil {
				return fmt.Errorf("扣减积分失败: %w", res.Error)
			}
			if res.RowsAffected == 0 {
				return errInsufficientCredits
			}
		}
		return nil
	})
	if err != nil {
		var limitErr *activeTaskLimitError
		switch {
		case errors.Is(err, errInsufficientCredits):
			utils.Error(c, 403, fmt.Sprintf("积分不足，需要 %d 积分，当前剩余 %d 积分", imageCount, user.Credits))
		case errors.As(err, &limitErr):
			utils.Error(c, http.StatusTooManyRequests,
				fmt.Sprintf("你已有 %d 个任务在排队或生成中（上限 %d 个），请等待完成后再提交", limitErr.active, limitErr.limit))
		default:
			h.log.Errorf("创建任务失败: user=%d, err=%v", userID, err)
			utils.InternalError(c, "创建任务失败")
		}
		return
	}

	// 入队（失败补偿：任务置 failed 并退还预扣积分）
	if err := h.queue.Enqueue(task.ID); err != nil {
		h.log.Errorf("任务入队失败: %v", err)
		_ = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&model.Task{}).Where("id = ?", task.ID).
				Updates(map[string]interface{}{
					"status":       model.TaskStatusFailed,
					"current_step": "任务入队失败",
				}).Error; err != nil {
				return err
			}
			if imageCount > 0 {
				return tx.Model(&model.User{}).Where("id = ?", userID).
					Updates(map[string]interface{}{
						"credits":      gorm.Expr("credits + ?", imageCount),
						"credits_used": gorm.Expr("credits_used - ?", imageCount),
					}).Error
			}
			return nil
		})
		utils.InternalError(c, "任务入队失败，请重试")
		return
	}

	h.log.Infof("任务创建成功: id=%s, target=%s, images=%d, user=%d, credits_after=%d", task.ID, task.Target, imageCount, userID, user.Credits-imageCount)

	// 通过 SSE 通知用户（携带权威计数，前端直接更新，无需回拉 /task/status）
	counts, _ := h.taskRepo.CountByUserID(userID)
	if counts == nil {
		counts = &model.TaskCountResponse{}
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"task_id": task.ID,
		"status":  "pending",
		"counts":  counts,
	})
	h.sseManager.PushUser(fmt.Sprintf("%d", userID), sse.SSEEvent{
		Event: "task_status",
		Data:  string(payload),
	})

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
	userID := getUserID(c)
	task, err := h.taskRepo.GetByID(id)
	if err != nil || task == nil {
		utils.Error(c, 404, "任务不存在")
		return
	}

	// 归属校验：只能查看自己的任务（防止按 ID 遍历读取他人设计稿与生成代码）
	if task.UserID != userID {
		utils.Error(c, 403, "无权查看该任务")
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

// ── GET /api/v1/tasks ───────────────────────

// GetTaskList 查询当前用户的任务列表
//
//	返回：任务列表（按创建时间倒序），精简字段
func (h *TaskHandler) GetTaskList(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		utils.Unauthorized(c, "请先登录")
		return
	}

	tasks, err := h.taskRepo.GetByUserID(userID)
	if err != nil {
		utils.InternalError(c, "查询任务列表失败")
		return
	}

	list := make([]model.TaskListItem, 0, len(tasks))
	for i := range tasks {
		list = append(list, tasks[i].ToTaskListItem())
	}

	utils.Success(c, list, "获取成功")
}

// ── DELETE /api/v1/task/:id ─────────────────

// DeleteTask 删除指定任务
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	userID := getUserID(c)
	if userID == 0 {
		utils.Unauthorized(c, "请先登录")
		return
	}

	task, err := h.taskRepo.GetByID(id)
	if err != nil || task == nil {
		utils.Error(c, 404, "任务不存在")
		return
	}

	if task.UserID != userID {
		utils.Error(c, 403, "无权删除该任务")
		return
	}

	// 执行中的任务不可删除：否则 worker 后续结果回写/进度更新会失靶
	if task.Status == model.TaskStatusPending || task.Status == model.TaskStatusRunning {
		utils.BadRequest(c, "任务排队/执行中，暂时无法删除")
		return
	}

	if err := h.taskRepo.DeleteWithResults(id); err != nil {
		utils.InternalError(c, "删除任务失败")
		return
	}

	h.log.Infof("任务已删除: id=%s, user=%d", id, userID)
	utils.Success(c, nil, "任务已删除")
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
		utils.InternalError(c, "数据库未连接")
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
