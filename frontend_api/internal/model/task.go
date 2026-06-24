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
	URL       string `json:"url"`
	Desc      string `json:"desc"`
	SortOrder int    `json:"sort_order"` // 排序，1-N
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
	ID          string      `json:"id" gorm:"primaryKey"`
	Target      string      `json:"target"`                                      // flutter | vue3 | react
	Platform    string      `json:"platform" gorm:"size:20;default:mobile"`      // mobile | desktop | tablet
	Images      []ImageItem `json:"images" gorm:"column:images;serializer:json"` // 上传的图片列表
	Quality     int         `json:"quality"`                                     // 质量要求 60-100
	Status      TaskStatus  `json:"status" gorm:"size:20"`
	Progress    int         `json:"progress"`                                            // 当前进度 0-100
	CurrentStep string      `json:"current_step"`                                        // 当前执行步骤
	TaskSteps   []TaskStep  `json:"task_steps" gorm:"column:task_steps;serializer:json"` // 执行步骤记录
	UserID      uint        `json:"user_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
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
		Status:      t.Status,
		Progress:    t.Progress,
		CurrentStep: t.CurrentStep,
		CanSSE:      canSSE,
		Images:      t.Images,
		TaskSteps:   t.TaskSteps,
		Result:      result,
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
