package api

import (
	"UGUAPI/api/middleware"
	"UGUAPI/auth"
	"UGUAPI/cmd/controllers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(
	AuthController *controllers.AuthController,
	JWTManager *auth.JWTManager) *gin.Engine {
	router := gin.Default()

	baseRouter := router.Group("api")
	{
		user := baseRouter.Group("user")
		{
			user.POST("login", AuthController.Login)
			user.POST("register", AuthController.Register)
			user.GET("get", middleware.SetAuthorizationCheck(JWTManager), AuthController.GetUser)
		}
	}

	return router
}
