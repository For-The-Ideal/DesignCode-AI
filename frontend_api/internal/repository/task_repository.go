package repository

import (
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"log"

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

// UpdateStatus 更新任务状态
func (r *TaskRepository) UpdateStatus(id string, status model.TaskStatus) error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	return db.Model(&model.Task{}).Where("id = ?", id).Update("status", status).Error
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
		}).Error
}

// UpdateProgressAndSteps 只更新 progress/current_step/task_steps 三列。
// 用于 pushProgress 高频进度落库：避免 db.Save(task) 重写包括 images 在内的整行（写放大）。
func (r *TaskRepository) UpdateProgressAndSteps(task *model.Task) error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	return db.Model(task).
		Select("progress", "current_step", "task_steps").
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

// ResetInterruptedTasks 服务启动时恢复残留任务：把 pending/running 标记为 failed 并退还预扣积分。
// pending 任务只存在于内存队列，进程重启后不会再被消费；若不清理会永久占用用户活跃名额
// （CountActive 统计 pending+running），且已预扣的积分无法退回。
func (r *TaskRepository) ResetInterruptedTasks() (int64, error) {
	db := mysql.GetDB()
	if db == nil {
		return 0, nil
	}

	interrupted := []model.TaskStatus{model.TaskStatusPending, model.TaskStatusRunning}
	var count int64
	err := db.Transaction(func(tx *gorm.DB) error {
		var tasks []model.Task
		if err := tx.Where("status IN ?", interrupted).Find(&tasks).Error; err != nil {
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

		result := tx.Model(&model.Task{}).Where("status IN ?", interrupted).
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
