package routes

import (
	"github.com/GdAyo19/FundsFlow/controllers"
	"github.com/GdAyo19/FundsFlow/middleware"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(router *gin.Engine) {

	api := router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.GET("/users", controllers.GetUsers)
		auth.POST("/login", controllers.Login)
	}

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	{
		protected.GET("/profile", controllers.Profile)
		protected.POST("/transactions", controllers.CreateTransaction)
		protected.GET("/transactions", controllers.GetTransactions)
		protected.PUT("/transactions/:id", controllers.UpdateTransaction)
		protected.DELETE("/transactions/:id", controllers.DeleteTransaction)
		protected.POST("/budgets", controllers.CreateBudget)
		protected.GET("/budgets", controllers.GetBudgets)
		protected.POST("/goals", controllers.CreateSavingsGoal)
		protected.POST("/goals/:id/contributions", controllers.AddContribution)
		protected.GET("/goals", controllers.GetSavingsGoals)

	}

}



