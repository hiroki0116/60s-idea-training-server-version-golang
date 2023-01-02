package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FetchUserFromCtx(ctx *gin.Context) primitive.ObjectID {
	uId, _ := ctx.Get("id")
	userId, err := primitive.ObjectIDFromHex(uId.(primitive.ObjectID).Hex())
	if err != nil {
		res := NewHttpResponse(http.StatusBadRequest, err)
		ctx.JSON(http.StatusBadRequest, res)
		return primitive.NilObjectID
	}
	return userId
}
