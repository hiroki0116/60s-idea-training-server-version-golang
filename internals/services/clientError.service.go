package services

import (
	"log"

	"github.com/gin-gonic/gin"
	errors "github.com/pkg/errors"
)

func LogClientError(ctx *gin.Context) {
	type RequestBody struct {
		Message string `json:"message"`
	}
	var req RequestBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Fatal(errors.Wrap(err, "Request body is not valid"))
		return
	}
	log.Println("Client Error====>", req.Message)
}
