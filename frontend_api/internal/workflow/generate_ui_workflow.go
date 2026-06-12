package workflow

import (
	"context"
	"encoding/base64"
	"fmt"
	"frontend_api/internal/model"
	"frontend_api/internal/repository"
	"frontend_api/internal/skills/code_to_html"
	"frontend_api/internal/skills/generator"
	"frontend_api/internal/skills/vision"
	"frontend_api/internal/sse"
	"frontend_api/pkg/ai"
	"frontend_api/pkg/logger"
	"io"
	"net/http"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════
//  GenerateUIWorkflow
//  流程：
//    创建任务 → 准备图片 → VisionAnalyzeSkill
//    → GeneratorSkill → CodeToHtmlSkill → 保存结果 → SSE通知
//
//  Progress 映射：
//    PrepareImages        0% → 20%
//    VisionAnalyzeSkill  20% → 50%
//    GeneratorSkill      50% → 75%
//    ConvertPreview      75% → 85%
//    SaveResult          85% → 95%
//    Done                100%
//
//  current_step 取值：
//    VisionAnalyzeSkill / FlutterGenerateSkill / Vue3GenerateSkill
//    / ReactGenerateSkill / CodeToHtmlSkill / Done
// ═══════════════════════════════════════════════

// step 常量
const (
	StepDownloadImages     = "PrepareImages"
	StepVisionAnalyzeSkill = "VisionAnalyzeSkill"
	StepGenerate           = "GenerateSkill"
	StepConvertPreview     = "CodeToHtmlSkill"
	StepSaveResult         = "SaveResult"
	StepDone               = "Done"
)

// GenerateUIWorkflow UI 生成工作流
type GenerateUIWorkflow struct {
	visionSkill     *vision.Skill
	codeToHtmlSkill *code_to_html.Skill
	taskRepo        *repository.TaskRepository
	resultRepo      *repository.ResultRepository
	sseManager      *sse.Manager
	log             *logger.Logger
}

// NewGenerateUIWorkflow 创建工作流
func NewGenerateUIWorkflow(
	visionSkill *vision.Skill,
	codeToHtmlSkill *code_to_html.Skill,
	taskRepo *repository.TaskRepository,
	resultRepo *repository.ResultRepository,
	sseManager *sse.Manager,
) *GenerateUIWorkflow {
	return &GenerateUIWorkflow{
		visionSkill:     visionSkill,
		codeToHtmlSkill: codeToHtmlSkill,
		taskRepo:        taskRepo,
		resultRepo:      resultRepo,
		sseManager:      sseManager,
		log:             logger.NewLogger("workflow"),
	}
}

// pushProgress 推送进度事件并持久化到 DB
func (w *GenerateUIWorkflow) pushProgress(taskID string, step string, progress int) {
	w.log.Infof("[Workflow] pushProgress: task=%s step=%s progress=%d", taskID, step, progress)
	// 持久化进度
	_ = w.taskRepo.UpdateProgress(taskID, progress, step)
	// SSE 推送
	w.sseManager.Push(taskID, sse.SSEEvent{
		Event: "progress",
		Data:  fmt.Sprintf(`{"progress":%d,"step":"%s"}`, progress, step),
	})
}

// simulateProgress 模拟进度递增推送（step 内渐进）
// 从 from 到 to，每次间隔 delay，最后精准推送 to%
func (w *GenerateUIWorkflow) simulateProgress(ctx context.Context, taskID string, step string, from, to, delay int) {
	stepSize := 3
	if to-from < stepSize*2 {
		stepSize = 1
	}
	for p := from + stepSize; p < to; p += stepSize {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(delay) * time.Millisecond):
			w.pushProgress(taskID, step, p)
		}
	}
	// 确保终点精确
	w.pushProgress(taskID, step, to)
}

// streamCode 将代码按行拆分为 SSE message 事件逐步推送
func (w *GenerateUIWorkflow) streamCode(ctx context.Context, taskID, code string) {
	if code == "" {
		return
	}
	// 按行拆分，并保留换行符
	lines := strings.SplitAfter(code, "\n")
	// 如果行数太少，直接一次性推送
	if len(lines) <= 3 {
		w.sseManager.Push(taskID, sse.SSEEvent{
			Event: "message",
			Data:  code,
		})
		return
	}

	// 每行拆成一行数据，逐行推送
	for _, line := range lines {
		select {
		case <-ctx.Done():
			return
		case <-time.After(20 * time.Millisecond):
			w.sseManager.Push(taskID, sse.SSEEvent{
				Event: "message",
				Data:  line,
			})
		}
	}
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

	// ── Step 1: 准备图片（取 URL + 下载转 base64）──
	w.pushProgress(task.ID, StepDownloadImages, 5)
	w.log.Infof("[Workflow] step 1/6: 准备图片 (%d 张)...", len(task.Images))
	imageURLs := make([]string, len(task.Images))
	imagesBase64 := make([]string, 0, len(task.Images))
	for i, img := range task.Images {
		imageURLs[i] = img.URL
		w.log.Infof("[Workflow]   图片[%d]: %s", i, img.URL)
		// 下载图片并转为 base64，供 Generator 等通用模型使用
		if b64, err := urlToBase64(img.URL); err == nil {
			imagesBase64 = append(imagesBase64, b64)
			w.log.Infof("[Workflow]   图片[%d] base64 转换完成: %d bytes", i, len(b64))
		} else {
			w.log.Infof("[Workflow]   [WARN] 图片[%d] 转 base64 失败: %v", i, err)
		}
	}
	w.simulateProgress(ctx, task.ID, StepDownloadImages, 5, 20, 400)

	// ── Step 2: VisionAnalyzeSkill → DSL ──────────
	w.pushProgress(task.ID, StepVisionAnalyzeSkill, 20)
	w.log.Infof("[Workflow] step 2/6: 视觉分析 (VisionAnalyzeSkill)...")
	w.simulateProgress(ctx, task.ID, StepVisionAnalyzeSkill, 20, 35, 300)

	visionOutput, err := w.visionSkill.Execute(ctx, vision.Input{Images: imageURLs})
	if err != nil {
		w.handleError(task.ID, fmt.Errorf("视觉分析失败: %w", err))
		return err
	}

	dsl := visionOutput.(vision.Output).DSL
	w.log.Infof("[Workflow] step 2/6 完成, DSL: %s", dsl)
	w.simulateProgress(ctx, task.ID, StepVisionAnalyzeSkill, 35, 50, 200)

	// ── Step 3: GeneratorSkill → Code + Preview ───
	genStep := getGeneratorStepName(task.Target)
	w.pushProgress(task.ID, genStep, 50)
	w.log.Infof("[Workflow] step 3/6: 代码生成 (%s)...", genStep)
	w.simulateProgress(ctx, task.ID, genStep, 50, 65, 300)

	genAI, err := ai.SelectClient("write")
	if err != nil {
		w.handleError(task.ID, fmt.Errorf("获取AI客户端失败: %w", err))
		return err
	}

	genSkill := generator.GeneratorFunc(task.Target, genAI)
	genOutput, err := genSkill.Execute(ctx, generator.Input{
		DSL:    dsl,
		Images: imagesBase64,
	})
	if err != nil {
		w.handleError(task.ID, fmt.Errorf("代码生成失败: %w", err))
		return err
	}

	out := genOutput.(generator.Output)
	w.log.Infof("[Workflow] step 3/6 完成, code=%d chars", len(out.Code))

	// 逐步推送 templateCode（逐行，SSE message 事件）
	w.streamCode(ctx, task.ID, out.Code)
	w.simulateProgress(ctx, task.ID, genStep, 65, 75, 200)

	// ── Step 4: CodeToHtmlSkill → HTML Preview ────
	w.pushProgress(task.ID, StepConvertPreview, 75)
	w.log.Infof("[Workflow] step 4/6: 代码转 HTML 预览 (CodeToHtmlSkill)...")

	// 将图像 URL 拼接成逗号分隔字符串
	imageURLStr := strings.Join(imageURLs, ",")
	htmlOutput, err := w.codeToHtmlSkill.Execute(ctx, code_to_html.Input{
		SourceCode: out.Code,
		Framework:  task.Target,
		DSL:        dsl,
		Images:     imageURLStr,
	})
	var previewHTML string
	if err != nil {
		w.log.Errorf("[Workflow] CodeToHtmlSkill 转换失败: %v", err)
		// 降级：使用 generator 返回的原始 preview
		previewHTML = out.Preview
	} else {
		previewHTML = htmlOutput.(code_to_html.Output).HTML
	}

	w.simulateProgress(ctx, task.ID, StepConvertPreview, 75, 85, 200)

	// 推送 preview 事件（替换之前 Step 3 末推送 out.Preview 的逻辑）
	w.sseManager.Push(task.ID, sse.SSEEvent{
		Event: "preview",
		Data:  previewHTML,
	})

	// ── Step 5: 保存结果到 DB ─────────────────────
	w.pushProgress(task.ID, StepSaveResult, 85)
	w.log.Infof("[Workflow] step 5/6: 保存结果...")
	result := &model.Result{
		TaskID:    task.ID,
		DSL:       dsl,
		Code:      out.Code,
		Preview:   previewHTML, // 使用转换后的 HTML
		Score:     out.Score,
		CodeType:  task.Target,
		CreatedAt: time.Now(),
	}
	if err := w.resultRepo.Create(result); err != nil {
		w.log.Errorf("[Workflow] 保存结果失败: %v", err)
	}
	w.simulateProgress(ctx, task.ID, StepSaveResult, 85, 95, 200)

	// ── Step 6: Done ───────────────────────────────
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

// urlToBase64 从 URL 下载图片并转换为 base64 字符串
func urlToBase64(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取失败: %w", err)
	}

	// 从 Content-Type 推断 MIME 类型
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/png"
	}

	// data:image/png;base64,... 格式
	encoded := base64.StdEncoding.EncodeToString(body)
	return fmt.Sprintf("data:%s;base64,%s", contentType, encoded), nil
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
