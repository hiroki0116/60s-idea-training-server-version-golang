package routes

import (
	"idea-training-version-go/internals/services"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	UserService services.IUserService
}

func NewUserRoutes(userService services.IUserService) UserRoutes {
	return UserRoutes{
		UserService: userService,
	}
}

func (ur *UserRoutes) UserRoutes(rg *gin.RouterGroup) {
	userroute := rg.Group("/users")

	userroute.POST("/signup", ur.UserService.SignUp)
	userroute.PUT("/update/:id", ur.UserService.UpdateUser)
}
