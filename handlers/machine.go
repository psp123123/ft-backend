package handlers

import (
	"net/http"
	"strconv"

	"ft-backend/database"
	"ft-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetMachineList 获取机器列表
func GetMachineList(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	name := c.Query("name")
	status := c.Query("status")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	db := database.DB.Model(&models.Machine{})

	// 添加过滤条件
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
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
	var machines []models.Machine
	db.Limit(pageSize).Offset(offset).Order("created_at DESC").Find(&machines)

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  machines,
			"total": total,
		},
		"msg": "success",
	})
}

// GetMachineDetail 获取机器详情
func GetMachineDetail(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的机器ID",
		})
		return
	}

	// 查询机器
	var machine models.Machine
	result := database.DB.First(&machine, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "机器不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询机器失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": machine,
		"msg":  "success",
	})
}

// AddMachine 添加机器
func AddMachine(c *gin.Context) {
	// 解析请求体
	var machine models.Machine
	if err := c.ShouldBindJSON(&machine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 保存机器
	result := database.DB.Create(&machine)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加机器失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": machine,
		"msg":  "success",
	})
}

// UpdateMachine 更新机器
func UpdateMachine(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的机器ID",
		})
		return
	}

	// 解析请求体
	var machine models.Machine
	if err := c.ShouldBindJSON(&machine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 查询机器是否存在
	var existingMachine models.Machine
	result := database.DB.First(&existingMachine, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "机器不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询机器失败",
		})
		return
	}

	// 更新机器
	machine.ID = uint(id)
	result = database.DB.Save(&machine)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新机器失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": machine,
		"msg":  "success",
	})
}

// DeleteMachine 删除机器
func DeleteMachine(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的机器ID",
		})
		return
	}

	// 查询机器是否存在
	var machine models.Machine
	result := database.DB.First(&machine, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "机器不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询机器失败",
		})
		return
	}

	// 删除机器（软删除）
	result = database.DB.Delete(&machine)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除机器失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

// BatchDeleteMachine 批量删除机器
func BatchDeleteMachine(c *gin.Context) {
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

	// 批量删除机器（软删除）
	result := database.DB.Delete(&models.Machine{}, request.IDs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "批量删除机器失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

// UpdateMachineStatus 更新机器状态
func UpdateMachineStatus(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的机器ID",
		})
		return
	}

	// 解析请求体
	var request struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 查询机器是否存在
	var machine models.Machine
	result := database.DB.First(&machine, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "机器不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询机器失败",
		})
		return
	}

	// 更新机器状态
	machine.Status = request.Status
	result = database.DB.Save(&machine)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新机器状态失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": machine,
		"msg":  "success",
	})
}
