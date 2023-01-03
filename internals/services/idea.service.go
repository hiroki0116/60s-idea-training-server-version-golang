package services

import (
	"idea-training-version-go/internals/controllers"
	"idea-training-version-go/internals/models"
	"idea-training-version-go/internals/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	errors "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IIdeaService interface {
	CreateIdea(ctx *gin.Context)
	GetAllIdeas(ctx *gin.Context)
	GetIdeaByID(ctx *gin.Context)
	UpdateIdea(ctx *gin.Context)
	DeleteIdea(ctx *gin.Context)
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
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Request body is not valid"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	idea.CreatedBy = userID
	newIdea, err := is.IdeaController.CreateIdea(&idea)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in creating idea"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.NewHttpResponse(http.StatusCreated, newIdea)
	ctx.JSON(http.StatusCreated, res)
}

func (is *IdeaService) GetAllIdeas(ctx *gin.Context) {
	userID := utils.FetchUserFromCtx(ctx)

	ideas, err := is.IdeaController.GetAllIdeas(userID)

	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in getting ideas"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.NewHttpResponse(http.StatusOK, ideas)
	ctx.JSON(http.StatusOK, res)
}

func (is *IdeaService) GetIdeaByID(ctx *gin.Context) {
	id := ctx.Param("id")
	ideaID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Invalid idea id"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	idea, err := is.IdeaController.GetIdeaByID(ideaID)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in getting idea"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.NewHttpResponse(http.StatusOK, idea)
	ctx.JSON(http.StatusOK, res)
}

func (is *IdeaService) UpdateIdea(ctx *gin.Context) {
	id := ctx.Param("id")
	ideaID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Invalid idea id"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	var idea models.Idea
	if err := ctx.ShouldBindJSON(&idea); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Request body is not valid"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	idea.ID = ideaID

	if err := is.IdeaController.UpdateIdea(&idea); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in updating idea"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	updatedIdea, err := is.IdeaController.GetIdeaByID(ideaID)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in getting updated idea"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.NewHttpResponse(http.StatusOK, updatedIdea)
	ctx.JSON(http.StatusOK, res)
}

func (is *IdeaService) DeleteIdea(ctx *gin.Context) {
	id := ctx.Param("id")
	ideaID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Invalid idea id"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	if err := is.IdeaController.DeleteIdea(ideaID); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in deleting idea"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.NewHttpResponse(http.StatusOK, "Idea deleted successfully")
	ctx.JSON(http.StatusOK, res)
}
