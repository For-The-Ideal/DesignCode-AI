package worker

import (
	"context"
	"fmt"
	"frontend_api/internal/model"
	"frontend_api/internal/queue"
	"frontend_api/internal/repository"
	"frontend_api/internal/skills/code_to_html"
	"frontend_api/internal/skills/vision"
	"frontend_api/internal/sse"
	"frontend_api/internal/workflow"
	"frontend_api/pkg/logger"
	"frontend_api/utils"
	"runtime/debug"
	"time"
)

// ═══════════════════════════════════════════════
//  GenerateWorker
//  职责：消费任务队列，驱动 GenerateUIWorkflow
// ═══════════════════════════════════════════════

// GenerateWorker 生成 worker
type GenerateWorker struct {
	queue    queue.Queue
	taskRepo *repository.TaskRepository
	workflow *workflow.GenerateUIWorkflow
	log      *logger.Logger
}

// NewGenerateWorker 创建 worker
func NewGenerateWorker(
	q queue.Queue,
	taskRepo *repository.TaskRepository,
	visionSkill *vision.Skill,
	codeToHtmlSkill *code_to_html.Skill,
	resultRepo *repository.ResultRepository,
	sseManager *sse.Manager,
) *GenerateWorker {
	return &GenerateWorker{
		queue:    q,
		taskRepo: taskRepo,
		workflow: workflow.NewGenerateUIWorkflow(visionSkill, codeToHtmlSkill, taskRepo, resultRepo, sseManager),
		log:      logger.NewLogger("worker"),
	}
}

// Run 启动 worker（阻塞，在 goroutine 中调用）
func (w *GenerateWorker) Run(ctx context.Context) {
	w.log.Info("[Worker] 启动，等待任务...")

	for {
		select {
		case <-ctx.Done():
			w.log.Info("[Worker] 收到退出信号")
			return
		default:
			taskID, err := w.queue.Dequeue()
			if err != nil || taskID == "" {
				continue
			}

			w.processTask(ctx, taskID)
		}
	}
}

// processTask 执行单个任务。
// gin 的 HTTP goroutine 有 Recovery 保护，worker goroutine 没有：
// 任何未预期 panic（如类型断言失败）都会直接打崩整个进程，
// 因此这里统一 recover 并复用 workflow.Fail 完成退积分 / 落库 / SSE 通知。
func (w *GenerateWorker) processTask(ctx context.Context, taskID string) {
	var task *model.Task

	defer func() {
		if r := recover(); r != nil {
			w.log.Errorf("[Worker] 任务 panic: %s, panic=%v\n%s", taskID, r, debug.Stack())
			if task != nil {
				w.workflow.Fail(task, fmt.Errorf("任务异常中断: %v", r))
			}
		}
	}()

	w.log.Infof("[Worker] 消费任务: %s（堆内存 %.1f MB）", taskID, utils.HeapMB())

	// 获取任务
	loaded, err := w.taskRepo.GetByID(taskID)
	if err != nil || loaded == nil {
		w.log.Errorf("[Worker] 获取任务失败: %s, err=%v", taskID, err)
		return
	}
	task = loaded

	// 执行期间周期刷新心跳：服务启动时的残留清理只处理心跳超时的 running 任务，
	// 这样另一个进程（旧实例/野进程）不会把本进程正在跑的任务误判为失败。
	hbCtx, hbCancel := context.WithCancel(ctx)
	defer hbCancel()
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-hbCtx.Done():
				return
			case <-ticker.C:
				if err := w.taskRepo.TouchHeartbeat(taskID); err != nil {
					w.log.Errorf("[Worker] 心跳刷新失败: %s, err=%v", taskID, err)
				}
			}
		}
	}()

	// 执行工作流
	start := time.Now()
	if err := w.workflow.Execute(ctx, task); err != nil {
		w.log.Errorf("[Worker] 工作流执行失败: %s, err=%v", taskID, err)
	}
	// 堆内存观测：用于评估 worker 扩容空间（峰值取整任务期间的 After 值，受 GC 时机影响）
	w.log.Infof("[Worker] 任务 %s 结束（耗时 %s，堆内存 %.1f MB）", taskID, time.Since(start).Round(time.Second), utils.HeapMB())
}
