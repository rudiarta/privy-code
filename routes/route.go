package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rudiarta/privy-code/controller"
)

// InitRoutes is function initialize all routes
func InitRoutes(app *gin.Engine) {
	router := app
	userRouter := router.Group("/user")
	{
		userRouter.POST("/add", controller.AddUserController)
	}
}
