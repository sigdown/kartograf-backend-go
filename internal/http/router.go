package http

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	r := gin.Default()

	RegisterHealthRoutes(r)

	return r
}