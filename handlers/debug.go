package handlers

import (
	"net/http"

	"ft-backend/common/config"
	"ft-backend/common/logger"
	"ft-backend/utils"

	"github.com/gin-gonic/gin"
)

// DebugGetToken 调试接口：获取测试token
// 注意：此接口仅用于调试，生产环境应该移除
func DebugGetToken(c *gin.Context) {
	logger.Debug("DebugGetToken called")

	// 获取配置
	cfg := c.MustGet("config").(*config.Config)
	logger.Debug("JWT Secret Key: %s", cfg.JWT.SecretKey)

	// 生成测试token
	token, err := utils.GenerateAccessToken(
		1,                   // 用户ID
		"admin",             // 用户名
		"admin@example.com", // 邮箱
		"admin",             // 角色
		cfg.JWT.SecretKey,
		60, // 1小时过期
	)

	if err != nil {
		logger.Error("Error generating token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成token失败",
		})
		return
	}

	logger.Debug("Generated token: %s", token)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取调试token成功",
		"data": gin.H{
			"token": token,
			"usage": "在Authorization头中使用: Bearer " + token,
		},
	})
}

// DebugTestAuth 调试接口：测试JWT认证
func DebugTestAuth(c *gin.Context) {
	logger.Debug("DebugTestAuth called")

	// 从上下文获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		logger.Debug("Error: userID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户未认证",
		})
		return
	}

	username, _ := c.Get("username")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	logger.Debug("Auth successful. UserID: %v, Username: %v, Email: %v, Role: %v",
		userID, username, email, role)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "认证成功",
		"data": gin.H{
			"userID":   userID,
			"username": username,
			"email":    email,
			"role":     role,
		},
	})
}