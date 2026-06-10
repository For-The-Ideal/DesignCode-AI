package routes

import (
	"frontend_api/controllers"
	"frontend_api/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(rg *gin.RouterGroup) {
	userController := &controllers.UserController{}

	user := rg.Group("/")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/user/info", userController.GetUserInfo)
		user.POST("/user/updatepassword", userController.UpdateUserPassword)
	}
}
