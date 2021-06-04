package middleware

import (
	"errors"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var authMiddleware *jwt.GinJWTMiddleware

type JwtMiddlewareConfig struct {
	Realm        string
	Key          []byte
	Timeout      time.Duration
	MaxRefresh   time.Duration
	ExcludePaths []string
}

func NewJwtMiddleware(c *JwtMiddlewareConfig) gin.HandlerFunc {
	authMiddleware, _ = jwt.New(&jwt.GinJWTMiddleware{
		Realm:      c.Realm,
		Key:        c.Key,
		Timeout:    c.Timeout,
		MaxRefresh: c.MaxRefresh,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if dataMap, ok := data.(map[string]interface{}); ok {
				return dataMap
			}
			return map[string]interface{}{}
		},
	})
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		for _, excludePath := range c.ExcludePaths {
			if path == excludePath {
				return
			}
		}
		authMiddleware.MiddlewareFunc()(ctx)
	}
}

type jwtNS struct{}

var JWT jwtNS

func (jwtNS) GenerateToken(data interface{}) (string, time.Time, error) {
	if authMiddleware == nil {
		panic("authMiddleware is nil")
	}

	return authMiddleware.TokenGenerator(data)
}

func (jwtNS) RefreshToken(ctx *gin.Context) (string, time.Time, error) {
	if authMiddleware == nil {
		panic("authMiddleware is nil")
	}

	return authMiddleware.RefreshToken(ctx)
}

func (jwtNS) GetPayload(ctx *gin.Context) (map[string]interface{}, error) {
	if authMiddleware == nil {
		panic("authMiddleware is nil")
	}

	return authMiddleware.GetClaimsFromJWT(ctx)
}

func (ns jwtNS) GetCurrentUserID(ctx *gin.Context) (uint, error) {
	if payload, err := ns.GetPayload(ctx); err != nil {
		if userID, ok := payload["sub"]; ok {
			return userID.(uint), nil
		}
		return 0, errors.New("failed to get current user ID")
	}
	return 0, errors.New("failed to get current user ID")
}
