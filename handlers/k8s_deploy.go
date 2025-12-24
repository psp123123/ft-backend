package handlers

import (
	"net/http"

	"ft-backend/database"
	"ft-backend/models"

	"github.com/gin-gonic/gin"
)

// GetK8sVersions 获取K8s版本列表
func GetK8sVersions(c *gin.Context) {
	// 从数据库中获取所有有效的Kubernetes版本
	versions, err := database.GetK8sVersions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to fetch K8s versions",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": versions,
	})
}

// GetK8sDeployMachines 获取可用于K8s部署的机器列表
func GetK8sDeployMachines(c *gin.Context) {
	var machines []models.Machine
	status := c.Query("status")

	// 构建查询条件
	query := database.DB
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 查询符合条件的机器
	if err := query.Find(&machines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to fetch machines",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": machines,
	})
}

// CheckClusterName 检查集群名称是否可用
func CheckClusterName(c *gin.Context) {
	clusterName := c.Query("params[clusterName]")

	if clusterName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Cluster name is required",
			"data": gin.H{
				"isAvailable": false,
			},
		})
		return
	}

	var count int64
	if err := database.DB.Model(&models.K8sCluster{}).Where("cluster_name = ?", clusterName).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to check cluster name",
			"data": gin.H{
				"isAvailable": false,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"isAvailable": count == 0,
		},
	})
}
