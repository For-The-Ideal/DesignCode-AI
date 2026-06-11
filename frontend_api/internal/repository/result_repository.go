package repository

import (
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"log"
)

// ═══════════════════════════════════════════════
//  Result Repository - 生成结果持久化
// ═══════════════════════════════════════════════

// ResultRepository 结果存储
type ResultRepository struct{}

// NewResultRepository 创建结果存储
func NewResultRepository() *ResultRepository {
	return &ResultRepository{}
}

// Create 创建结果
func (r *ResultRepository) Create(result *model.Result) error {
	db := mysql.GetDB()
	if db == nil {
		log.Printf("[ResultRepo] DB not initialized, skipping create for task: %s", result.TaskID)
		return nil
	}
	return db.Create(result).Error
}

// GetByTaskID 根据任务 ID 获取结果
func (r *ResultRepository) GetByTaskID(taskID string) (*model.Result, error) {
	db := mysql.GetDB()
	if db == nil {
		return nil, nil
	}
	var result model.Result
	err := db.Where("task_id = ?", taskID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}
