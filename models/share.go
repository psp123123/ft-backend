package models

import (
	"time"
)

type Share struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FileID    uint      `gorm:"index;not null" json:"file_id"`
	ShareKey  string    `gorm:"uniqueIndex;size:64;not null" json:"share_key"`
	ExpiresAt time.Time `json:"expires_at"`
	AccessCount int     `gorm:"default:0" json:"access_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	File File `gorm:"foreignKey:FileID" json:"file,omitempty"`
}