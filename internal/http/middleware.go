package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	apiauth "github.com/sigdown/kartograf-backend-go/internal/auth"
	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

const claimsContextKey = "auth_claims"

func AuthRequired(tokens *apiauth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			writeError(c, http.StatusUnauthorized, "missing bearer token")
			c.Abort()
			return
		}

		claims, err := tokens.Parse(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			writeError(c, statusFromError(err), err.Error())
			c.Abort()
			return
		}

		c.Set(claimsContextKey, claims)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := CurrentClaims(c)
		if !ok {
			writeError(c, http.StatusUnauthorized, "missing auth context")
			c.Abort()
			return
		}

		if claims.Role != domain.RoleAdmin {
			writeError(c, http.StatusForbidden, "admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

func CurrentClaims(c *gin.Context) (apiauth.Claims, bool) {
	raw, ok := c.Get(claimsContextKey)
	if !ok {
		return apiauth.Claims{}, false
	}

	claims, ok := raw.(apiauth.Claims)
	return claims, ok
}

func statusFromError(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func writeError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": message,
	})
}
