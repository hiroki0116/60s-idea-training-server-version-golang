package routes

import (
	"idea-training-version-go/internals/services"

	"github.com/gin-gonic/gin"
)

func UtilsRoutes(rg *gin.RouterGroup) {
	rg.POST("/error-message", services.LogClientError)
	rg.GET("/healthcheck", services.HealthCheck)
}
