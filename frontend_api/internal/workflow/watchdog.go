package workflow

import (
	"errors"
	"fmt"
)

// staleTaskErr 看门狗判定中断的原因：执行进程失联（心跳超时），
// 与任务级超时（worker 的 context.WithTimeout）区分开。
var staleTaskErr = errors.New("任务执行中断（执行进程失联），已自动结束")

// ReapStaleTasks 运行期看门狗：扫描心跳超时的 running 任务，按失败收尾。
//
// 服务启动时的 ResetInterruptedTasks 只覆盖「重启前残留」；服务运行期间 worker
// 若异常死亡（未预期的崩溃路径 / OOM 被杀），任务会永远停在 running，前端一直
// 显示“生成中”。这里周期性复用 handleError 完成退积分、失败结果与 SSE 通知。
// 返回本次处理的任务数。
func (w *GenerateUIWorkflow) ReapStaleTasks() int {
	tasks, err := w.taskRepo.ListStaleRunning()
	if err != nil {
		w.log.Errorf("[Workflow] 看门狗查询心跳超时任务失败: %v", err)
		return 0
	}
	for i := range tasks {
		task := tasks[i]
		last := "null"
		if task.HeartbeatAt != nil {
			last = task.HeartbeatAt.Format("2006-01-02 15:04:05")
		}
		w.log.Warn(fmt.Sprintf("[Workflow] 看门狗结束心跳超时任务: %s（最后心跳 %s）", task.ID, last))
		w.handleError(&task, staleTaskErr)
	}
	return len(tasks)
}
