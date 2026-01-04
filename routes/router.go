package routes

import (
	"ft-backend/common/config"
	"ft-backend/handlers"
	"ft-backend/iotservice"
	"ft-backend/middleware"

	"github.com/gin-gonic/gin"
) // SetupRouter 设置路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin路由
	r := gin.Default()

	// 添加CORS中间件
	r.Use(middleware.CORS())

	// 添加配置信息到上下文
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	})

	// 健康检查
	r.GET("/health", handlers.HealthCheck)

	// 公开路由组 - API路径与前端保持一致
	public := r.Group("/api")
	{
		// 用户认证
		public.POST("/auth/login", handlers.Login)
		public.POST("/auth/logout", handlers.Logout)

		// 文件下载（公开访问）
		public.GET("/files/download/:file_id", handlers.DownloadFile)

		// 调试接口 - 仅用于开发环境
		public.GET("/debug/token", handlers.DebugGetToken)

		// client接口相关
		// // 客户端心跳
		public.POST("/v1/heartbeats", iotservice.HeatbeatCheck)
	}

	// 受保护路由组
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth(cfg.JWT.SecretKey))
	{
		// 仪表盘数据
		protected.GET("/dashboard/data", handlers.GetDashboardData)

		// 用户管理
		protected.GET("/auth/info", handlers.GetUserProfile)
		protected.PUT("/users/profile", handlers.UpdateUserProfile)
		protected.GET("/user", handlers.GetUserList)
		protected.GET("/user/:id", handlers.GetUserDetail)
		protected.POST("/user", handlers.AddUser)
		protected.PUT("/user/:id", handlers.UpdateUser)
		protected.DELETE("/user/:id", handlers.DeleteUser)
		protected.DELETE("/user/batch", handlers.BatchDeleteUser)
		protected.PATCH("/user/:id/role", handlers.UpdateUserRole)

		// 机器管理
		protected.GET("/machine", handlers.GetMachineList)
		protected.GET("/machine/:id", handlers.GetMachineDetail)
		protected.POST("/machine", handlers.AddMachine)
		protected.PUT("/machine/:id", handlers.UpdateMachine)
		protected.DELETE("/machine/:id", handlers.DeleteMachine)
		protected.DELETE("/machine/batch", handlers.BatchDeleteMachine)
		protected.PATCH("/machine/:id/status", handlers.UpdateMachineStatus)

		// 文件管理
		protected.POST("/files/upload", handlers.UploadFile)
		protected.GET("/files/list", handlers.ListFiles)
		protected.GET("/files/:file_id", handlers.GetFileInfo)
		protected.DELETE("/files/:file_id", handlers.DeleteFile)

		// 文件分享
		protected.POST("/files/share/:file_id", handlers.ShareFile)
		protected.GET("/files/shared", handlers.GetSharedFiles)

		// 传输记录
		protected.GET("/transfers", handlers.GetTransferHistory)

		// 安全与审计
		// 操作日志
		protected.GET("/security-audit/operation-logs", handlers.GetOperationLogs)
		protected.GET("/security-audit/operation-logs/:id", handlers.GetOperationLogDetail)

		// 权限管理
		protected.GET("/security-audit/permissions", handlers.GetPermissions)
		protected.GET("/security-audit/permissions/:id", handlers.GetPermissionDetail)
		protected.POST("/security-audit/permissions", handlers.AddPermission)
		protected.PUT("/security-audit/permissions/:id", handlers.UpdatePermission)
		protected.DELETE("/security-audit/permissions/:id", handlers.DeletePermission)
		protected.DELETE("/security-audit/permissions/batch", handlers.BatchDeletePermissions)

		// 角色权限
		protected.GET("/security-audit/roles/:role/permissions", handlers.GetRolePermissions)
		protected.POST("/security-audit/roles/:role/permissions", handlers.AssignRolePermissions)

		// 高级功能
		// 备份与恢复
		protected.GET("/advanced/backup", handlers.GetBackupList)
		protected.GET("/advanced/backup/:id", handlers.GetBackupDetail)
		protected.POST("/advanced/backup", handlers.Backup)
		protected.POST("/advanced/restore/:id", handlers.Restore)

		// 性能分析
		protected.GET("/advanced/performance", handlers.GetPerformanceData)
		protected.POST("/advanced/performance/report", handlers.GeneratePerformanceReport)

		// 调试接口 - 仅用于开发环境
		protected.GET("/debug/test-auth", handlers.DebugTestAuth)

		// K8s部署相关接口
		protected.GET("/k8s/deploy/versions", handlers.GetK8sVersions)
		protected.GET("/k8s/deploy/machines", handlers.GetK8sDeployMachines)
		protected.GET("/k8s/deploy/check-name", handlers.CheckClusterName)

	}

	// WebSocket路由
	r.GET("/ws/:user_id", handlers.WebSocketHandler)

	// 静态文件服务
	r.Static("/uploads", cfg.File.UploadDir)

	return r
}
