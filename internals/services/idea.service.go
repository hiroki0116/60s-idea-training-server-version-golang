package services

import (
	"idea-training-version-go/internals/controllers"
	"idea-training-version-go/internals/models"
	"idea-training-version-go/internals/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	paginate "github.com/gobeam/mongo-go-pagination"
	errors "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IIdeaService interface {
	CreateIdea(ctx *gin.Context)
	GetAllIdeas(ctx *gin.Context)
	GetIdeaByID(ctx *gin.Context)
	UpdateIdea(ctx *gin.Context)
	DeleteIdea(ctx *gin.Context)
	GetTotalIdeasOfToday(ctx *gin.Context)
	GetTotalIdeasOfAllTime(ctx *gin.Context)
	GetTotalConsecutiveDays(ctx *gin.Context)
	GetRecentIdeas(ctx *gin.Context)
	GetWeeklyIdeas(ctx *gin.Context)
	SearchIdeas(ctx *gin.Context)
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

func (is *IdeaService) GetTotalIdeasOfToday(ctx *gin.Context) {
	userID := utils.FetchUserFromCtx(ctx)

	result, err := is.IdeaController.GetTotalIdeasOfToday(userID)

	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in getting ideas"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.NewHttpResponse(http.StatusOK, result)
	ctx.JSON(http.StatusOK, res)
}

func (is *IdeaService) GetTotalIdeasOfAllTime(ctx *gin.Context) {
	userID := utils.FetchUserFromCtx(ctx)

	result, err := is.IdeaController.GetTotalIdeasOfAllTime(userID)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in getting ideas of all time"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.NewHttpResponse(http.StatusOK, result)
	ctx.JSON(http.StatusOK, res)
}

func (is *IdeaService) GetTotalConsecutiveDays(ctx *gin.Context) {
	userID := utils.FetchUserFromCtx(ctx)

	result, err := is.IdeaController.GetTotalConsecutiveDays(userID)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in getting total consecutive days"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.NewHttpResponse(http.StatusOK, result)
	ctx.JSON(http.StatusOK, res)
}

func (is *IdeaService) GetRecentIdeas(ctx *gin.Context) {
	userID := utils.FetchUserFromCtx(ctx)

	result, err := is.IdeaController.GetRecentIdeas(userID)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in getting total ideas of 5 recent ideas"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.NewHttpResponse(http.StatusOK, result)
	ctx.JSON(http.StatusOK, res)
}

func (is *IdeaService) GetWeeklyIdeas(ctx *gin.Context) {
	userID := utils.FetchUserFromCtx(ctx)

	result, lastMonday, err := is.IdeaController.GetWeeklyIdeas(userID)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in getting weekly ideas"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	type ResponseBody struct {
		WeeklyRecords []primitive.M `json:"weeklyRecords"`
		LastMonday    time.Time     `json:"lastMonday"`
	}

	res := utils.NewHttpResponse(http.StatusOK, &ResponseBody{WeeklyRecords: result, LastMonday: lastMonday})
	ctx.JSON(http.StatusOK, res)
}

func (is *IdeaService) SearchIdeas(ctx *gin.Context) {

	type RequestBody struct {
		SearchInput   string    `json:"searchInput,omitempty"`
		Category      string    `json:"category,omitempty"`
		CreatedAtFrom time.Time `json:"createdAtFrom,omitempty"`
		CreatedAtTo   time.Time `json:"createdAtTo,omitempty"`
		Pagesize      int       `json:"pageSize,omitempty"`
		Current       int       `json:"current"`
		SortByRecent  bool      `json:"sortByRecent,omitempty"`
		IsLiked       bool      `json:"isLiked,omitempty"`
	}

	var req RequestBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Request body is not valid"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	userID := utils.FetchUserFromCtx(ctx)

	// matchStage
	filter := bson.M{}
	filter["createdBy"] = userID
	// filtering createdAt duration
	if req.CreatedAtFrom.IsZero() && !req.CreatedAtTo.IsZero() {
		filter["createdAt"] = bson.M{"$lte": req.CreatedAtTo}
	}

	if !req.CreatedAtFrom.IsZero() && req.CreatedAtTo.IsZero() {
		filter["createdAt"] = bson.M{"$gte": req.CreatedAtFrom}
	}

	if !req.CreatedAtFrom.IsZero() && !req.CreatedAtTo.IsZero() {
		filter["createdAt"] = bson.M{"$gte": req.CreatedAtFrom, "$lte": req.CreatedAtTo}
	}

	if req.Category != "" {
		filter["category"] = req.Category
	}

	if req.IsLiked {
		filter["isLiked"] = true
	}

	if req.SearchInput != "" {
		filter["$or"] = []bson.M{
			{"topicTitle": bson.M{"$regex": req.SearchInput, "$options": "i"}},
			{"ideas": bson.M{"$regex": req.SearchInput, "$options": "i"}},
			{"category": bson.M{"$regex": req.SearchInput, "$options": "i"}},
		}
	}

	// sort
	var sort int
	if req.SortByRecent {
		sort = -1
	} else {
		sort = 1
	}

	// page and size
	if req.Current == 0 {
		req.Current = 1
	}
	if req.Pagesize == 0 {
		req.Pagesize = 9
	}

	results, paginateData, err := is.IdeaController.Search(filter, sort, req.Current, req.Pagesize)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in searching ideas"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	type ResponseBody struct {
		Ideas        []models.Idea           `json:"ideas"`
		PaginateData *paginate.PaginatedData `json:"paginateData"`
	}

	resBody := ResponseBody{
		Ideas:        results,
		PaginateData: paginateData,
	}

	res := utils.NewHttpResponse(http.StatusOK, resBody)
	ctx.JSON(http.StatusOK, res)
}
