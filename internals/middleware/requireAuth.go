package middleware

import (
	"context"
	"fmt"
	"idea-training-version-go/internals/controllers"
	"idea-training-version-go/internals/utils"
	"idea-training-version-go/internals/utils/firebase"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type RequireAuth struct {
	UserController controllers.IUserController
}

func NewRequireAuth(usercontroller controllers.IUserController) RequireAuth {
	return RequireAuth{
		UserController: usercontroller,
	}
}

func (r *RequireAuth) AllowIfLogIn(ctx *gin.Context) {
	var email string
	auth := ctx.GetHeader("Authorization")
	if auth == "" {
		res := utils.NewHttpResponse(http.StatusUnauthorized, "Invalid authorization token provided...")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}
	tokenString := strings.TrimPrefix(auth, "Bearer ")
	if os.Getenv("STAGE") == "test" {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			res := utils.NewHttpResponse(http.StatusUnauthorized, "Invalid token auth for test env...")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			email = claims["email"].(string)
		}
	} else {
		// For dev and prod environment
		decodedToken, err := firebase.Client.VerifyIDToken(ctx, tokenString)
		if err != nil {
			res := utils.NewHttpResponse(http.StatusUnauthorized, errors.Wrap(err, "Invalid authorization token provided..."))
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		email = decodedToken.Claims["email"].(string)
	}

	user, err := r.UserController.GetUserByEmail(email)
	if err != nil {
		res := utils.NewHttpResponse(http.StatusUnauthorized, errors.Wrap(err, "Failed to get user from mongo db in auth middleware function"))
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}
	// for graphql
	c := context.WithValue(ctx.Request.Context(), "id", user.ID)
	ctx.Request = ctx.Request.WithContext(c)

	// for rest api
	ctx.Set("id", user.ID)
	ctx.Set("email", user.Email)
	ctx.Next()
}
