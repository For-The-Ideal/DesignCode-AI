package repository

import (
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"log"
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
