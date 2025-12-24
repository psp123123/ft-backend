package models

import (
	"time"
)

type File struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	Filename    string    `gorm:"size:255;not null" json:"filename"`
	OriginalName string   `gorm:"size:255;not null" json:"original_name"`
	Size        int64     `json:"size"`
	Path        string    `gorm:"size:255;not null" json:"path"`
	MimeType    string    `gorm:"size:100" json:"mime_type"`
	Extension   string    `gorm:"size:20" json:"extension"`
	Hash        string    `gorm:"size:64" json:"hash"`
	Status      string    `gorm:"size:20;default:'available'" json:"status"`
	Visibility  string    `gorm:"size:20;default:'private'" json:"visibility"`
	DownloadCount int     `gorm:"default:0" json:"download_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`

	// 关联
	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Transfers  []Transfer `gorm:"foreignKey:FileID" json:"transfers,omitempty"`
	Shares     []Share    `gorm:"foreignKey:FileID" json:"shares,omitempty"`
}