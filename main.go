package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ft-backend/config"
	"ft-backend/database"
	"ft-backend/routes"
	"ft-backend/utils"
)

func main() {
	// 初始化全局WebSocket管理器
	utils.GlobalWebSocketManager = utils.NewWebSocketManager()
	go utils.GlobalWebSocketManager.Start()

	// 启动机器状态监控器
	go utils.StartMachineStatusMonitor()

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 数据库迁移
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建上传目录
	if err := os.MkdirAll(cfg.File.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// 设置路由
	router := routes.SetupRouter(cfg)

	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
