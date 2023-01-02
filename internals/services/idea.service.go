package services

import (
	"idea-training-version-go/internals/controllers"
	"idea-training-version-go/internals/models"
	"idea-training-version-go/internals/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IIdeaService interface {
	CreateIdea(ctx *gin.Context)
}

type IdeaService struct {
	IdeaController controllers.IIdeaController
}

func NewIdeaService(ideaController controllers.IIdeaController) IIdeaService {
	return &IdeaService{
		IdeaController: ideaController,
	}
}

func (is *IdeaService) CreateIdea(ctx *gin.Context) {
	userID := utils.FetchUserFromCtx(ctx)
	var idea models.Idea
	var newIdea *models.Idea
	if err := ctx.ShouldBindJSON(&idea); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	idea.CreatedBy = userID
	newIdea, err := is.IdeaController.CreateIdea(&idea)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, err)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.NewHttpResponse(http.StatusCreated, newIdea)
	ctx.JSON(http.StatusCreated, res)
}
