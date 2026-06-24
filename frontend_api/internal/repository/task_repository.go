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
