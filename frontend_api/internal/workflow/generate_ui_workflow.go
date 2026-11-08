package workflow

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"frontend_api/internal/model"
	"frontend_api/internal/quality"
	"frontend_api/internal/repository"
	"frontend_api/internal/skills/code_to_html"
	"frontend_api/internal/skills/generator"
	"frontend_api/internal/skills/vision"
	"frontend_api/internal/sse"
	"frontend_api/pkg/ai"
	"frontend_api/pkg/logger"
	"frontend_api/pkg/mysql"
	"frontend_api/utils"
	"io"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
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

// pushProgress 推送进度事件并持久化到 DB（含 task_steps）
func (w *GenerateUIWorkflow) pushProgress(task *model.Task, step string, progress int) {
	w.log.Infof("[Workflow] pushProgress: task=%s step=%s progress=%d", task.ID, step, progress)
	// 记录步骤到 task_steps 并落库
	task.AddStep(step, progress, "done")
	task.Progress = progress
	task.CurrentStep = step
	_ = w.taskRepo.UpdateProgressAndSteps(task)
	// SSE 推送（任务级通道：详情页）
	w.sseManager.Push(task.ID, sse.SSEEvent{
		Event: "progress",
		Data:  fmt.Sprintf(`{"progress":%d,"step":"%s"}`, progress, step),
	})
	// 用户级通道：任务列表页实时进度（零轮询，与 task_status 共用同一条长连接）
	progPayload, _ := json.Marshal(map[string]interface{}{
		"task_id":      task.ID,
		"status":       task.Status,
		"progress":     progress,
		"current_step": step,
	})
	w.sseManager.PushUser(fmt.Sprintf("%d", task.UserID), sse.SSEEvent{
		Event: "task_progress",
		Data:  string(progPayload),
	})
}

// pushUserStatus 向用户级 SSE 推送任务状态变更（携带权威计数，前端直接更新，无需回拉 /task/status）
func (w *GenerateUIWorkflow) pushUserStatus(task *model.Task, status model.TaskStatus) {
	counts, _ := w.taskRepo.CountByUserID(task.UserID)
	if counts == nil {
		counts = &model.TaskCountResponse{}
	}
	payload, _ := json.Marshal(map[string]interface{}{
		"task_id": task.ID,
		"status":  status,
		"counts":  counts,
	})
	w.sseManager.PushUser(fmt.Sprintf("%d", task.UserID), sse.SSEEvent{
		Event: "task_status",
		Data:  string(payload),
	})
}

// updateProgressDB 只更新 DB 的 progress + current_step（不追加 task_steps，用于中间态）。
// 同时推任务级与用户级 SSE：详情页只订阅用户级 task_progress，只推任务级时
// 5→20、80→95 这些平滑进度在页面上不可见，会误以为后端没在跑。
func (w *GenerateUIWorkflow) updateProgressDB(task *model.Task, step string, progress int) {
	_ = w.taskRepo.UpdateProgress(task.ID, progress, step)
	w.sseManager.Push(task.ID, sse.SSEEvent{
		Event: "progress",
		Data:  fmt.Sprintf(`{"progress":%d,"step":"%s"}`, progress, step),
	})
	payload, _ := json.Marshal(map[string]interface{}{
		"task_id":      task.ID,
		"status":       task.Status,
		"progress":     progress,
		"current_step": step,
	})
	w.sseManager.PushUser(fmt.Sprintf("%d", task.UserID), sse.SSEEvent{
		Event: "task_progress",
		Data:  string(payload),
	})
}

// simulateProgress 模拟进度递增推送（step 内渐进）
// 从 from 到 to，每次间隔 delay，最后精准推送 to%
func (w *GenerateUIWorkflow) simulateProgress(ctx context.Context, task *model.Task, step string, from, to, delay int) {
	stepSize := 3
	if to-from < stepSize*2 {
		stepSize = 1
	}
	for p := from + stepSize; p < to; p += stepSize {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(delay) * time.Millisecond):
			w.updateProgressDB(task, step, p)
		}
	}
	// 确保终点精准
	w.updateProgressDB(task, step, to)
}

// tickProgress 长时间阻塞阶段的“心跳进度”：每 3 秒推进 stepPct%，封顶 to-1。
// 非流式 AI 调用（38→71、71→80）期间没有任何真实进度事件，前端会长时间停在
// 同一个数字上，看起来像任务已经死了；这里补一条单调递增的可见进度。
// 返回的 stop 必须在该阶段结束时调用，返回值是最后一次推送的进度，
// 供真实回填时作为起点，避免进度倒退。
func (w *GenerateUIWorkflow) tickProgress(ctx context.Context, task *model.Task, step string, from, to, stepPct int) func() int {
	done := make(chan struct{})
	var once sync.Once
	var mu sync.Mutex
	p := from
	w.log.Infof("[Workflow] 心跳进度开始: task=%s step=%s %d→%d（每 3s +%d%%）", task.ID, step, from, to, stepPct)
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			case <-ticker.C:
				mu.Lock()
				if p+stepPct >= to {
					mu.Unlock()
					return // 已推到上限（to-1），等真实回填
				}
				p += stepPct
				cur := p
				mu.Unlock()
				w.updateProgressDB(task, step, cur)
			}
		}
	}()
	return func() int {
		once.Do(func() { close(done) })
		mu.Lock()
		defer mu.Unlock()
		w.log.Infof("[Workflow] 心跳进度结束: task=%s step=%s 停在 %d%%", task.ID, step, p)
		return p
	}
}

// progressTo 单页任务保留平滑模拟；多页任务直接推送里程碑，
// 避免逐页 sleep 累加导致总耗时被拉长。
func (w *GenerateUIWorkflow) progressTo(ctx context.Context, task *model.Task, step string, from, to int, multiPage bool) {
	if multiPage {
		w.pushProgress(task, step, to)
		return
	}
	if from < 0 || from > to {
		from = to
	}
	w.simulateProgress(ctx, task, step, from, to, 150)
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
	task.Status = model.TaskStatusRunning
	w.pushUserStatus(task, model.TaskStatusRunning)

	// ── Step 1: 准备图片（按页下载 + 转 base64）──
	w.pushProgress(task, StepDownloadImages, 5)
	pages := OrderedPages(task.Images)
	multiPage := len(pages) > 1
	w.log.Infof("[Workflow] step 1/6: 准备图片（%d 页）...", len(pages))
	// 逐页下载并转 base64：多页任务每页生成只带自己的图片，避免输入膨胀
	b64ByPage := make(map[string]string, len(pages))
	for i, img := range pages {
		w.log.Infof("[Workflow]   页面[%d] %s: %s", i+1, img.PageID, img.URL)
		if b64, err := urlToBase64(img.URL); err == nil {
			b64ByPage[img.PageID] = b64
			w.log.Infof("[Workflow]   页面[%d] base64 转换完成: %d bytes", i+1, len(b64))
		} else {
			w.log.Infof("[Workflow]   [WARN] 页面[%d] 转 base64 失败: %v", i+1, err)
		}
	}
	w.simulateProgress(ctx, task, StepDownloadImages, 5, 20, 400)

	// ── Step 2~4: 逐页生成（DSL → 代码 → 预览）──
	// 单页任务与旧行为完全等价（走同一条逐页路径）；多页任务每页独立生成，
	// 进度 20% → 80% 按页均分，current_step 带 (i/N) 便于前端展示当前页。
	genStep := getGeneratorStepName(task.Target)
	genAI, err := ai.SelectClient("write")
	if err != nil {
		w.handleError(task, fmt.Errorf("获取AI客户端失败: %w", err))
		return err
	}
	genSkill := generator.GeneratorFunc(task.Target, genAI)

	pageResults := make([]model.PageResult, 0, len(pages))
	total := len(pages)
	for i, page := range pages {
		segStart := 20 + 60*i/total
		segEnd := 20 + 60*(i+1)/total
		dslEnd := segStart + (segEnd-segStart)*30/100
		genEnd := segStart + (segEnd-segStart)*85/100
		suffix := ""
		if multiPage {
			suffix = fmt.Sprintf(" %d/%d", i+1, total)
		}
		visionStep := StepVisionAnalyzeSkill + suffix
		pageGenStep := genStep + suffix

		// ── 2a. DSL：有标注树本地构建，无树对该页图片单独跑 Vision ──
		w.pushProgress(task, visionStep, segStart)
		dsl, ok := BuildPageDSL(page)
		if ok {
			w.log.Infof("[Workflow] 页面 %s(%s): 使用结构标注生成 DSL（跳过 Vision）", page.PageID, page.PageName)
		} else {
			w.log.Infof("[Workflow] 页面 %s(%s): 无结构标注，回退视觉分析", page.PageID, page.PageName)
			visionOutput, verr := w.visionSkill.Execute(ctx, vision.Input{Images: []string{page.URL}})
			if verr != nil {
				w.handleError(task, fmt.Errorf("第 %d 页视觉分析失败: %w", i+1, verr))
				return verr
			}
			dsl = visionOutput.(vision.Output).DSL
		}
		w.progressTo(ctx, task, visionStep, segStart, dslEnd, multiPage)

		// ── 3. 代码生成：每页只带自己的图片 ──
		w.pushProgress(task, pageGenStep, dslEnd)
		pageImages := make([]string, 0, 1)
		if b64 := b64ByPage[page.PageID]; b64 != "" {
			pageImages = append(pageImages, b64)
		}
		genInput := generator.Input{
			DSL:      dsl,
			Images:   pageImages,
			Platform: task.Platform,
			Options:  task.Options,
		}
		if task.ComponentLib != "" {
			genInput.ComponentLib = task.ComponentLib
		}
		// 非流式生成期间没有真实进度：用心跳进度让前端看到“还在跑”。
		// stopTick 返回最后推送值，真实回填时从它继续，避免进度倒退。
		stopTick := w.tickProgress(ctx, task, pageGenStep, dslEnd, genEnd, 1)
		genOutput, gerr := genSkill.Execute(ctx, genInput)
		ticked := stopTick()
		if gerr != nil {
			w.handleError(task, fmt.Errorf("第 %d 页代码生成失败: %w", i+1, gerr))
			return gerr
		}
		out := genOutput.(generator.Output)
		// 清洗 markdown 代码块包裹 (```dart ... ``` → 纯代码)
		out.Code = cleanCodeBlock(out.Code)
		// 空结果拦截：模型偶尔会因思维链(reasoning_content)吃满 max_tokens 而返回空 content。
		// 此时不得当作成功入库（否则出现 code="" 的 success 记录），直接判失败并退积分。
		if strings.TrimSpace(out.Code) == "" {
			emptyErr := fmt.Errorf("第 %d 页代码生成结果为空（模型未返回有效代码，疑似思维链耗尽 max_tokens）", i+1)
			w.handleError(task, emptyErr)
			return emptyErr
		}
		w.log.Infof("[Workflow] 页面 %s 生成完成, code=%d chars", page.PageID, len(out.Code))

		// 确定性质量检查：规则编号与 prompts/*/generate.txt 的硬性规则一一对应。
		// 只做扫描与标记，不阻断生成（fail 由前端展示，不在这里回滚任务）。
		report := quality.Scan(out.Code, task.Target, dsl, task.ComponentLib)
		if report.Status != "pass" {
			w.log.Infof("[Workflow] 页面 %s 质量检查: %s (%d/%d 通过)", page.PageID, report.Status, report.Passed, report.Total)
		}

		// 单页任务保留逐行推代码的生成动画；多页任务不推，避免多段代码在前端拼接混杂
		if !multiPage {
			w.streamCode(ctx, task.ID, out.Code)
		}
		w.progressTo(ctx, task, pageGenStep, ticked, genEnd, multiPage)

		// ── 4. 转 HTML 预览（失败降级用模型自带 preview）──
		htmlStep := pageGenStep
		if !multiPage {
			// 单页任务保留独立的 CodeToHtmlSkill 阶段，与旧版步骤序列一致
			htmlStep = StepConvertPreview
			w.pushProgress(task, htmlStep, genEnd)
		}
		stopHTML := w.tickProgress(ctx, task, htmlStep, genEnd, segEnd, 1)
		htmlOutput, herr := w.codeToHtmlSkill.Execute(ctx, code_to_html.Input{
			SourceCode: out.Code,
			Framework:  task.Target,
			DSL:        dsl,
			Images:     page.URL,
		})
		htmlTicked := stopHTML()
		previewHTML := out.Preview
		if herr != nil {
			w.log.Errorf("[Workflow] 页面 %s 转 HTML 失败: %v（降级用模型自带 preview）", page.PageID, herr)
		} else {
			previewHTML = htmlOutput.(code_to_html.Output).HTML
		}
		route := page.Route
		if route == "" {
			route = "/" + page.PageID
		}
		pageResults = append(pageResults, model.PageResult{
			PageID:        page.PageID,
			PageName:      page.PageName,
			Route:         route,
			DSL:           dsl,
			Code:          out.Code,
			Preview:       previewHTML,
			Status:        "success",
			QualityReport: &report,
		})
		w.progressTo(ctx, task, htmlStep, htmlTicked, segEnd, multiPage)
	}

	// ── Step 5: 合成 + 保存结果 ─────────────────────
	// 入口页（OrderedPages 已前置）作为旧字段 Code/Preview/DSL 的兼容取值，
	// 保证未升级的前端与旧数据读取路径不受影响。
	entry := pageResults[0]
	appPreview := BuildAppPreview(pageResults, entry.PageID)
	navMap := BuildNavMap(task.PageLinks, pageResults)

	w.pushProgress(task, StepSaveResult, 80)
	w.log.Infof("[Workflow] step 5/6: 保存结果（%d 页）...", len(pageResults))
	result := &model.Result{
		TaskID:           task.ID,
		UserID:           task.UserID,
		DSL:              entry.DSL,
		Code:             entry.Code,
		Preview:          entry.Preview,
		CodeType:         task.Target,
		GenerationStatus: "success",
		PointsDeducted:   task.RequiredPoints,
		CreatedAt:        time.Now(),
		Pages:            pageResults,
		AppPreview:       appPreview,
		NavMap:           navMap,
		QualityReport:    entry.QualityReport,
	}
	if err := w.resultRepo.Create(result); err != nil {
		w.log.Errorf("[Workflow] 保存结果失败: %v", err)
	}

	// 推送 preview 事件：入口页 HTML（前端生成动画/预览区即时展示）
	w.sseManager.Push(task.ID, sse.SSEEvent{
		Event: "preview",
		Data:  entry.Preview,
	})

	// 推送 quality 事件：入口页质量检查报告（与 done 解耦，前端可先展示检查结果）
	if entry.QualityReport != nil {
		if payload, err := json.Marshal(entry.QualityReport); err == nil {
			w.sseManager.Push(task.ID, sse.SSEEvent{Event: "quality", Data: string(payload)})
		}
	}

	// 积分已在创建任务时预扣，此处无需重复扣减
	w.simulateProgress(ctx, task, StepSaveResult, 80, 95, 200)

	// ── Step 6: Done ───────────────────────────────
	task.Status = model.TaskStatusSuccess
	_ = w.taskRepo.UpdateStatus(task.ID, model.TaskStatusSuccess)
	w.pushUserStatus(task, model.TaskStatusSuccess)
	w.pushProgress(task, StepDone, 100)

	// 推送 done 事件
	w.sseManager.Push(task.ID, sse.SSEEvent{
		Event: "done",
		Data:  fmt.Sprintf(`{"progress":100,"step":"Done"}`),
	})

	w.log.Infof("[Workflow] 任务完成: %s", task.ID)
	return nil
}

// maxImageBytes 图片下载上限，防止异常大图拖垮内存
const maxImageBytes = 20 << 20 // 20MB

// urlToBase64 从 URL 下载图片，缩放压缩后转为 base64（data URL）
// 原图（设计稿常见 2700px+）直接 base64 会让单任务内存峰值达到几十 MB，
// 缩到长边 1600 并转 JPEG 后通常只剩几百 KB；解码失败（如 webp）时回退原图
func urlToBase64(url string) (string, error) {
	// 防 SSRF：仅限白名单域名（COS 桶 / server.allowed_image_hosts）+ https
	if err := utils.ValidateImageURL(url); err != nil {
		return "", err
	}
	// 使用专用客户端：60s 超时 + 重定向逐跳校验
	resp, err := utils.ImageHTTPClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxImageBytes))
	if err != nil {
		return "", fmt.Errorf("读取失败: %w", err)
	}

	// 优先压缩：解码 → 缩放到长边上限 → JPEG 编码（公用实现见 utils/helps.go）
	if compressed, err := utils.CompressImageToBase64(body); err == nil {
		return compressed, nil
	}

	// 回退：不支持的格式按原样返回，保证生成流程不中断
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/png"
	}
	encoded := base64.StdEncoding.EncodeToString(body)
	return fmt.Sprintf("data:%s;base64,%s", contentType, encoded), nil
}

// Fail 处理失败收尾：退还预扣积分、更新任务状态、保存失败结果、推送 SSE。
// 导出给 worker 的 panic 兜底复用（保证任何异常路径都执行完整善后）。
func (w *GenerateUIWorkflow) Fail(task *model.Task, err error) {
	w.handleError(task, err)
}

func (w *GenerateUIWorkflow) handleError(task *model.Task, err error) {
	// 任务级超时（worker 的 context.WithTimeout）会在各 skill 错误里以
	// context deadline exceeded 形式冒出，统一换成用户能读懂的文案。
	if errors.Is(err, context.DeadlineExceeded) {
		err = fmt.Errorf("任务执行超时，已自动中断（可减少页数或稍后重试）")
	}

	// 先复核任务现状：任务可能已被另一个执行者成功完成（例如启动时的残留清理
	// 误判后旧进程继续跑完），此时不能把 success 覆盖成 failed，也不再向
	// 前端推 failed——否则表现为“生成中突然失败，回来又莫名其妙成功”。
	if latest, gerr := w.taskRepo.GetByID(task.ID); gerr == nil && latest != nil && latest.Status == model.TaskStatusSuccess {
		w.log.Warn(fmt.Sprintf("[Workflow] 任务 %s 已是 success，忽略本次失败: %v", task.ID, err))
		return
	}

	w.log.Errorf("[Workflow] 任务失败: %s, error=%v", task.ID, err)

	// 退还预扣积分（创建任务时已扣）
	if task.RequiredPoints > 0 {
		if db := mysql.GetDB(); db != nil {
			db.Model(&model.User{}).Where("id = ?", task.UserID).
				Updates(map[string]interface{}{
					"credits":      gorm.Expr("credits + ?", task.RequiredPoints),
					"credits_used": gorm.Expr("credits_used - ?", task.RequiredPoints),
				})
		}
	}

	// 更新任务内存状态并落库
	task.Status = model.TaskStatusFailed
	_ = w.taskRepo.Save(task)
	w.pushUserStatus(task, model.TaskStatusFailed)

	// 保存失败结果（含错误信息和积分记录）
	result := &model.Result{
		TaskID:           task.ID,
		UserID:           task.UserID,
		CodeType:         task.Target,
		GenerationStatus: "failed",
		PointsDeducted:   task.RequiredPoints,
		ErrorMsg:         err.Error(),
		CreatedAt:        time.Now(),
	}
	if createErr := w.resultRepo.Create(result); createErr != nil {
		w.log.Errorf("[Workflow] 保存失败结果出错: %v", createErr)
	}

	// 推送 error 事件（用 json.Marshal 构造，避免错误信息中的引号破坏 JSON）
	errPayload, _ := json.Marshal(map[string]string{"message": err.Error()})
	w.sseManager.Push(task.ID, sse.SSEEvent{
		Event: "error",
		Data:  string(errPayload),
	})
}

// cleanCodeBlock 去掉 AI 返回值的 ```language ... ``` 包裹
func cleanCodeBlock(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```") && strings.HasSuffix(raw, "```") {
		if idx := strings.Index(raw, "\n"); idx != -1 {
			raw = raw[idx+1:]
		}
		raw = strings.TrimSuffix(raw, "```")
		raw = strings.TrimSpace(raw)
	}
	return raw
}
