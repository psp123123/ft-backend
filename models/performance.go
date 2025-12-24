package models

import (
	"time"
)

type PerformanceData struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	MachineID    uint      `gorm:"not null" json:"machineId"`
	MachineName  string    `gorm:"size:100;not null" json:"machineName"`
	CPUUsage     float64   `gorm:"not null" json:"cpuUsage"`
	MemoryUsage  float64   `gorm:"not null" json:"memoryUsage"`
	DiskUsage    float64   `gorm:"not null" json:"diskUsage"`
	NetworkIn    float64   `json:"networkIn"`
	NetworkOut   float64   `json:"networkOut"`
	Timestamp    time.Time `json:"timestamp"`
	CreatedAt    time.Time `json:"created_at"`
}
