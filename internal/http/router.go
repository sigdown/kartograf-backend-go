package http

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sigdown/kartograf-backend-go/internal/auth"
	"github.com/sigdown/kartograf-backend-go/internal/usecase"
)

type Services struct {
	Auth   *usecase.AuthService
	Points *usecase.PointService
	Maps   *usecase.MapService
	Tokens *auth.TokenManager
}

func NewRouter(services Services) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	h := newHandler(services)

	r.POST("/auth/register", h.registerUser)
	r.POST("/auth/login", h.loginUser)
	r.GET("/maps", h.listMaps)
	r.GET("/maps/:slug", h.getMap)

	authGroup := r.Group("/")
	authGroup.Use(AuthRequired(services.Tokens))
	authGroup.GET("/maps/by-id/:id/download", h.downloadMap)
	authGroup.GET("/points", h.listPoints)
	authGroup.POST("/points", h.createPoint)
	authGroup.PATCH("/points/:id", h.updatePoint)
	authGroup.DELETE("/points/:id", h.deletePoint)
	authGroup.PATCH("/account", h.updateAccount)
	authGroup.DELETE("/account", h.deleteAccount)

	adminGroup := authGroup.Group("/admin")
	adminGroup.Use(AdminOnly())
	adminGroup.POST("/maps", h.createMap)
	adminGroup.PATCH("/maps/:id", h.updateMapMetadata)
	adminGroup.PUT("/maps/:id/archive", h.replaceMapArchive)
	adminGroup.DELETE("/maps/:id", h.deleteMap)

	return r
}
