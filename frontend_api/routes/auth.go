package routes

import (
	"frontend_api/controllers"
	"frontend_api/utils"

	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(rg *gin.RouterGroup) {
	authController := controllers.NewAuthController(utils.DB)

	auth := rg.Group("/")
	{
		auth.GET("/captcha", authController.Captcha)
		auth.GET("/template", authController.Template)
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}
}
