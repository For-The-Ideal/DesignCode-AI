package routes

import (
	"frontend_api/controllers"
	"frontend_api/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(rg *gin.RouterGroup) {
	userController := &controllers.UserController{}

	user := rg.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/info", userController.GetUserInfo)
		user.POST("/updatepassword", userController.UpdateUserPassword)
	}
}
