package middleware

import (
	"net/http"
	"strings"
	"time"

	"mindhelp-backend/internal/config"
	"mindhelp-backend/internal/vo"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 聲明結構
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// AuthMiddleware JWT 認證中間件
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Authorization header 獲取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
				"unauthorized",
				"Authorization header is required",
				"UNAUTHORIZED",
				nil,
				c.Request.URL.Path,
			))
			c.Abort()
			return
		}

		// 檢查 Bearer 前綴
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
				"unauthorized",
				"Invalid authorization header format",
				"UNAUTHORIZED",
				nil,
				c.Request.URL.Path,
			))
			c.Abort()
			return
		}

		// 解析 JWT token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// 檢查簽名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
				"unauthorized",
				"Invalid or expired token",
				"UNAUTHORIZED",
				nil,
				c.Request.URL.Path,
			))
			c.Abort()
			return
		}

		// 驗證 token
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// 檢查 token 是否過期
			if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
					"unauthorized",
					"Token has expired",
					"UNAUTHORIZED",
					nil,
					c.Request.URL.Path,
				))
				c.Abort()
				return
			}

			// 將使用者資訊存入 context
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, vo.NewErrorResponse(
				"unauthorized",
				"Invalid token claims",
				"UNAUTHORIZED",
				nil,
				c.Request.URL.Path,
			))
			c.Abort()
			return
		}
	}
}

// GenerateToken 生成 JWT token
func GenerateToken(cfg *config.Config, userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.Expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// GenerateRefreshToken 生成刷新 token
func GenerateRefreshToken(cfg *config.Config, userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// GetUserID 從 context 獲取使用者 ID
func GetUserID(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return ""
	}
	return userID.(string)
}

// GetEmail 從 context 獲取使用者 email
func GetEmail(c *gin.Context) string {
	email, exists := c.Get("email")
	if !exists {
		return ""
	}
	return email.(string)
}
