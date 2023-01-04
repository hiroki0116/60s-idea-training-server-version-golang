package routes

import (
	"idea-training-version-go/internals/services"

	"github.com/gin-gonic/gin"
)

func ClientErrorRoutes(rg *gin.RouterGroup) {
	clienterrorroute := rg.Group("/error-message")

	clienterrorroute.POST("/", services.LogClientError)
}
