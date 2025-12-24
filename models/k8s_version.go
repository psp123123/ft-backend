package models

import "time"

// K8sVersion 代表Kubernetes版本信息
type K8sVersion struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Version   string    `gorm:"unique;not null" json:"version"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 设置表名
func (K8sVersion) TableName() string {
	return "k8s_versions"
}
