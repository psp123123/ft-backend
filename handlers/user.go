package handlers

import (
	"net/http"
	"strconv"

	"ft-backend/database"
	"ft-backend/models"
	"ft-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserProfile 获取当前用户信息
func GetUserProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 查询用户信息
	var user models.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "用户不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询用户失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
		"msg":  "success",
	})
}

// UpdateUserProfile 更新当前用户信息
func UpdateUserProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 解析请求体
	var request struct {
		Phone    string `json:"phone"`
		FullName string `json:"full_name"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 查询用户
	var user models.User
	result := database.DB.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "用户不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询用户失败",
		})
		return
	}

	// 更新用户信息
	if request.Phone != "" {
		user.Phone = request.Phone
	}
	if request.FullName != "" {
		user.FullName = request.FullName
	}
	if request.Avatar != "" {
		user.Avatar = request.Avatar
	}

	// 保存更新
	result = database.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新用户信息失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
		"msg":  "success",
	})
}

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	username := c.Query("username")
	role := c.Query("role")

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	db := database.DB.Model(&models.User{})

	// 添加过滤条件
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if role != "" {
		db = db.Where("role = ?", role)
	}

	// 获取总数
	var total int64
	db.Count(&total)

	// 获取数据
	var users []models.User
	db.Limit(pageSize).Offset(offset).Order("created_at DESC").Find(&users)

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  users,
			"total": total,
		},
		"msg": "success",
	})
}

// GetUserDetail 获取用户详情
func GetUserDetail(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}

	// 查询用户
	var user models.User
	result := database.DB.First(&user, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "用户不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询用户失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
		"msg":  "success",
	})
}

// AddUser 添加用户
func AddUser(c *gin.Context) {
	// 解析请求体
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	result := database.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "检查用户名失败",
		})
		return
	}

	// 检查邮箱是否已存在
	result = database.DB.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "邮箱已存在",
		})
		return
	} else if result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "检查邮箱失败",
		})
		return
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}
	user.Password = hashedPassword

	// 设置默认角色
	if user.Role == "" {
		user.Role = "user"
	}

	// 保存用户
	result = database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加用户失败",
		})
		return
	}

	// 清除密码字段
	user.Password = ""

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
		"msg":  "success",
	})
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}

	// 解析请求体
	var request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Role     string `json:"role"`
		FullName string `json:"full_name"`
		Avatar   string `json:"avatar"`
		Password string `json:"password,omitempty"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 查询用户
	var user models.User
	result := database.DB.First(&user, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "用户不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询用户失败",
		})
		return
	}

	// 更新用户信息
	if request.Username != "" && request.Username != user.Username {
		// 检查用户名是否已存在
		var existingUser models.User
		result = database.DB.Where("username = ?", request.Username).First(&existingUser)
		if result.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "用户名已存在",
			})
			return
		} else if result.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "检查用户名失败",
			})
			return
		}
		user.Username = request.Username
	}

	if request.Email != "" && request.Email != user.Email {
		// 检查邮箱是否已存在
		var existingUser models.User
		result = database.DB.Where("email = ?", request.Email).First(&existingUser)
		if result.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "邮箱已存在",
			})
			return
		} else if result.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "检查邮箱失败",
			})
			return
		}
		user.Email = request.Email
	}

	if request.Phone != "" {
		user.Phone = request.Phone
	}
	if request.Role != "" {
		user.Role = request.Role
	}
	if request.FullName != "" {
		user.FullName = request.FullName
	}
	if request.Avatar != "" {
		user.Avatar = request.Avatar
	}
	if request.Password != "" {
		// 哈希新密码
		hashedPassword, err := utils.HashPassword(request.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "密码加密失败",
			})
			return
		}
		user.Password = hashedPassword
	}

	// 保存更新
	result = database.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新用户失败",
		})
		return
	}

	// 清除密码字段
	user.Password = ""

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
		"msg":  "success",
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}

	// 查询用户是否存在
	var user models.User
	result := database.DB.First(&user, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "用户不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询用户失败",
		})
		return
	}

	// 不能删除管理员用户
	if user.Role == "admin" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "不能删除管理员用户",
		})
		return
	}

	// 删除用户（软删除）
	result = database.DB.Delete(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除用户失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

// BatchDeleteUser 批量删除用户
func BatchDeleteUser(c *gin.Context) {
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

	// 检查是否包含管理员用户
	var adminCount int64
	database.DB.Model(&models.User{}).Where("id IN ? AND role = ?", request.IDs, "admin").Count(&adminCount)
	if adminCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "不能删除管理员用户",
		})
		return
	}

	// 批量删除用户（软删除）
	result := database.DB.Delete(&models.User{}, request.IDs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "批量删除用户失败",
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}

// UpdateUserRole 更新用户角色
func UpdateUserRole(c *gin.Context) {
	// 解析ID参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的用户ID",
		})
		return
	}

	// 解析请求体
	var request struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的请求参数",
		})
		return
	}

	// 查询用户
	var user models.User
	result := database.DB.First(&user, uint(id))
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"code": 404,
				"msg":  "用户不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询用户失败",
		})
		return
	}

	// 更新角色
	user.Role = request.Role

	// 保存更新
	result = database.DB.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新用户角色失败",
		})
		return
	}

	// 清除密码字段
	user.Password = ""

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
		"msg":  "success",
	})
}