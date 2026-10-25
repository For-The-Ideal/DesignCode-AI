package model

import "time"

// ═══════════════════════════════════════════════
//  任务领域模型
// ═══════════════════════════════════════════════

// TaskStatus 任务状态
//   pending   - 刚创建，等待队列
//   running   - 执行中（VisionAnalyzeSkill / GeneratorSkill / ReviewSkill）
//   success   - 完成
//   failed    - 失败
type TaskStatus string

const (
	TaskStatusPending TaskStatus = "pending"
	TaskStatusRunning TaskStatus = "running"
	TaskStatusSuccess TaskStatus = "success"
	TaskStatusFailed  TaskStatus = "failed"
)

// ImageItem 上传图片项
type ImageItem struct {
	URL       string   `json:"url"`
	Title     string   `json:"title,omitempty"`  // 每图简短标题（列表/详情回显用，如文件名去扩展名）
	Tree      []UINode `json:"tree,omitempty"`   // 结构编辑器确认的嵌套树（做法 B：有则优先据此生成，跳过 Vision）
	Width     int      `json:"width,omitempty"`  // 原图像素宽（供 DSL canvas 坐标换算）
	Height    int      `json:"height,omitempty"` // 原图像素高
	SortOrder int      `json:"sort_order"`       // 排序，1-N（多图页序，生成时按其升序消费）
}

// TaskStep 任务执行步骤记录
type TaskStep struct {
	Step      string    `json:"step"`      // 步骤名称,  如 VisionAnalyzeSkill
	Progress  int       `json:"progress"`  // 该步骤完成时的进度 0-100
	Status    string    `json:"status"`    // pending | running | done | error
	Timestamp time.Time `json:"timestamp"` // 记录时间
}

// Task 生成任务
type Task struct {
	ID             string      `json:"id" gorm:"primaryKey"`
	Target         string      `json:"target"`                                      // flutter | vue3 | react
	Platform       string      `json:"platform" gorm:"size:20;default:mobile"`      // mobile | desktop | tablet
	Images         []ImageItem `json:"images" gorm:"column:images;serializer:json"` // 上传的图片列表
	Options        []string    `json:"options" gorm:"column:options;serializer:json"`
	Advanced       []string    `json:"advanced" gorm:"column:advanced;serializer:json"`
	ComponentLib   string      `json:"component_lib" gorm:"size:50"`
	Quality        int         `json:"quality"`         // 质量要求 60-100
	RequiredPoints int         `json:"required_points"` // 本次任务消耗的积分数
	Status         TaskStatus  `json:"status" gorm:"size:20"`
	Progress       int         `json:"progress"`                                            // 当前进度 0-100
	CurrentStep    string      `json:"current_step"`                                        // 当前执行步骤
	TaskSteps      []TaskStep  `json:"task_steps" gorm:"column:task_steps;serializer:json"` // 执行步骤记录
	UserID         uint        `json:"user_id" gorm:"index"`                                // 列表/统计按用户过滤，加索引避免全表扫
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// ── 任务状态查询响应 ──────────────────────────

// TaskStatusResponse GET /api/v1/task/:id 的完整响应
//   task_id       - 任务唯一标识
//   status        - pending | running | success | failed
//   target        - flutter | vue3 | react
//   progress      - 执行进度 0-100
//   current_step  - 当前执行步骤，如 VisionAnalyzeSkill / FlutterGenerateSkill
//   can_sse       - 是否可以连接 SSE（true=运行中可连，false=已结束无需连）
//   result        - 生成结果（仅 success 时有值）
//   images        - 上传的图片列表（刷新后回显用）
type TaskStatusResponse struct {
	TaskID      string      `json:"task_id"`
	Target      string      `json:"target"`
	Platform    string      `json:"platform"`
	Status      TaskStatus  `json:"status"`
	Progress    int         `json:"progress"`
	CurrentStep string      `json:"current_step"`
	CanSSE      bool        `json:"can_sse"`
	Images      []ImageItem `json:"images"`
	TaskSteps   []TaskStep  `json:"task_steps"` // 执行步骤记录
	Result      *Result     `json:"result,omitempty"`
}

// ToTaskStatusResponse 将 Task 转换为对外响应
func (t *Task) ToTaskStatusResponse(result *Result) *TaskStatusResponse {
	canSSE := t.Status == TaskStatusPending || t.Status == TaskStatusRunning
	return &TaskStatusResponse{
		TaskID:      t.ID,
		Target:      t.Target,
		Platform:    t.Platform,
		Status:      t.Status,
		Progress:    t.Progress,
		CurrentStep: t.CurrentStep,
		CanSSE:      canSSE,
		Images:      t.Images,
		TaskSteps:   t.TaskSteps,
		Result:      result,
	}
}

// ── 用户任务统计响应 ──────────────────────────

// TaskCountResponse GET /api/v1/task/status 的状态数量统计
type TaskCountResponse struct {
	Pending int `json:"pending"` // 排队中
	Running int `json:"running"` // 执行中
	Success int `json:"success"` // 已完成
	Failed  int `json:"failed"`  // 已失败
}

// ── 任务列表响应 ──────────────────────────

// ListItemImage 列表用轻量图片项：只回缩略图 URL 与标题，不回整棵 tree
// （tree 只在详情 TaskStatusResponse 中返回，供按标注位置回显/重开编辑）
type ListItemImage struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// TaskListItem GET /api/v1/tasks 的列表项（精简字段）
type TaskListItem struct {
	ID           string          `json:"id"`
	Target       string          `json:"framework"` // 前端字段名兼容：framework
	Platform     string          `json:"platform"`
	Status       TaskStatus      `json:"status"`
	Progress     int             `json:"progress"`
	Images       []ListItemImage `json:"images"`
	Options      []string        `json:"options"`
	Advanced     []string        `json:"advanced"`
	ComponentLib string          `json:"component_lib"`
	CreatedAt    string          `json:"created_at"`
}

// ToTaskListItem 将 Task 转为列表项（图片瘦身为 url+title，剥离 tree）
func (t *Task) ToTaskListItem() TaskListItem {
	imgs := make([]ListItemImage, 0, len(t.Images))
	for _, im := range t.Images {
		imgs = append(imgs, ListItemImage{URL: im.URL, Title: im.Title})
	}
	return TaskListItem{
		ID:           t.ID,
		Target:       t.Target,
		Platform:     t.Platform,
		Status:       t.Status,
		Progress:     t.Progress,
		Images:       imgs,
		Options:      t.Options,
		Advanced:     t.Advanced,
		ComponentLib: t.ComponentLib,
		CreatedAt:    t.CreatedAt.Format("2006-01-02 15:04"),
	}
}

// AddStep 追加执行步骤记录
func (t *Task) AddStep(step string, progress int, status string) {
	t.TaskSteps = append(t.TaskSteps, TaskStep{
		Step:      step,
		Progress:  progress,
		Status:    status,
		Timestamp: time.Now(),
	})
}
