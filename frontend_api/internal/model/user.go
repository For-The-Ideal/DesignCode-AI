package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Email               string     `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password            string     `gorm:"size:255;not null" json:"-"`
	Nickname            string     `gorm:"size:50" json:"nickname"`
	Avatar              string     `gorm:"size:255" json:"avatar"`
	Phone               string     `gorm:"size:20" json:"phone"`
	Status              string     `gorm:"size:20;default:active" json:"status"`
	Credits             int        `gorm:"default:100" json:"credits"`
	CreditsUsed         int        `gorm:"default:0" json:"credits_used"`
	Level               int        `gorm:"default:0" json:"level"`
	ResetToken          *string    `gorm:"size:255;uniqueIndex" json:"-"`
	ResetTokenExpiresAt *time.Time `gorm:"default:null" json:"-"`
	ResetTokenUsedAt    *time.Time `gorm:"default:null" json:"-"`
}
