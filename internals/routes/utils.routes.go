package routes

import (
	"idea-training-version-go/internals/services"

	"github.com/gin-gonic/gin"
)

func UtilsRoutes(rg *gin.RouterGroup) {
	clienterrorroute := rg.Group("/error-message")
	clienterrorroute.POST("/", services.LogClientError)
	clienterrorroute.GET("/healthcheck", services.HealthCheck)
}
