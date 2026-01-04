package main

import (
	"fmt"

	"ft-backend/common/config"
	"ft-backend/common/logger"
	"ft-backend/database"
	"ft-backend/routes"
	"ft-backend/utils"
	"net/http"
	"os"
)

func main() {
	const configPath = "conf/config.yaml"
	// 1️. 检查或创建配置文件
	if err := config.EnsureConfigExists(configPath); err != nil {
		logger.Error("Failed to ensure config exists: %v", err)
		return
	}

	// 2️. 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config: %v", err)
		return
	}

	config.GlobalCfg = cfg
	logger.InitLogger(config.GlobalCfg.Log.Level, nil)
	logger.Info("Loaded config: %+v", config.GlobalCfg)

	// 初始化全局WebSocket管理器
	utils.GlobalWebSocketManager = utils.NewWebSocketManager()
	go utils.GlobalWebSocketManager.Start()

	// 启动机器状态监控器
	go utils.StartMachineStatusMonitor()

	// 连接数据库
	if err := database.Connect(&cfg.Database); err != nil {
		logger.Error("Failed to connect to database: %v", err)
		return
	}
	defer database.Close()

	// 数据库迁移
	if err := database.Migrate(); err != nil {
		logger.Error("Failed to migrate database: %v", err)
		return
	}

	// 创建上传目录
	if err := os.MkdirAll(cfg.File.UploadDir, 0755); err != nil {
		logger.Error("Failed to create upload directory: %v", err)
		return
	}

	// 设置路由
	router := routes.SetupRouter(cfg)

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Server starting on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		logger.Error("Failed to start server: %v", err)
	}
}
