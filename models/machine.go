package models

import (
	"time"
)

type Machine struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	IP        string    `gorm:"size:50;not null" json:"ip"`
	CPU       int       `gorm:"not null" json:"cpu"`
	Memory    int       `gorm:"not null" json:"memory"`
	Disk      int       `gorm:"not null" json:"disk"`
	Status    string    `gorm:"size:20;default:'offline'" json:"status"`
	CreatedAt time.Time `json:"createTime"`
	UpdatedAt time.Time `json:"updateTime"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}
