package routes

import (
	"idea-training-version-go/internals/middleware"
	"idea-training-version-go/internals/services"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	UserService services.IUserService
	RequireAuth middleware.RequireAuth
}

func NewUserRoutes(userService services.IUserService, requireAuth middleware.RequireAuth) UserRoutes {
	return UserRoutes{
		UserService: userService,
		RequireAuth: requireAuth,
	}
}

func (ur *UserRoutes) UserRoutes(rg *gin.RouterGroup) {
	userroute := rg.Group("/users")

	userroute.POST("/signup", ur.UserService.SignUp)
	userroute.GET("/", ur.UserService.GetUserByEmail)
	userroute.PUT("/:id", ur.RequireAuth.AllowIfLogIn, ur.UserService.UpdateUser)
	userroute.POST("/images", ur.RequireAuth.AllowIfLogIn, ur.UserService.UploadImageCloudinary)
	userroute.DELETE("/images", ur.RequireAuth.AllowIfLogIn, ur.UserService.RemoveImageCloudinary)
}
