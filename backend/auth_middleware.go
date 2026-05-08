package main

import (
	"net"
	"net/http"
	"net/url"
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

func (m *AuthMiddleware) matchesPublicHosts(ip net.IP, hostname string) bool {
	for _, host := range m.cfg.PublicHosts {
		host = strings.TrimSpace(host)
		if strings.Contains(host, "/") {
			_, network, err := net.ParseCIDR(host)
			if err == nil && ip != nil && network.Contains(ip) {
				return true
			}
		} else {
			if ip != nil && ip.String() == host {
				return true
			}
			if hostname == host {
				return true
			}
		}
	}
	return false
}

func (m *AuthMiddleware) isPublicHost(c *gin.Context) bool {
	ip := net.ParseIP(c.ClientIP())
	hostname := strings.Split(c.Request.Host, ":")[0]
	return m.matchesPublicHosts(ip, hostname)
}

func (m *AuthMiddleware) isPublicOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	ip := net.ParseIP(u.Hostname())
	return m.matchesPublicHosts(ip, u.Hostname())
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
