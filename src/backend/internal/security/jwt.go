package security

import (
	"errors"
	"fmt"
	"time"

	"biolitmanager/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// TokenExpireTime Token 过期时间（2小时）
	TokenExpireTime = time.Hour * 2
)

// Claims JWT 载荷
type Claims struct {
	UserID      uint     `json:"user_id"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username string, role string, permissions []string) (string, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return "", errors.New("config not initialized")
	}

	secret := cfg.JWT.Secret
	if secret == "" {
		return "", errors.New("jwt secret is empty")
	}

	nowTime := time.Now()
	expireTime := nowTime.Add(TokenExpireTime)

	claims := Claims{
		UserID:      userID,
		Username:    username,
		Role:        role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "biolitmanager",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()
	if cfg == nil {
		return nil, errors.New("config not initialized")
	}

	secret := cfg.JWT.Secret
	if secret == "" {
		return nil, errors.New("jwt secret is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
