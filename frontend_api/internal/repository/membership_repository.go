package repository

import (
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"log"

	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════
//  Membership Repository - 会员套餐 / 积分包
//  提供：首次启动灌入默认数据（AutoMigrate 只建表，不建数据）
// ═══════════════════════════════════════════════

// MembershipRepository 会员基础数据存储
type MembershipRepository struct{}

// NewMembershipRepository 创建会员基础数据存储
func NewMembershipRepository() *MembershipRepository {
	return &MembershipRepository{}
}

// SeedIfEmpty 表为空时灌入默认套餐与积分包（幂等，已有数据时完全跳过）
//
// 与 db/migrations/0002_seed_membership_plans.sql 保持同一份数据：
// 该 SQL 供生产环境手工执行/核对，本方法负责新环境首次启动自动初始化。
func (r *MembershipRepository) SeedIfEmpty() error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var planCnt int64
		if err := tx.Model(&model.MembershipPlan{}).Count(&planCnt).Error; err != nil {
			return err
		}
		if planCnt == 0 {
			plans := defaultMembershipPlans()
			if err := tx.Create(&plans).Error; err != nil {
				return err
			}
			log.Printf("[Membership] 已初始化会员套餐 %d 条", len(plans))
		}

		var pkgCnt int64
		if err := tx.Model(&model.CreditsPackage{}).Count(&pkgCnt).Error; err != nil {
			return err
		}
		if pkgCnt == 0 {
			pkgs := defaultCreditsPackages()
			if err := tx.Create(&pkgs).Error; err != nil {
				return err
			}
			log.Printf("[Membership] 已初始化积分包 %d 条", len(pkgs))
		}
		return nil
	})
}

// defaultMembershipPlans 默认套餐
// level 必须与 handler 的 planFeatures 对应：0=免费 1=专业 2=旗舰
func defaultMembershipPlans() []model.MembershipPlan {
	return []model.MembershipPlan{
		{Level: model.LevelFree, Name: "免费版", CreditsPerMonth: 100, Price: 0, OriginalPrice: 0, IsActive: true, SortOrder: 0},
		{Level: model.LevelPro, Name: "专业版", CreditsPerMonth: 500, Price: 29, OriginalPrice: 59, IsActive: true, SortOrder: 1},
		{Level: model.LevelPremium, Name: "旗舰版", CreditsPerMonth: 2000, Price: 99, OriginalPrice: 199, IsActive: true, SortOrder: 2},
	}
}

// defaultCreditsPackages 默认积分包
func defaultCreditsPackages() []model.CreditsPackage {
	return []model.CreditsPackage{
		{Name: "100 次加油包", Credits: 100, Price: 9, OriginalPrice: 12},
		{Name: "500 次加油包", Credits: 500, Price: 39, OriginalPrice: 59},
		{Name: "2000 次加油包", Credits: 2000, Price: 128, OriginalPrice: 199},
	}
}
