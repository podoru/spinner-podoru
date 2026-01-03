package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/podoru/podoru/internal/usecase/auth"
	"github.com/podoru/podoru/pkg/response"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	UserIDKey           = "user_id"
	UserEmailKey        = "user_email"
)

type AuthMiddleware struct {
	authUseCase *auth.UseCase
}

func NewAuthMiddleware(authUseCase *auth.UseCase) *AuthMiddleware {
	return &AuthMiddleware{authUseCase: authUseCase}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)
		if token == "" {
			response.Unauthorized(c, "Token is required")
			c.Abort()
			return
		}

		claims, err := m.authUseCase.ValidateAccessToken(token)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)
		if token == "" {
			c.Next()
			return
		}

		claims, err := m.authUseCase.ValidateAccessToken(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserEmailKey, claims.Email)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return uuid.UUID{}, false
	}

	id, ok := userID.(uuid.UUID)
	return id, ok
}

func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}

	e, ok := email.(string)
	return e, ok
}
