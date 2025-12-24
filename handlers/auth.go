package handlers

import (
	"net/http"

	"ft-backend/config"
	"ft-backend/database"
	"ft-backend/models"
	"ft-backend/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"omitempty,max=100"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid input", "error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "Username already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Database error"})
		return
	}

	// 检查邮箱是否已存在
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "Email already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Database error"})
		return
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to hash password"})
		return
	}

	// 创建新用户
	newUser := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Role:     "user",
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "User registered successfully",
		"user": gin.H{
			"id":        newUser.ID,
			"username":  newUser.Username,
			"email":     newUser.Email,
			"full_name": newUser.FullName,
			"role":      newUser.Role,
		},
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid input", "error": err.Error()})
		return
	}

	// 获取用户配置
	cfg := c.MustGet("config").(*config.Config)

	// 查找用户
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Database error"})
		}
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid credentials"})
		return
	}

	// 生成访问令牌
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Username,
		user.Email,
		user.Role,
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenExp,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to generate access token"})
		return
	}

	// 返回结果，与前端类型定义匹配
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"token": accessToken,
			"user": gin.H{
				"id":         user.ID,
				"username":   user.Username,
				"email":      user.Email,
				"phone":      user.Phone,
				"role":       user.Role,
				"full_name":  user.FullName,
				"avatar":     user.Avatar,
				"createTime": user.CreatedAt,
				"updateTime": user.UpdatedAt,
			},
		},
		"msg": "success",
	})
}

// RefreshToken 刷新Token
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid input", "error": err.Error()})
		return
	}

	// 获取用户配置
	cfg := c.MustGet("config").(*config.Config)

	// 解析刷新令牌
	claims, err := utils.ValidateToken(req.RefreshToken, cfg.JWT.SecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid refresh token"})
		return
	}

	// 查找用户
	var user models.User
	if err := database.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Database error"})
		}
		return
	}

	// 生成新的访问令牌
	newAccessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.Username,
		user.Email,
		user.Role,
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenExp,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "Token refreshed successfully",
		"data": gin.H{
			"access_token": newAccessToken,
			"token_type":   "Bearer",
			"expires_in":   cfg.JWT.AccessTokenExp * 60,
		},
	})
}

// Logout 用户登出
func Logout(c *gin.Context) {
	// 在这个简单实现中，我们只需要客户端删除token即可
	// 更复杂的实现可以将token添加到黑名单
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": nil,
		"msg":  "success",
	})
}
