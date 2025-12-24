package models

import (
	"time"
)

type Transfer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	FileID    uint      `gorm:"index;not null" json:"file_id"`
	Type      string    `gorm:"size:20;not null" json:"type"` // upload/download
	Status    string    `gorm:"size:20;default:'pending'" json:"status"`
	Progress  int       `gorm:"default:0" json:"progress"`
	Speed     int64     `json:"speed"` // bytes per second
	IpAddress string    `gorm:"size:50" json:"ip_address"`
	UserAgent string    `gorm:"size:255" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	File File `gorm:"foreignKey:FileID" json:"file,omitempty"`
}