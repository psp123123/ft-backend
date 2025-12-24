package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Password  string    `gorm:"size:100;not null" json:"-"`
	FullName  string    `gorm:"size:100" json:"full_name"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Role      string    `gorm:"size:20;default:'user'" json:"role"`
	CreatedAt time.Time `json:"createTime"`
	UpdatedAt time.Time `json:"updateTime"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	Files     []File     `gorm:"foreignKey:UserID" json:"files,omitempty"`
	Transfers []Transfer `gorm:"foreignKey:UserID" json:"transfers,omitempty"`
}