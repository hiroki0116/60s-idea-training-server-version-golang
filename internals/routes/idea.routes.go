package routes

import (
	"idea-training-version-go/internals/middleware"
	"idea-training-version-go/internals/services"

	"github.com/gin-gonic/gin"
)

type IdeaRoutes struct {
	IdeaService services.IIdeaService
	RequireAuth middleware.RequireAuth
}

func NewIdeaRoutes(ideaService services.IIdeaService, requireAuth middleware.RequireAuth) IdeaRoutes {
	return IdeaRoutes{
		IdeaService: ideaService,
		RequireAuth: requireAuth,
	}
}

func (ir *IdeaRoutes) IdeaRoutes(rg *gin.RouterGroup) {
	idearoute := rg.Group("/ideas")

	idearoute.POST("/", ir.RequireAuth.AllowIfLogIn, ir.IdeaService.CreateIdea)
	idearoute.GET("/", ir.RequireAuth.AllowIfLogIn, ir.IdeaService.GetAllIdeas)
	idearoute.GET("/:id", ir.RequireAuth.AllowIfLogIn, ir.IdeaService.GetIdeaByID)
	idearoute.PUT("/:id", ir.RequireAuth.AllowIfLogIn, ir.IdeaService.UpdateIdea)
	idearoute.DELETE("/:id", ir.RequireAuth.AllowIfLogIn, ir.IdeaService.DeleteIdea)
	idearoute.GET("/total-today", ir.RequireAuth.AllowIfLogIn, ir.IdeaService.GetTotalIdeasOfToday)
}
