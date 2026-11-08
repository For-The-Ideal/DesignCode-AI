package repository

import (
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"log"
	"time"

	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════
//  Task Repository - 任务数据持久化
// ═══════════════════════════════════════════════

// TaskRepository 任务存储
type TaskRepository struct{}

// NewTaskRepository 创建任务存储
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

// Create 创建任务
func (r *TaskRepository) Create(task *model.Task) error {
	db := mysql.GetDB()
	if db == nil {
		log.Printf("[TaskRepo] DB not initialized, skipping create: %s", task.ID)
		return nil
	}
	return db.Create(task).Error
}

// GetByID 根据 ID 获取任务
func (r *TaskRepository) GetByID(id string) (*model.Task, error) {
	db := mysql.GetDB()
	if db == nil {
		return nil, nil
	}
	var task model.Task
	err := db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// staleHeartbeat 心跳超时阈值：超过该时长未上报心跳的 running 任务才视为已中断。
// worker 执行期间每 20 秒刷新一次，正常任务不会命中该阈值。
const staleHeartbeat = 3 * time.Minute

// UpdateStatus 更新任务状态（顺带刷新心跳：running 一开始就应有存活证据，
// 否则首次进度上报前的窗口里会被其他实例的残留清理当作已中断）
func (r *TaskRepository) UpdateStatus(id string, status model.TaskStatus) error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	return db.Model(&model.Task{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"heartbeat_at": time.Now(),
		}).Error
}

// TouchHeartbeat 刷新任务执行心跳（worker 执行期间周期性调用）
func (r *TaskRepository) TouchHeartbeat(id string) error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	return db.Model(&model.Task{}).Where("id = ?", id).
		Update("heartbeat_at", time.Now()).Error
}

// UpdateProgress 更新任务进度
//
//	progress    - 当前进度 0-100
//	currentStep - 当前执行步骤，如 "VisionAnalyzeSkill"
func (r *TaskRepository) UpdateProgress(id string, progress int, currentStep string) error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	return db.Model(&model.Task{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"progress":     progress,
			"current_step": currentStep,
			"heartbeat_at": time.Now(),
		}).Error
}

// UpdateProgressAndSteps 只更新 progress/current_step/task_steps 三列。
// 用于 pushProgress 高频进度落库：避免 db.Save(task) 重写包括 images 在内的整行（写放大）。
func (r *TaskRepository) UpdateProgressAndSteps(task *model.Task) error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	// 进度落库同时刷新心跳：详情页的每一项进度都是“这个执行者还活着”的证据
	now := time.Now()
	task.HeartbeatAt = &now
	return db.Model(task).
		Select("progress", "current_step", "task_steps", "heartbeat_at").
		Updates(task).Error
}

// Save 完整保存Task（含 task_steps JSON）
func (r *TaskRepository) Save(task *model.Task) error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	return db.Save(task).Error
}

// Delete 删除任务
func (r *TaskRepository) Delete(id string) error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	return db.Where("id = ?", id).Delete(&model.Task{}).Error
}

// DeleteWithResults 删除任务的同时删除其关联生成结果（同一事务，避免 results 孤儿记录）
func (r *TaskRepository) DeleteWithResults(id string) error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("task_id = ?", id).Delete(&model.Result{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&model.Task{}).Error
	})
}

// GetByUserID 根据用户 ID 获取所有任务（按创建时间倒序）
func (r *TaskRepository) GetByUserID(userID uint) ([]model.Task, error) {
	db := mysql.GetDB()
	if db == nil {
		return nil, nil
	}
	var tasks []model.Task
	err := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// CountByUserID 按状态统计用户任务数量
func (r *TaskRepository) CountByUserID(userID uint) (*model.TaskCountResponse, error) {
	db := mysql.GetDB()
	if db == nil {
		return &model.TaskCountResponse{}, nil
	}

	var counts []struct {
		Status model.TaskStatus
		Count  int
	}
	err := db.Model(&model.Task{}).
		Select("status, count(*) as count").
		Where("user_id = ?", userID).
		Group("status").
		Find(&counts).Error
	if err != nil {
		return nil, err
	}

	resp := &model.TaskCountResponse{}
	for _, c := range counts {
		switch c.Status {
		case model.TaskStatusPending:
			resp.Pending = c.Count
		case model.TaskStatusRunning:
			resp.Running = c.Count
		case model.TaskStatusSuccess:
			resp.Success = c.Count
		case model.TaskStatusFailed:
			resp.Failed = c.Count
		}
	}
	return resp, nil
}

// ResetInterruptedTasks 服务启动时恢复残留任务，并退还预扣积分。
//
// 只清理确定已中断的任务：
//   - pending：只存在于内存队列，进程重启后不可能再被消费；
//   - running 且心跳超时：执行者已消失。心跳仍新鲜的 running 任务会跳过，
//     避免把另一个进程（旧实例/野进程）正在执行的任务误判为失败——
//     否则它会先被写成 failed，旧进程跑完又写回 success，前端表现为
//     “生成中突然失败，回来又莫名其妙成功”。
func (r *TaskRepository) ResetInterruptedTasks() (int64, error) {
	db := mysql.GetDB()
	if db == nil {
		return 0, nil
	}

	// 条件必须对查询与更新保持一致，避免退积分与置 failed 的目标集不同
	cond := "status = ? OR (status = ? AND (heartbeat_at IS NULL OR heartbeat_at < ?))"
	cut := time.Now().Add(-staleHeartbeat)
	args := []interface{}{model.TaskStatusPending, model.TaskStatusRunning, cut}

	var count int64
	err := db.Transaction(func(tx *gorm.DB) error {
		var tasks []model.Task
		if err := tx.Where(cond, args...).Find(&tasks).Error; err != nil {
			return err
		}

		// 逐个退还预扣积分（与状态更新处于同一事务，避免重复退款）
		for _, task := range tasks {
			if task.RequiredPoints <= 0 {
				continue
			}
			if err := tx.Model(&model.User{}).Where("id = ?", task.UserID).
				Updates(map[string]interface{}{
					"credits":      gorm.Expr("credits + ?", task.RequiredPoints),
					"credits_used": gorm.Expr("credits_used - ?", task.RequiredPoints),
				}).Error; err != nil {
				return err
			}
		}

		result := tx.Model(&model.Task{}).Where(cond, args...).
			Updates(map[string]interface{}{
				"status":       model.TaskStatusFailed,
				"current_step": "服务重启，任务中断",
			})
		if result.Error != nil {
			return result.Error
		}
		count = result.RowsAffected
		return nil
	})
	return count, err
}
