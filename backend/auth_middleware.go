package main

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	cfg *Config
}

func NewAuthMiddleware(cfg *Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

func (m *AuthMiddleware) isPublicHost(c *gin.Context) bool {
	clientIP := net.ParseIP(c.ClientIP())
	for _, host := range m.cfg.PublicHosts {
		host = strings.TrimSpace(host)
		if strings.Contains(host, "/") {
			_, network, err := net.ParseCIDR(host)
			if err == nil && clientIP != nil && network.Contains(clientIP) {
				return true
			}
		} else {
			if clientIP != nil && clientIP.String() == host {
				return true
			}
			if c.Request.Host == host || strings.Split(c.Request.Host, ":")[0] == host {
				return true
			}
		}
	}
	return false
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	if m.isPublicHost(c) {
		c.Next()
		return
	}

	tokenStr := c.Query("token")
	if tokenStr == "" {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr = strings.TrimPrefix(header, "Bearer ")
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return []byte(m.cfg.JWTSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	c.Set("userID", claims["sub"])
	c.Set("email", claims["email"])
	c.Next()
}
