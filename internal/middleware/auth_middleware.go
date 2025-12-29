package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		tokenString := parts[0]
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			tokenString = parts[1]
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil {
			fmt.Printf("DEBUG [Auth]: JWT Parse Error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			uid, okUID := claims["user_id"].(float64)
			role, okRole := claims["role"].(string)

			if !okUID || !okRole {
				fmt.Println("DEBUG [Auth]: Token missing required claims (user_id or role)")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				return
			}

			c.Set("user_id", int(uid))
			c.Set("role", role)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		}
	}
}

func RoleBlock(authorizedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Identity not found"})
			return
		}

		roleStr := userRole.(string)
		
		if roleStr == "admin" {
			c.Next()
			return
		}

		for _, role := range authorizedRoles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": fmt.Sprintf("Access denied. Role '%s' is not authorized for this action", roleStr),
		})
	}
}