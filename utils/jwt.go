package utils

import (
	"errors"
	"time"

	"ft-backend/common/logger"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT claims结构
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken 生成访问令牌
func GenerateAccessToken(userID uint, username, email, role, secretKey string, expiresIn int) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint, username, secretKey string, expiresIn int) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   username,
		ID:        string(rune(userID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ValidateToken 验证令牌
func ValidateToken(tokenString, secretKey string) (*JWTClaims, error) {
	logger.Debug("正在验证JWT令牌")
	logger.Debug("使用的密钥: %s", secretKey)

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		logger.Debug("令牌签名方法: %v", token.Method)
		return []byte(secretKey), nil
	})

	if err != nil {
		logger.Error("令牌解析错误: %v", err)
		return nil, err
	}

	logger.Debug("令牌解析成功, 有效性: %t", token.Valid)

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		logger.Debug("JWT声明信息: %+v", claims)
		return claims, nil
	}

	logger.Warn("无效的令牌声明或令牌无效")
	return nil, errors.New("invalid token")
}

// ExtractUserIDFromToken 从令牌中提取用户ID
func ExtractUserIDFromToken(tokenString, secretKey string) (uint, error) {
	claims, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}