package worker

import (
	"context"
	"frontend_api/internal/queue"
	"frontend_api/internal/repository"
	"frontend_api/internal/skills/vision"
	"frontend_api/internal/sse"
	"frontend_api/internal/workflow"
	"frontend_api/pkg/logger"
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
	resultRepo *repository.ResultRepository,
	sseManager *sse.Manager,
) *GenerateWorker {
	return &GenerateWorker{
		queue:    q,
		taskRepo: taskRepo,
		workflow: workflow.NewGenerateUIWorkflow(visionSkill, taskRepo, resultRepo, sseManager),
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

			w.log.Infof("[Worker] 消费任务: %s", taskID)

			// 获取任务
			task, err := w.taskRepo.GetByID(taskID)
			if err != nil || task == nil {
				w.log.Errorf("[Worker] 获取任务失败: %s, err=%v", taskID, err)
				continue
			}

			// 执行工作流
			if err := w.workflow.Execute(ctx, task); err != nil {
				w.log.Errorf("[Worker] 工作流执行失败: %s, err=%v", taskID, err)
			}
		}
	}
}
