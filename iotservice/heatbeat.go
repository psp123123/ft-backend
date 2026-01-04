package iotservice

import (
	"ft-backend/common/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 客户端心跳信息顶层结构体（对应整个JSON对象）
type ClientHeartbeat struct {
	ClientID       string     `json:"client_id"`
	HeartbeatTime  int64      `json:"heartbeat_time"` // 毫秒级时间戳，用int64避免溢出
	ClientVersion  string     `json:"client_version"`
	PID            int        `json:"process_id"` // 进程ID，对应JSON中的process_id
	Status         string     `json:"status"`
	LocalIP        string     `json:"local_ip"`
	BusinessModule string     `json:"business_module"`
	TaskCount      int        `json:"task_count"`
	TaskLeft       int        `json:"task_left"`
	LastTaskTime   int64      `json:"last_task_time"`  // 毫秒级时间戳，int64类型
	PrimaryHost    HostInfo   `json:"primary_host"`    // 主主机信息，嵌套HostInfo结构体
	SecondaryHosts []HostInfo `json:"secondary_hosts"` // 副主机列表，切片对应JSON数组
}

// 主机信息结构体（主/副主机共用，字段完全一致）
type HostInfo struct {
	IP               string  `json:"ip"`
	Hostname         string  `json:"hostname"`
	OSInfo           string  `json:"os_info"`
	CPUUsage         float64 `json:"cpu_usage"`     // CPU使用率，浮点型对应JSON中的小数
	MemoryUsage      int64   `json:"memory_usage"`  // 内存占用（字节），int64避免大数值溢出
	DiskUsage        string  `json:"disk_usage"`    // 磁盘可用空间，字符串类型（带单位）
	NetworkDelay     int     `json:"network_delay"` // 网络延迟（毫秒）
	NetworkInterface string  `json:"network_interface"`
	Status           string  `json:"status"` // 主机状态（up/down等）
}

func HeatbeatCheck(c *gin.Context) {
	var clientInfo ClientHeartbeat
	if err := c.ShouldBindJSON(&clientInfo); err != nil {
		// 返回错误信息：
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Error("get client info error:%v", err.Error())
		return
	}
	logger.Debug("get client info:%v", clientInfo)
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
