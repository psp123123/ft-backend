package utils

import (
	"time"

	"ft-backend/common/logger"
	"ft-backend/database"
	"ft-backend/models"
)

// StartMachineStatusMonitor 启动机器状态监控器
func StartMachineStatusMonitor() {
	// 每5秒检查一次机器状态
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	logger.Info("Machine status monitor started")

	for range ticker.C {
		// 检查机器状态
		checkMachineStatus()
	}

}

// checkMachineStatus 检查所有机器的状态并广播更新
func checkMachineStatus() {
	// 获取所有机器
	var machines []models.Machine
	result := database.DB.Find(&machines)
	if result.Error != nil {
		logger.Error("Failed to get machines: %v", result.Error)
		return
	}

	// 这里可以添加实际的状态检查逻辑
	// 例如，通过ping检查机器是否在线
	// 为了演示目的，我们随机改变一些机器的状态

	// 广播机器状态更新
	statusUpdateMsg := WebSocketMessage{
		Type:    "machine_status_update",
		Message: "Machine status updated",
		Data:    machines,
	}

	GlobalWebSocketManager.Broadcast(statusUpdateMsg)

	logger.Debug("Broadcasted machine status update to %d machines", len(machines))
}