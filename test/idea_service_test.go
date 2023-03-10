package test

import (
	"fmt"
	"idea-training-version-go/internals/models"
	"testing"
	"time"

	unitTest "github.com/Valiben/gin_unit_test"
	"github.com/Valiben/gin_unit_test/utils"
	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getSampleIdea() (*models.Idea, error) {
	// get one test user id
	var ideas []*models.Idea
	cursor, err := ideacollection.Find(ctx, bson.D{{}}, options.Find().SetLimit(1))
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var idea *models.Idea
		err := cursor.Decode(&idea)
		if err != nil {
			return nil, err
		}
		ideas = append(ideas, idea)
	}
	return ideas[0], nil
}

func TestCreateIdea(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int         `json:"status"`
		Success    bool        `json:"success"`
		Message    string      `json:"message"`
		Data       models.Idea `json:"data"`
	}

	type IdeaParams struct {
		TopicTitle string    `json:"topicTitle"`
		Ideas      *[]string `json:"ideas"`
		Category   string    `json:"category"`
	}

	var res HTTPResponse

	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestCreateIdea: Fails to add auth header %v\n", err)
		return
	}

	params := IdeaParams{
		TopicTitle: "test_topic_title_11",
		Ideas:      &[]string{"test_idea_1", "test_idea_2"},
		Category:   "test_category_11",
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.POST, "/api/ideas/", "json", params, &res); err != nil {
		t.Errorf("TestCreateIdea: %v\n", err)
		return
	}

	if !res.Success {
		t.Errorf("TestCreateIdea: %v\n", res.Success)
		return
	}

	if res.Data.TopicTitle != params.TopicTitle {
		t.Errorf("TestCreateIdea: expected email %v, got %v\n", params.TopicTitle, res.Data.TopicTitle)
	}

	t.Log("passed")
}

func TestGetAllIdeas(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int           `json:"status"`
		Success    bool          `json:"success"`
		Message    string        `json:"message"`
		Data       []models.Idea `json:"data"`
	}

	var res HTTPResponse

	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestCreateIdea: Fails to add auth header %v\n", err)
		return
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/api/ideas/", "json", nil, &res); err != nil {
		t.Errorf("TestGetAllIdeas: %v\n", err)
		return
	}

	if !res.Success {
		t.Errorf("TestGetAllIdeas: %v\n", res.Success)
		return
	}

	if len(res.Data) == 0 {
		t.Errorf("TestGetAllIdeas: expected ideas count > 0, got %v\n", len(res.Data))
	}

	t.Log("passed")
}

func TestGetIdeaByID(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int         `json:"status"`
		Success    bool        `json:"success"`
		Message    string      `json:"message"`
		Data       models.Idea `json:"data"`
	}

	var res HTTPResponse

	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestCreateIdea: Fails to add auth header %v\n", err)
		return
	}

	idea, err := getSampleIdea()
	if err != nil {
		t.Errorf("TestGetIdeaByID: Failed to get sample idea data...%v\n", err)
		return
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.GET, fmt.Sprintf("/api/ideas/%v", idea.ID.Hex()), "json", nil, &res); err != nil {
		t.Errorf("TestGetIdeaByID: %v\n", err)
		return
	}

	if !res.Success {
		t.Errorf("TestGetIdeaByID: %v\n", res.Success)
		return
	}

	if len(*res.Data.Ideas) == 0 {
		t.Errorf("TestGetIdeaByID: expected ideas count > 0, got %v\n", len(*res.Data.Ideas))
		return
	}

	if res.Data.CreatedBy != idea.CreatedBy {
		t.Errorf("TestGetIdeaByID: expected created by %v, got %v\n", idea.CreatedBy, res.Data.CreatedBy)
		return
	}

	t.Log("passed")
}

func TestUpdateIdea(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int         `json:"status"`
		Success    bool        `json:"success"`
		Message    string      `json:"message"`
		Data       models.Idea `json:"data"`
	}

	type IdeaParams struct {
		TopicTitle string `json:"topicTitle"`
		Category   string `json:"category"`
		Viewed     bool   `json:"viewed"`
		IsLiked    bool   `json:"isLiked"`
	}

	var res HTTPResponse
	var params IdeaParams

	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestCreateIdea: Fails to add auth header %v\n", err)
		return
	}

	params = IdeaParams{
		TopicTitle: "updated title",
		Category:   "updated category",
		Viewed:     true,
		IsLiked:    true,
	}

	idea, err := getSampleIdea()
	if err != nil {
		t.Errorf("TestGetIdeaByID: Failed to get sample idea data...%v\n", err)
		return
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.PUT, fmt.Sprintf("/api/ideas/%v", idea.ID.Hex()), "json", params, &res); err != nil {
		t.Errorf("TestGetIdeaByID: %v\n", err)
		return
	}

	if !res.Success {
		t.Errorf("TestGetIdeaByID: %v\n", res.Success)
		return
	}

	if !*res.Data.IsLiked {
		t.Errorf("TestGetIdeaByID: expected isLiked %v, got %v\n", true, *res.Data.IsLiked)
		return
	}

	if !*res.Data.Viewed {
		t.Errorf("TestGetIdeaByID: expected viewed %v, got %v\n", true, *res.Data.Viewed)
		return
	}

	if res.Data.Category != "updated category" {
		t.Errorf("TestGetIdeaByID: expected category %v, got %v\n", "updated category", res.Data.Category)
		return
	}

	if res.Data.TopicTitle != "updated title" {
		t.Errorf("TestGetIdeaByID: expected topic title %v, got %v\n", "updated title", res.Data.TopicTitle)
		return
	}

	t.Log("passed")

}

func TestDeleteIdea(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int    `json:"status"`
		Success    bool   `json:"success"`
		Message    string `json:"message"`
		Data       string `json:"data"`
	}

	var res HTTPResponse

	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestCreateIdea: Fails to add auth header %v\n", err)
		return
	}

	idea, err := getSampleIdea()
	if err != nil {
		t.Errorf("TestGetIdeaByID: Failed to get sample idea data...%v\n", err)
		return
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.DELETE, fmt.Sprintf("/api/ideas/%v", idea.ID.Hex()), "json", nil, &res); err != nil {
		t.Errorf("TestGetIdeaByID: %v\n", err)
		return
	}

	if !res.Success {
		t.Errorf("TestGetIdeaByID: %v\n", res.Success)
		return
	}

	if res.Data != "Idea deleted successfully" {
		t.Errorf("TestGetIdeaByID: expected message %v, got %v\n", "Idea deleted successfully", res.Data)
		return
	}

	t.Log("passed")
}

func TestGetTotalIdeasOfToday(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int           `json:"status"`
		Success    bool          `json:"success"`
		Message    string        `json:"message"`
		Data       []primitive.M `json:"data"`
	}
	var res HTTPResponse

	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestGetTotalIdeasOfToday: Fails to add auth header %v\n", err)
		return
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/api/ideas/total/today", "json", nil, &res); err != nil {
		t.Errorf("TestGetTotalIdeasOfToday: %v\n", err)
		return
	}

	if !res.Success {
		t.Errorf("TestGetTotalIdeasOfToday: %v\n", res.Success)
		return
	}

	t.Log("passed")
}

func TestGetTotalIdeasOfAllTime(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int           `json:"status"`
		Success    bool          `json:"success"`
		Message    string        `json:"message"`
		Data       []primitive.M `json:"data"`
	}
	var res HTTPResponse

	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestGetTotalIdeasOfAllTime: Fails to add auth header %v\n", err)
		return
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/api/ideas/total/all", "json", nil, &res); err != nil {
		t.Errorf("TestGetTotalIdeasOfAllTime: %v\n", err)
		return
	}

	if !res.Success {
		t.Errorf("TestGetTotalIdeasOfAllTime: %v\n", res.Success)
		return
	}

	t.Log("passed")
}

func TestGetTotalConsecutiveDays(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int    `json:"status"`
		Success    bool   `json:"success"`
		Message    string `json:"message"`
		Data       int    `json:"data"`
	}
	var res HTTPResponse

	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestGetTotalConsecutiveDays: Fails to add auth header %v\n", err)
		return
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/api/ideas/total/consecutive", "json", nil, &res); err != nil {
		t.Errorf("TestGetTotalConsecutiveDays: %v\n", err)
		return
	}

	if res.Data != 0 {
		t.Errorf("TestGetTotalConsecutiveDays: expected data %v, got %v\n", 0, res.Data)
		return
	}

	if !res.Success {
		t.Errorf("TestGetTotalConsecutiveDays: %v\n", res.Success)
		return
	}

	t.Log("passed")
}

func TestGetRecentIdeas(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int           `json:"status"`
		Success    bool          `json:"success"`
		Message    string        `json:"message"`
		Data       []models.Idea `json:"data"`
	}
	var res HTTPResponse
	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestGetRecentIdeas: Fails to add auth header %v\n", err)
		return
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/api/ideas/recent", "json", nil, &res); err != nil {
		t.Errorf("TestGetRecentIdeas: %v\n", err)
		return
	}

	if !res.Success {
		t.Errorf("TestGetRecentIdeas: %v\n", res.Success)
		return
	}

	if len(res.Data) != 5 {
		t.Errorf("TestGetRecentIdeas: expected %v, got %v\n", 5, len(res.Data))
		return
	}

	t.Log("passed")
}

func TestGetWeeklyIdeas(t *testing.T) {
	type WeeklyBody struct {
		TotalIdeas    int    `json:"totalIdeas"`
		TotalSessions int    `json:"totalSessions"`
		ID            string `json:"_id"`
	}
	type ResponseBody struct {
		WeeklyRecords []WeeklyBody `json:"weeklyRecords"`
		LastMonday    time.Time    `json:"lastMonday"`
	}
	type HTTPResponse struct {
		StatusCode int          `json:"status"`
		Success    bool         `json:"success"`
		Message    string       `json:"message"`
		Data       ResponseBody `json:"data"`
	}

	var res HTTPResponse
	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestGetWeeklyIdeas: Fails to add auth header %v\n", err)
		return
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.GET, "/api/ideas/weekly", "json", nil, &res); err != nil {
		t.Errorf("TestGetWeeklyIdeas: %v\n", err)
		return
	}

	if res.Data.WeeklyRecords[0].TotalIdeas < 1 {
		t.Errorf("TestGetWeeklyIdeas: expected %v, got %v\n", 1, res.Data.WeeklyRecords[0].TotalIdeas)
		return
	}

	if res.Data.WeeklyRecords[0].TotalSessions < 1 {
		t.Errorf("TestGetWeeklyIdeas: expected %v, got %v\n", 1, res.Data.WeeklyRecords[0].TotalSessions)
		return
	}

	if !res.Success {
		t.Errorf("TestGetWeeklyIdeas: %v\n", res.Success)
		return
	}

	t.Log("passed")

}

func TestSearchIdeas(t *testing.T) {
	type ResponseBody struct {
		Ideas        []models.Idea           `json:"ideas"`
		PaginateData *paginate.PaginatedData `json:"paginateData"`
	}

	type SearchParams struct {
		SearchInput   string    `json:"searchInput,omitempty"`
		Category      string    `json:"category,omitempty"`
		CreatedAtFrom time.Time `json:"createdAtFrom,omitempty"`
		CreatedAtTo   time.Time `json:"createdAtTo,omitempty"`
		Pagesize      int       `json:"pageSize,omitempty"`
		Current       int       `json:"current"`
		SortByRecent  bool      `json:"sortByRecent,omitempty"`
		IsLiked       bool      `json:"isLiked,omitempty"`
	}

	type HTTPResponse struct {
		StatusCode int          `json:"status"`
		Success    bool         `json:"success"`
		Message    string       `json:"message"`
		Data       ResponseBody `json:"data"`
	}
	var res HTTPResponse
	if _, err := AddAuthHeader(); err != nil {
		t.Errorf("TestSearchIdeas: Fails to add auth header %v\n", err)
		return
	}

	params := SearchParams{
		SearchInput:  "test",
		Pagesize:     9,
		Current:      1,
		IsLiked:      false,
		SortByRecent: true,
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.POST, "/api/ideas/search", "json", params, &res); err != nil {
		t.Errorf("TestSearchIdeas: %v\n", err)
		return
	}

	if res.Data.PaginateData.Pagination.PerPage != 9 {
		t.Errorf("TestSearchIdeas: expected %v, got %v\n", 9, len(res.Data.Ideas))
		return
	}

	if len(res.Data.Ideas) != 9 {
		t.Errorf("TestSearchIdeas: expected %v, got %v\n", 9, len(res.Data.Ideas))
		return
	}

	if !res.Success {
		t.Errorf("TestSearchIdeas: %v\n", res.Success)
		return
	}

	t.Log("passed")
}
