package services

import (
	"context"
	"idea-training-version-go/internals/controllers"
	"idea-training-version-go/internals/db"
	"idea-training-version-go/internals/models"
	"idea-training-version-go/internals/utils"
	"idea-training-version-go/internals/utils/firebase"
	"net/http"

	"github.com/gin-gonic/gin"
	errors "github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserService interface {
	SignUp(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
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
	var newUser *models.User
	var session mongo.Session
	var err error

	if err = ctx.ShouldBindJSON(&req); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Request body is not valid"))
		ctx.JSON(http.StatusBadRequest, res)
	}

	if session, err = db.MongoDB.StartSession(); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in creating mongo session"))
		ctx.JSON(http.StatusBadRequest, res)
	}

	if err = session.StartTransaction(); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in starting transaction"))
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

	if err := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		// register in mongodb
		newUser, err = us.UserController.CreateUser(&user)
		if err != nil {
			res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in creating user in mongodb"))
			ctx.JSON(http.StatusBadRequest, res)
			return err
		}

		if err = session.CommitTransaction(sc); err != nil {
			res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in committing transaction"))
			ctx.JSON(http.StatusBadRequest, res)
		}
		return nil
	}); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in creating session"))
		ctx.JSON(http.StatusBadRequest, res)
	}

	session.EndSession(context.Background())

	res := utils.NewHttpResponse(http.StatusOK, newUser)
	ctx.JSON(http.StatusOK, res)
}

func (us *UserService) UpdateUser(ctx *gin.Context) {
	type RequestBody struct {
		Email     string `json:"email,omitempty"`
		FirstName string `json:"firstName,omitempty"`
		LastName  string `json:"lastName,omitempty"`
	}

	var req RequestBody
	var err error
	var updatedUser *models.User

	// Convert id to primitive.ObjectID
	id := ctx.Param("id")
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in converting id to primitive.ObjectID"))
		ctx.JSON(http.StatusBadRequest, res)
	}

	// Bind json
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Request body is not valid"))
		ctx.JSON(http.StatusBadRequest, res)
	}

	user := &models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// update in mongodb
	if err = us.UserController.UpdateUser(userID, user); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in updating user in mongodb"))
		ctx.JSON(http.StatusBadRequest, res)
	}

	if updatedUser, err = us.UserController.GetUserByID(userID); err != nil {
		res := utils.NewHttpResponse(http.StatusBadRequest, errors.Wrap(err, "Error in getting updateduser from mongodb"))
		ctx.JSON(http.StatusBadRequest, res)
	}

	res := utils.NewHttpResponse(http.StatusOK, updatedUser)
	ctx.JSON(http.StatusOK, res)
}
