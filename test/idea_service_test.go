package test

import (
	"fmt"
	"idea-training-version-go/internals/models"
	"log"
	"testing"

	unitTest "github.com/Valiben/gin_unit_test"
	"github.com/Valiben/gin_unit_test/utils"
)

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

	user := getSampleUser()
	tokenString, err := GenerateJWTToken(user.Email)
	if err != nil {
		log.Fatal("Error in generating JWT token: ", err)
		return
	}
	unitTest.AddHeader("Authorization", fmt.Sprintf("Bearer %s", tokenString))

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
