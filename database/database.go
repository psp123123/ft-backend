package database

import (
	"fmt"
	"log"

	"ft-backend/config"
	"ft-backend/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// Migrate 数据库迁移
func Migrate() error {
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
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 初始化K8s版本数据
	initK8sVersions()

	log.Println("Database migration completed successfully")
	return nil
}

// initK8sVersions 初始化K8s版本数据
func initK8sVersions() {
	// 检查数据库中是否已有K8s版本记录
	var count int64
	DB.Model(&models.K8sVersion{}).Count(&count)
	
	// 如果数据库中没有版本记录，则初始化默认版本
	if count == 0 {
		// 初始化默认的Kubernetes版本
		defaultVersions := []string{"v1.35.0", "v1.32.11", "v1.34.3", "v1.32.6", "v1.28.15"}
		
		// 批量插入默认版本
		for _, version := range defaultVersions {
			DB.Create(&models.K8sVersion{Version: version, IsActive: true})
		}
	}
}

// GetK8sVersions 从数据库中获取所有有效的Kubernetes版本
func GetK8sVersions() ([]models.K8sVersion, error) {
	var versions []models.K8sVersion
	
	// 查询所有激活的K8s版本
	result := DB.Where("is_active = ?", true).Find(&versions)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get k8s versions: %w", result.Error)
	}
	
	return versions, nil
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("Database connection closed")
	return nil
}
