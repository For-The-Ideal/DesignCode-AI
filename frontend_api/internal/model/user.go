package model

import "gorm.io/gorm"

// User 用户模型
type User struct {
	gorm.Model
	Email    string `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Nickname string `gorm:"size:50" json:"nickname"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Phone    string `gorm:"size:20" json:"phone"`
	Status   string `gorm:"size:20;default:active" json:"status"` // active | disabled
}
