package services

import (
	"idea-training-version-go/internals/controllers"
	"idea-training-version-go/internals/models"
	"idea-training-version-go/internals/utils"
	"idea-training-version-go/internals/utils/firebase"
	"net/http"

	"github.com/gin-gonic/gin"
	errors "github.com/pkg/errors"
)

type IUserService interface {
	SignUp(ctx *gin.Context)
}

type UserService struct {
	UserController controllers.IUserController
}

func NewUserService(userController controllers.IUserController) IUserService {
	return &UserService{
		UserController: userController,
	}
}

func (us *UserService) SignUp(ctx *gin.Context) {
	type RequestBody struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	var req RequestBody
	var user models.User

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Request body is not valid"))
		ctx.JSON(http.StatusBadRequest, res)
	}

	// register in firebase
	firebaseUID, err := firebase.CreateUserInFirebase(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in creating user in firebase"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user.FirebaseUID = firebaseUID
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email

	// register in mongodb
	newUser, err := us.UserController.CreateUser(&user)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in creating user in mongodb"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.NewHttpResponse(http.StatusOK, newUser)
	ctx.JSON(http.StatusOK, res)
}
