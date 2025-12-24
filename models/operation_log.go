package models

import (
	"time"
)

type OperationLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Username      string    `gorm:"size:50;not null" json:"username"`
	Operation     string    `gorm:"size:100;not null" json:"operation"`
	Resource      string    `gorm:"size:100;not null" json:"resource"`
	ResourceID    uint      `json:"resourceId"`
	IP            string    `gorm:"size:50;not null" json:"ip"`
	UserAgent     string    `gorm:"size:255" json:"userAgent"`
	Status        string    `gorm:"size:20;not null" json:"status"`
	ErrorMessage  string    `gorm:"size:255" json:"errorMessage,omitempty"`
	CreatedAt     time.Time `json:"createTime"`
}
