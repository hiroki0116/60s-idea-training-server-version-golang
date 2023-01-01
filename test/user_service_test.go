package test

import (
	"idea-training-version-go/internals/models"
	"testing"

	unitTest "github.com/Valiben/gin_unit_test"
	"github.com/Valiben/gin_unit_test/utils"
)

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
