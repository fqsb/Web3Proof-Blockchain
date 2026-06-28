package middleware

import (
	"net/http"
	"strings"

	"web3proof/backend/internal/model"
	"web3proof/backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Claims struct {
	UserID        uint64 `json:"user_id"`
	WalletAddress string `json:"wallet_address"`
	ActiveRole    string `json:"active_role"`
	jwt.RegisteredClaims
}

func CORS(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func JWTAuth(secret string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, http.StatusUnauthorized, 40101, "unauthorized")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			response.Fail(c, http.StatusUnauthorized, 40101, "invalid token")
			c.Abort()
			return
		}
		var user model.User
		if err := db.First(&user, claims.UserID).Error; err != nil || user.Status != 1 {
			response.Fail(c, http.StatusUnauthorized, 40101, "user unavailable")
			c.Abort()
			return
		}
		c.Set("user_id", user.ID)
		c.Set("wallet_address", user.WalletAddress)
		c.Set("active_role", user.ActiveRole)
		c.Next()
	}
}

func OptionalJWTAuth(secret string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.Next()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.Next()
			return
		}
		var user model.User
		if err := db.First(&user, claims.UserID).Error; err != nil || user.Status != 1 {
			c.Next()
			return
		}
		c.Set("user_id", user.ID)
		c.Set("wallet_address", user.WalletAddress)
		c.Set("active_role", user.ActiveRole)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, _ := c.Get("active_role")
		roleStr, _ := role.(string)
		if _, ok := allowed[roleStr]; !ok {
			response.Fail(c, http.StatusForbidden, 40301, "forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
