package model

import (
	"fmt"
	"sort"
	"time"
)

// ═══════════════════════════════════════════════
//  任务领域模型
// ═══════════════════════════════════════════════

// TaskStatus 任务状态
//
//	pending   - 刚创建，等待队列
//	running   - 执行中（VisionAnalyzeSkill / GeneratorSkill / ReviewSkill）
//	success   - 完成
//	failed    - 失败
type TaskStatus string

const (
	TaskStatusPending TaskStatus = "pending"
	TaskStatusRunning TaskStatus = "running"
	TaskStatusSuccess TaskStatus = "success"
	TaskStatusFailed  TaskStatus = "failed"
)

// AttachedAsset 主图归属的附图素材（真实素材：Logo/图标/插画/照片）
// 树节点通过 assetId 引用，生成时直接使用对应图片而不再由 AI 重绘。
type AttachedAsset struct {
	AssetID   string `json:"asset_id"`
	URL       string `json:"url"`
	Label     string `json:"label"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	SortOrder int    `json:"sort_order"`
}

// ImageItem 上传图片项
// 多图模式下每张图即一个页面：PageID/PageName/PageDesc/Route/IsEntry 描述页面身份，
// 页面之间的跳转关系统一存在 Task.PageLinks（边表），不在此处重复表达。
type ImageItem struct {
	URL            string          `json:"url"`
	Title          string          `json:"title,omitempty"`           // 每图简短标题（列表/详情回显用，如文件名去扩展名）
	Tree           []UINode        `json:"tree,omitempty"`            // 结构编辑器确认的嵌套树（做法 B：有则优先据此生成，跳过 Vision）
	AttachedAssets []AttachedAsset `json:"attached_assets,omitempty"` // 该主图名下的附图素材（节点 assetId 引用）
	Width          int             `json:"width,omitempty"`           // 原图像素宽（供 DSL canvas 坐标换算）
	Height         int             `json:"height,omitempty"`          // 原图像素高
	SortOrder      int             `json:"sort_order"`                // 排序，1-N（多图页序，生成时按其升序消费）

	// ── 页面身份（多图模式；单图缺省时由后端按 sort_order 兜底填充）──
	PageID   string `json:"page_id,omitempty"`   // 页面唯一标识（p1/p2...），PageLink 引用它
	PageName string `json:"page_name,omitempty"` // 页面名（如「商品列表」），生成时作为页面标题
	PageDesc string `json:"page_desc,omitempty"` // 页面说明（用户填写的页面级 brief，喂给生成模型）
	Route    string `json:"route,omitempty"`     // 路由（如 /list、/detail/:id），缺省由生成模型推断
	IsEntry  bool   `json:"is_entry,omitempty"`  // 是否入口页（App 启动页）；多页时至少一个
}

// PageLink 页面跳转边（多图模式的页面关系图）
//
//	trigger = node_click  某节点点击进入目标页（NodeKey 为空表示由生成模型推断触发元素）
//	trigger = bottom_tab  底部导航切换（对应 tabbar 项）
//	trigger = back        返回上一页（生成 AppBar 返回按钮）
//	trigger = auto        仅声明层级，触发方式交给生成模型
//
// 页面关系是图（多入边/多出边），因此用边表而非每页一个 relation 字段。
type PageLink struct {
	FromPageID string `json:"from_page_id"`
	ToPageID   string `json:"to_page_id"`
	Trigger    string `json:"trigger"`            // node_click | bottom_tab | back | auto
	NodeKey    string `json:"node_key,omitempty"` // 挂到哪个标注节点（空=整页/自动推断）
	Label      string `json:"label,omitempty"`    // 触发元素描述，如「商品卡片」
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
	PageLinks      []PageLink  `json:"page_links" gorm:"column:page_links;serializer:json"` // 页面跳转边（多图模式）
	RequiredPoints int         `json:"required_points"`                                     // 本次任务消耗的积分数
	Status         TaskStatus  `json:"status" gorm:"size:20"`
	Progress       int         `json:"progress"`                                            // 当前进度 0-100
	CurrentStep    string      `json:"current_step"`                                        // 当前执行步骤
	TaskSteps      []TaskStep  `json:"task_steps" gorm:"column:task_steps;serializer:json"` // 执行步骤记录
	UserID         uint        `json:"user_id" gorm:"index"`                                // 列表/统计按用户过滤，加索引避免全表扫
	// 执行者心跳：worker 执行期间周期性刷新。服务启动时的残留清理只处理心跳超时的
	// running 任务，避免把另一个进程正在执行的任务误判为失败（旧进程仍会写回 success）。
	HeartbeatAt *time.Time `json:"heartbeat_at,omitempty" gorm:"index"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// NormalizePages 入库前补齐页面身份并按页序排列：
//   - 缺 PageID 按 p1..pN 生成，并保证唯一；
//   - 缺 PageName 回退 Title，再回退「页面 N」；
//   - 缺省入口页时，首张视为入口页。
//
// 兼容旧请求（无页面字段）与单图任务，保证多页生成/回显始终拿到完整元信息。
func (t *Task) NormalizePages() {
	if len(t.Images) == 0 {
		return
	}
	sort.SliceStable(t.Images, func(i, j int) bool {
		return t.Images[i].SortOrder < t.Images[j].SortOrder
	})
	seen := make(map[string]bool, len(t.Images))
	for i := range t.Images {
		img := &t.Images[i]
		if img.SortOrder <= 0 {
			img.SortOrder = i + 1
		}
		if img.PageID == "" {
			img.PageID = fmt.Sprintf("p%d", i+1)
		}
		for seen[img.PageID] { // 外部传入重复 ID 时追加后缀消重
			img.PageID += "x"
		}
		seen[img.PageID] = true
		if img.PageName == "" {
			img.PageName = img.Title
		}
		if img.PageName == "" {
			img.PageName = fmt.Sprintf("页面 %d", i+1)
		}
	}
	hasEntry := false
	for _, img := range t.Images {
		if img.IsEntry {
			hasEntry = true
			break
		}
	}
	if !hasEntry {
		t.Images[0].IsEntry = true
	}
	// 路由兜底：入口页 /，其余 /<page_id>（供合成预览导航与详情页回显）
	for i := range t.Images {
		if t.Images[i].Route == "" {
			if t.Images[i].IsEntry {
				t.Images[i].Route = "/"
			} else {
				t.Images[i].Route = "/" + t.Images[i].PageID
			}
		}
	}
}

// ── 任务状态查询响应 ──────────────────────────

// TaskStatusResponse GET /api/v1/task/:id 的完整响应
//
//	task_id       - 任务唯一标识
//	status        - pending | running | success | failed
//	target        - flutter | vue3 | react
//	progress      - 执行进度 0-100
//	current_step  - 当前执行步骤，如 VisionAnalyzeSkill / FlutterGenerateSkill
//	can_sse       - 是否可以连接 SSE（true=运行中可连，false=已结束无需连）
//	result        - 生成结果（仅 success 时有值）
//	images        - 上传的图片列表（刷新后回显用）
type TaskStatusResponse struct {
	TaskID      string      `json:"task_id"`
	Target      string      `json:"target"`
	Platform    string      `json:"platform"`
	Status      TaskStatus  `json:"status"`
	Progress    int         `json:"progress"`
	CurrentStep string      `json:"current_step"`
	CanSSE      bool        `json:"can_sse"`
	Images      []ImageItem `json:"images"`
	PageLinks   []PageLink  `json:"page_links"` // 页面跳转边（详情页回显页面关系）
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
		PageLinks:   t.PageLinks,
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
