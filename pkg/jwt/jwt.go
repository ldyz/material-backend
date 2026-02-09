package jwtpkg

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getJWTSecret())

var (
	ErrTokenExpired = errors.New("token expired")
)

func getJWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "replace_with_a_random_jwt_secret_key"
}

// GenerateToken creates a JWT token containing user id and username
func GenerateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken parses and validates a token string and returns claims or an error
func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	parser := jwt.NewParser()
	t, err := parser.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		// Map jwt errors to clearer sentinel errors
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	if mc, ok := t.Claims.(jwt.MapClaims); ok {
		return mc, nil
	}
	return nil, errors.New("invalid claims")
}

func extractBearer(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("bad auth header")
	}
	return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer ")), nil
}

// TokenMiddleware validates token and aborts on invalid token
func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		var err error

		// Try Authorization header first
		authHeader := c.GetHeader("Authorization")
		tokenStr, err = extractBearer(authHeader)

		// If no valid token in header, try query parameter (for WebSocket connections)
		if err != nil {
			tokenStr = c.Query("token")
			if tokenStr == "" {
				c.AbortWithStatusJSON(401, gin.H{"error": "Token格式错误"})
				return
			}
		}
		claims, err := ParseToken(tokenStr)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				c.AbortWithStatusJSON(401, gin.H{"error": "Token已过期"})
				return
			}
			c.AbortWithStatusJSON(401, gin.H{"error": "Token无效"})
			return
		}
		if uid, ok := claims["user_id"].(float64); ok {
			c.Set("current_user_id", int64(uid))
		}
		if un, ok := claims["username"].(string); ok {
			c.Set("current_username", un)
		}
		c.Set("token_payload", claims)
		c.Next()
	}
}

// TokenOnlyMiddleware validates token but does not enforce DB checks; useful for operations that only require a valid token
func TokenOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		var err error

		// Try Authorization header first
		authHeader := c.GetHeader("Authorization")
		tokenStr, err = extractBearer(authHeader)

		// If no valid token in header, try query parameter (for WebSocket connections)
		if err != nil {
			tokenStr = c.Query("token")
			if tokenStr == "" {
				c.AbortWithStatusJSON(401, gin.H{"error": "Token格式错误"})
				return
			}
		}
		claims, err := ParseToken(tokenStr)
		if err != nil {
			if errors.Is(err, ErrTokenExpired) {
				c.AbortWithStatusJSON(401, gin.H{"error": "Token已过期"})
				return
			}
			c.AbortWithStatusJSON(401, gin.H{"error": "Token无效"})
			return
		}
		// set claims into context for handlers that will fetch DB later
		if uid, ok := claims["user_id"].(float64); ok {
			c.Set("current_user_id", int64(uid))
		}
		if un, ok := claims["username"].(string); ok {
			c.Set("current_username", un)
		}
		c.Set("token_payload", claims)
		c.Next()
	}
}
