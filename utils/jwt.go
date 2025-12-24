package utils

import (
	"errors"
	"log"
	"time"

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
	log.Printf("[JWT Utils] Validating token: %s", tokenString)
	log.Printf("[JWT Utils] Using secret key: %s", secretKey)

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		log.Printf("[JWT Utils] Token method: %v", token.Method)
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("[JWT Utils] Error parsing token: %v", err)
		return nil, err
	}

	log.Printf("[JWT Utils] Token parsed successfully. Valid: %t", token.Valid)

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		log.Printf("[JWT Utils] Claims: %+v", claims)
		return claims, nil
	}

	log.Printf("[JWT Utils] Invalid token claims or token not valid")
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
