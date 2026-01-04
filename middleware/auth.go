package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"ft-backend/common/logger"
	"ft-backend/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debug("Processing request: %s %s", c.Request.Method, c.Request.URL.Path)

		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		logger.Debug("Authorization header: %s", authHeader)

		if authHeader == "" {
			logger.Debug("Error: Authorization header is empty")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Authorization header is required",
			})
			c.Abort()
			return
		}

		// 检查token格式
		parts := strings.SplitN(authHeader, " ", 2)
		logger.Debug("Token parts: %v", parts)

		if !(len(parts) == 2 && parts[0] == "Bearer") {
			logger.Debug("Error: Invalid token format. Expected 'Bearer {token}'")
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Authorization header format must be Bearer {token}",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		logger.Debug("Token string: %s", tokenString)

		// 使用专门的验证函数来解析token
		logger.Debug("Validating token with secret: %s", secretKey)
		claims, err := utils.ValidateToken(tokenString, secretKey)
		if err != nil {
			logger.Debug("Error validating token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  fmt.Sprintf("Invalid or expired token: %v", err),
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文，使用驼峰式命名
		logger.Debug("Token valid. UserID: %d, Username: %s, Role: %s", claims.UserID, claims.Username, claims.Role)
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}