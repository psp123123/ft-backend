package database

import (
	"fmt"
	"time"

	"ft-backend/common/config"
	"ft-backend/common/logger"
	"ft-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var DB *gorm.DB

// Connect 连接数据库
func Connect(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	logger.Info("正在连接数据库: %s@%s:%s", cfg.User, cfg.Host, cfg.Port)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		logger.Error("数据库连接失败: %v", err)
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层 sql.DB 对象以配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Error("获取数据库实例失败: %v", err)
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		logger.Error("数据库连接测试失败: %v", err)
		return fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("数据库连接成功")
	return nil
}

// Migrate 数据库迁移
func Migrate() error {
	logger.Info("开始数据库迁移")

	err := DB.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.Transfer{},
		&models.Share{},
		&models.Machine{},
		&models.OperationLog{},
		&models.Permission{},
		&models.RolePermission{},
		&models.PerformanceData{},
		&models.K8sVersion{},
		&models.K8sCluster{},
	)

	if err != nil {
		logger.Error("数据库迁移失败: %v", err)
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 初始化K8s版本数据
	initK8sVersions()

	logger.Info("数据库迁移完成")
	return nil
}

// initK8sVersions 初始化K8s版本数据
func initK8sVersions() {
	logger.Debug("检查K8s版本数据")

	// 检查数据库中是否已有K8s版本记录
	var count int64
	result := DB.Model(&models.K8sVersion{}).Count(&count)
	if result.Error != nil {
		logger.Error("查询K8s版本数量失败: %v", result.Error)
		return
	}

	// 如果数据库中没有版本记录，则初始化默认版本
	if count == 0 {
		logger.Info("初始化默认K8s版本数据")

		// 初始化默认的Kubernetes版本
		defaultVersions := []string{
			"v1.35.0", "v1.32.11", "v1.34.3",
			"v1.32.6", "v1.28.15", "v1.30.0",
		}

		// 批量插入默认版本
		var versions []models.K8sVersion
		for _, version := range defaultVersions {
			versions = append(versions, models.K8sVersion{
				Version:  version,
				IsActive: true,
			})
		}

		if err := DB.CreateInBatches(versions, 5).Error; err != nil {
			logger.Error("批量插入K8s版本失败: %v", err)
			return
		}

		logger.Info("成功初始化 %d 个K8s版本", len(versions))
	} else {
		logger.Debug("K8s版本数据已存在，跳过初始化")
	}
}

// GetK8sVersions 从数据库中获取所有有效的Kubernetes版本
func GetK8sVersions() ([]models.K8sVersion, error) {
	logger.Debug("获取K8s版本列表")

	var versions []models.K8sVersion

	// 查询所有激活的K8s版本
	result := DB.Where("is_active = ?", true).Find(&versions)
	if result.Error != nil {
		logger.Error("查询K8s版本失败: %v", result.Error)
		return nil, fmt.Errorf("failed to get k8s versions: %w", result.Error)
	}

	logger.Debug("成功获取 %d 个K8s版本", len(versions))
	return versions, nil
}

// GetDBStatus 获取数据库连接状态
func GetDBStatus() map[string]interface{} {
	status := make(map[string]interface{})

	sqlDB, err := DB.DB()
	if err != nil {
		status["status"] = "error"
		status["error"] = fmt.Sprintf("获取数据库实例失败: %v", err)
		return status
	}

	// 获取连接池统计信息
	stats := sqlDB.Stats()
	status["status"] = "connected"
	status["idle"] = stats.Idle
	status["in_use"] = stats.InUse
	status["max_open_connections"] = stats.MaxOpenConnections
	status["open_connections"] = stats.OpenConnections
	status["wait_count"] = stats.WaitCount
	status["wait_duration"] = stats.WaitDuration.String()

	return status
}

// Close 关闭数据库连接
func Close() error {
	logger.Info("正在关闭数据库连接")

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Error("获取数据库实例失败: %v", err)
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		logger.Error("关闭数据库连接失败: %v", err)
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	logger.Info("数据库连接已关闭")
	return nil
}

// GetGormDB 获取GORM数据库实例（用于特殊操作）
func GetGormDB() *gorm.DB {
	return DB
}

// ExecRawSQL 执行原生SQL（仅用于特殊情况）
func ExecRawSQL(query string, args ...interface{}) error {
	logger.Warn("执行原生SQL: %s", utils.ToString(query))

	result := DB.Exec(query, args...)
	if result.Error != nil {
		logger.Error("原生SQL执行失败: %v", result.Error)
		return fmt.Errorf("failed to execute SQL: %w", result.Error)
	}

	logger.Debug("原生SQL执行成功，影响行数: %d", result.RowsAffected)
	return nil
}
