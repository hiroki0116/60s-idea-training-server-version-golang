package test

import (
	"fmt"
	"idea-training-version-go/internals/models"
	"log"
	"testing"

	unitTest "github.com/Valiben/gin_unit_test"
	"github.com/Valiben/gin_unit_test/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddAuthHeader() (*models.User, error) {
	// get one test user id
	var users []*models.User
	cursor, err := usercollection.Find(ctx, bson.D{{}}, options.Find().SetLimit(1))
	if err != nil {
		log.Fatal("Error getting sample users: ", err)
		return nil, err
	}

	for cursor.Next(ctx) {
		var user *models.User
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal("Error decoding sample users: ", err)
			return nil, err
		}
		users = append(users, user)
	}
	tokenString, err := GenerateJWTToken(users[0].Email)
	if err != nil {
		log.Fatal("Error in generating JWT token: ", err)
		return nil, err
	}
	unitTest.AddHeader("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	return users[0], nil
}

func TestSignup(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int         `json:"status"`
		Success    bool        `json:"success"`
		Message    string      `json:"message"`
		Data       models.User `json:"data"`
	}

	type UserParams struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	var res HTTPResponse

	params := UserParams{
		Email:     "test_email11@test.com",
		Password:  "password_11",
		FirstName: "test_first_name_11",
		LastName:  "test_last_name_11",
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.POST, "/api/users/signup", "json", params, &res); err != nil {
		t.Errorf("TestSignup: %v\n", err)
		return
	}

	if !res.Success {
		t.Errorf("TestSignup: %v\n", res.Success)
		return
	}

	if res.Data.Email != params.Email {
		t.Errorf("TestSignup: expected email %v, got %v\n", params.Email, res.Data.Email)
	}

	t.Log("passed")
}

func TestUpdateUser(t *testing.T) {
	type HTTPResponse struct {
		StatusCode int         `json:"status"`
		Success    bool        `json:"success"`
		Message    string      `json:"message"`
		Data       models.User `json:"data"`
	}

	type UserParams struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	var res HTTPResponse

	params := UserParams{
		Email:     "updated_test_email11@test.com",
		FirstName: "updated_test_first_name_11",
		LastName:  "updated_test_last_name_11",
	}

	user, err := AddAuthHeader()
	if err != nil {
		t.Errorf("TestUpdateUser: %v/n", err)
		return
	}

	if err := unitTest.TestHandlerUnMarshalResp(utils.PUT, fmt.Sprintf("/api/users/update/%v", user.ID.Hex()), "json", params, &res); err != nil {
		t.Errorf("TestUpdateUser: %v/n", err)
		return
	}

	if !res.Success {
		t.Errorf("TestUpdateUser: %v\n", res.Success)
		return
	}

	if res.Data.Email != params.Email {
		t.Errorf("TestUpdateUser: expected email %v, got %v\n", params.Email, res.Data.Email)
		return
	}

	t.Log("passed")
}
