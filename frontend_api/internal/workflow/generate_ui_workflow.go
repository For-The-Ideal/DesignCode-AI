package workflow

import (
	"context"
	"fmt"
	"frontend_api/internal/model"
	"frontend_api/internal/repository"
	"frontend_api/internal/skills/generator"
	"frontend_api/internal/skills/vision"
	"frontend_api/internal/sse"
	"frontend_api/pkg/logger"
	"time"
)

// ═══════════════════════════════════════════════
//  GenerateUIWorkflow
//  流程：
//    创建任务 → 下载图片 → VisionAnalyzeSkill
//    → GeneratorSkill → 保存结果 → SSE通知
//
//  Progress 映射：
//    DownloadImages       0% → 20%
//    VisionAnalyzeSkill  20% → 50%
//    GeneratorSkill      50% → 80%
//    SaveResult          80% → 90%
//    Done                100%
//
//  current_step 取值：
//    VisionAnalyzeSkill / FlutterGenerateSkill
//    / Vue3GenerateSkill / ReactGenerateSkill / Done
// ═══════════════════════════════════════════════

// step 常量，确保前后端 step 名称一致
const (
	StepDownloadImages     = "DownloadImages"
	StepVisionAnalyzeSkill = "VisionAnalyzeSkill"
	StepGenerate           = "GenerateSkill"
	StepSaveResult         = "SaveResult"
	StepDone               = "Done"
)

// GenerateUIWorkflow UI 生成工作流
type GenerateUIWorkflow struct {
	visionSkill *vision.Skill
	taskRepo    *repository.TaskRepository
	resultRepo  *repository.ResultRepository
	sseManager  *sse.Manager
	log         *logger.Logger
}

// NewGenerateUIWorkflow 创建工作流
func NewGenerateUIWorkflow(
	visionSkill *vision.Skill,
	taskRepo *repository.TaskRepository,
	resultRepo *repository.ResultRepository,
	sseManager *sse.Manager,
) *GenerateUIWorkflow {
	return &GenerateUIWorkflow{
		visionSkill: visionSkill,
		taskRepo:    taskRepo,
		resultRepo:  resultRepo,
		sseManager:  sseManager,
		log:         logger.NewLogger("workflow"),
	}
}

// pushProgress 推送进度事件并持久化到 DB
func (w *GenerateUIWorkflow) pushProgress(taskID string, step string, progress int) {
	// 持久化进度
	_ = w.taskRepo.UpdateProgress(taskID, progress, step)
	// SSE 推送
	w.sseManager.Push(taskID, sse.SSEEvent{
		Event: "progress",
		Data:  fmt.Sprintf(`{"progress":%d,"step":"%s"}`, progress, step),
	})
}

// getGeneratorStepName 根据 target 返回 generator step 名称
func getGeneratorStepName(target string) string {
	switch target {
	case "flutter":
		return "FlutterGenerateSkill"
	case "vue3":
		return "Vue3GenerateSkill"
	case "react":
		return "ReactGenerateSkill"
	default:
		return "GenerateSkill"
	}
}

// Execute 执行完整工作流
func (w *GenerateUIWorkflow) Execute(ctx context.Context, task *model.Task) error {
	w.log.Infof("[Workflow] 开始执行任务: %s, target=%s", task.ID, task.Target)

	// ── 0. 更新状态为 running ──
	_ = w.taskRepo.UpdateStatus(task.ID, model.TaskStatusRunning)

	// ── Step 1: 下载图片 ──────────────────────────
	w.pushProgress(task.ID, StepDownloadImages, 5)
	w.log.Infof("[Workflow] step 1/5: 下载图片 (%d 张)...", len(task.Images))
	for i, img := range task.Images {
		w.log.Infof("[Workflow]   图片[%d]: %s", i, img.URL)
		// TODO: 实际下载图片到临时目录
	}
	w.pushProgress(task.ID, StepDownloadImages, 20)

	// ── Step 2: VisionAnalyzeSkill → DSL ──────────
	w.pushProgress(task.ID, StepVisionAnalyzeSkill, 20)
	w.log.Infof("[Workflow] step 2/5: 视觉分析 (VisionAnalyzeSkill)...")

	visionOutput, err := w.visionSkill.Execute(ctx, vision.Input{Images: task.Images})
	if err != nil {
		w.handleError(task.ID, fmt.Errorf("视觉分析失败: %w", err))
		return err
	}

	dsl := visionOutput.(vision.Output).DSL
	w.log.Infof("[Workflow] step 2/5 完成, DSL: %s", dsl)
	w.pushProgress(task.ID, StepVisionAnalyzeSkill, 50)

	// ── Step 3: GeneratorSkill → Code + Preview ───
	genStep := getGeneratorStepName(task.Target)
	w.pushProgress(task.ID, genStep, 50)
	w.log.Infof("[Workflow] step 3/5: 代码生成 (%s)...", genStep)

	genSkill := generator.GeneratorFunc(task.Target)
	genOutput, err := genSkill.Execute(ctx, generator.Input{DSL: dsl})
	if err != nil {
		w.handleError(task.ID, fmt.Errorf("代码生成失败: %w", err))
		return err
	}

	out := genOutput.(generator.Output)
	w.log.Infof("[Workflow] step 3/5 完成, code=%d chars", len(out.Code))
	w.pushProgress(task.ID, genStep, 80)

	// ── Step 4: 保存结果到 DB ─────────────────────
	w.pushProgress(task.ID, StepSaveResult, 85)
	w.log.Infof("[Workflow] step 4/5: 保存结果...")
	result := &model.Result{
		TaskID:    task.ID,
		DSL:       dsl,
		Code:      out.Code,
		Preview:   out.Preview,
		Score:     out.Score,
		CodeType:  task.Target,
		CreatedAt: time.Now(),
	}
	if err := w.resultRepo.Create(result); err != nil {
		w.log.Errorf("[Workflow] 保存结果失败: %v", err)
	}
	w.pushProgress(task.ID, StepSaveResult, 95)

	// ── Step 5: Done ───────────────────────────────
	_ = w.taskRepo.UpdateStatus(task.ID, model.TaskStatusSuccess)
	w.pushProgress(task.ID, StepDone, 100)

	// 推送 done 事件
	w.sseManager.Push(task.ID, sse.SSEEvent{
		Event: "done",
		Data:  fmt.Sprintf(`{"progress":100,"step":"Done"}`),
	})

	w.log.Infof("[Workflow] 任务完成: %s", task.ID)
	return nil
}

// handleError 处理工作流错误
func (w *GenerateUIWorkflow) handleError(taskID string, err error) {
	w.log.Errorf("[Workflow] 任务失败: %s, error=%v", taskID, err)
	_ = w.taskRepo.UpdateStatus(taskID, model.TaskStatusFailed)

	// 推送 error 事件
	w.sseManager.Push(taskID, sse.SSEEvent{
		Event: "error",
		Data:  fmt.Sprintf(`{"message":"%s"}`, err.Error()),
	})
}
