package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rudiarta/privy-code/controller"
	"github.com/rudiarta/privy-code/middleware"
)

// InitRoutes is function initialize all routes
func InitRoutes(app *gin.Engine) {
	router := app

	//userRouter is routes for user endpoint
	userRouter := router.Group("/user")
	{
		userRouter.POST("/add", middleware.AuthMiddleware, controller.AddUserController)
		userRouter.POST("/login", controller.LoginUserController)
	}

	balanceRouter := router.Group("balance")
	{
		balanceRouter.POST("/add", middleware.AuthMiddleware, controller.AddBalanceController)
		balanceRouter.POST("/takeout", middleware.AuthMiddleware, controller.TakeOutBalanceController)
		balanceRouter.POST("/transfer", middleware.AuthMiddleware, controller.TransferBalanceController)
	}
}
