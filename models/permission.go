package models

import (
	"time"
)

type Permission struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Code        string    `gorm:"size:100;not null;uniqueIndex" json:"code"`
	Description string    `gorm:"size:255" json:"description,omitempty"`
	CreatedAt   time.Time `json:"createTime"`
	UpdatedAt   time.Time `json:"updateTime"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
