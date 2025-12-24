package models

import "time"

// K8sCluster 代表Kubernetes集群信息
type K8sCluster struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ClusterName string    `gorm:"unique;not null" json:"cluster_name"`
	Status      string    `gorm:"default:'pending'" json:"status"`
	Version     string    `json:"version"`
	MasterNode  string    `json:"master_node"`
	WorkerNodes string    `json:"worker_nodes"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 设置表名
func (K8sCluster) TableName() string {
	return "k8s_clusters"
}
