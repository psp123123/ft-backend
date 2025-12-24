package handlers

import (
	"net/http"
	"strconv"

	"ft-backend/database"
	"ft-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetOperationLogs 获取操作日志列表
func GetOperationLogs(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	username := c.Query("username")
	operation := c.Query("operation")
	resource := c.Query("resource")
	status := c.Query("status")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	db := database.DB.Model(&models.OperationLog{})

	// 添加过滤条件
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if operation != "" {
		db = db.Where("operation LIKE ?", "%"+operation+"%")
	}
	if resource != "" {
		db = db.Where("resource LIKE ?", "%"+resource+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if startDate != "" {
		db = db.Where("created_at >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("created_at <= ?", endDate)
	}

	// 获取总数
	var total int64
	db.Count(&total)

	// 获取数据
	var logs []models.OperationLog
	db.Limit(pageSize).Offset(offset).Order("created_at DESC").Find(&logs)

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  logs,
			"total": total,
		},
		"msg": "success",
	})
}

// GetOperationLogDetail 获取操作日志详情
func GetOperationLogDetail(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的日志ID",
		})
		return
	}

	// 查询日志
	var log models.OperationLog
	result := database.DB.First(&log, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "日志不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询日志失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": log,
		"msg":  "success",
	})
}

// GetPermissions 获取权限列表
func GetPermissions(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	code := c.Query("code")

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	db := database.DB.Model(&models.Permission{})

	// 添加过滤条件
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		db = db.Where("code LIKE ?", "%"+code+"%")
	}

	// 获取总数
	var total int64
	db.Count(&total)

	// 获取数据
	var permissions []models.Permission
	db.Limit(pageSize).Offset(offset).Order("created_at DESC").Find(&permissions)

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  permissions,
			"total": total,
		},
		"msg": "success",
	})
}

// GetPermissionDetail 获取权限详情
func GetPermissionDetail(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的权限ID",
		})
		return
	}

	// 查询权限
	var permission models.Permission
	result := database.DB.First(&permission, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "权限不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询权限失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": permission,
		"msg":  "success",
	})
}

// AddPermission 添加权限
func AddPermission(c *gin.Context) {
	// 解析请求体
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 检查权限代码是否已存在
	var existingPermission models.Permission
	result := database.DB.Where("code = ?", permission.Code).First(&existingPermission)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "权限代码已存在",
		})
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "检查权限失败",
		})
		return
	}

	// 保存权限
	result = database.DB.Create(&permission)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加权限失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": permission,
		"msg":  "success",
	})
}

// UpdatePermission 更新权限
func UpdatePermission(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的权限ID",
		})
		return
	}

	// 解析请求体
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 查询权限是否存在
	var existingPermission models.Permission
	result := database.DB.First(&existingPermission, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "权限不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询权限失败",
		})
		return
	}

	// 检查权限代码是否已存在（排除当前权限）
	if permission.Code != existingPermission.Code {
		result = database.DB.Where("code = ? AND id != ?", permission.Code, id).First(&models.Permission{})
		if result.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "权限代码已存在",
			})
			return
		} else if result.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "检查权限失败",
			})
			return
		}
	}

	// 更新权限
	permission.ID = uint(id)
	result = database.DB.Save(&permission)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新权限失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": permission,
		"msg":  "success",
	})
}

// DeletePermission 删除权限
func DeletePermission(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的权限ID",
		})
		return
	}

	// 查询权限是否存在
	var permission models.Permission
	result := database.DB.First(&permission, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "权限不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询权限失败",
		})
		return
	}

	// 删除权限
	result = database.DB.Delete(&permission)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除权限失败",
		})
		return
	}

	// 删除关联的角色权限
	database.DB.Where("permission_id = ?", id).Delete(&models.RolePermission{})

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

// BatchDeletePermissions 批量删除权限
func BatchDeletePermissions(c *gin.Context) {
	// 解析请求体
	var request struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 批量删除权限
	result := database.DB.Delete(&models.Permission{}, request.IDs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "批量删除权限失败",
		})
		return
	}

	// 删除关联的角色权限
	database.DB.Where("permission_id IN ?", request.IDs).Delete(&models.RolePermission{})

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

// GetRolePermissions 获取角色权限列表
func GetRolePermissions(c *gin.Context) {
	// 解析角色参数
	role := c.Param("role")

	// 查询角色权限
	var rolePermissions []models.RolePermission
	result := database.DB.Where("role_id = ?", role).Find(&rolePermissions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询角色权限失败",
		})
		return
	}

	// 提取权限ID
	permissionIDs := make([]uint, 0, len(rolePermissions))
	for _, rp := range rolePermissions {
		permissionIDs = append(permissionIDs, rp.PermissionID)
	}

	// 查询权限详情
	var permissions []models.Permission
	if len(permissionIDs) > 0 {
		database.DB.Where("id IN ?", permissionIDs).Find(&permissions)
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  permissions,
			"total": len(permissions),
		},
		"msg": "success",
	})
}

// AssignRolePermissions 分配角色权限
func AssignRolePermissions(c *gin.Context) {
	// 解析角色参数
	role := c.Param("role")

	// 解析请求体
	var request struct {
		PermissionIds []uint `json:"permissionIds"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 开始事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除现有角色权限
	if err := tx.Where("role_id = ?", role).Delete(&models.RolePermission{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除现有角色权限失败",
		})
		return
	}

	// 添加新的角色权限
	for _, permissionId := range request.PermissionIds {
		rolePermission := models.RolePermission{
			RoleID:       role,
			PermissionID: permissionId,
		}
		if err := tx.Create(&rolePermission).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "分配角色权限失败",
			})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "提交事务失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
