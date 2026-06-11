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

// TaskImage 任务关联的图片
type TaskImage struct {
	URL  string `json:"url"`
	Desc string `json:"desc"`
}

// Task 生成任务
type Task struct {
	ID          string      `json:"id" gorm:"primaryKey"`
	Target      string      `json:"target"` // flutter | vue3 | react
	Images      []TaskImage `json:"images" gorm:"serializer:json"`
	Status      TaskStatus  `json:"status" gorm:"size:20"`
	Progress    int         `json:"progress"`     // 当前进度 0-100
	CurrentStep string      `json:"current_step"` // 当前执行步骤
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
type TaskStatusResponse struct {
	TaskID      string     `json:"task_id"`
	Target      string     `json:"target"`
	Status      TaskStatus `json:"status"`
	Progress    int        `json:"progress"`
	CurrentStep string     `json:"current_step"`
	CanSSE      bool       `json:"can_sse"`
	Result      *Result    `json:"result,omitempty"`
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
		Result:      result,
	}
}
