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

		// Transactions
		protected.POST("/transactions", controllers.CreateTransaction)
		protected.GET("/transactions", controllers.GetTransactions)
		protected.PUT("/transactions/:id", controllers.UpdateTransaction)
		protected.DELETE("/transactions/:id", controllers.DeleteTransaction)

		// Budgets
		protected.POST("/budgets", controllers.CreateBudget)
		protected.GET("/budgets", controllers.GetBudgets)
		protected.PUT("/budgets/:id", controllers.UpdateBudget)
		protected.DELETE("/budgets/:id", controllers.DeleteBudget)

		// Savings Goals
		protected.POST("/goals", controllers.CreateSavingsGoal)
		protected.GET("/goals", controllers.GetSavingsGoals)
		protected.PUT("/goals/:id", controllers.UpdateSavingsGoal)
		protected.DELETE("/goals/:id", controllers.DeleteSavingsGoal)

		// Goal Contributions
		protected.POST("/goals/:id/contributions", controllers.AddContribution)
		protected.DELETE("/goals/:id/contributions/:contributionId", controllers.DeleteContribution)

		// Dashboard
		protected.GET("/dashboard", controllers.Dashboard)

		// Reports
		protected.GET("/reports/monthly", controllers.MonthlyReport)
		protected.GET("/reports/category", controllers.CategoryReportHandler)
		protected.GET("/reports/export/csv", controllers.ExportCSV)
		protected.GET("/reports/export/pdf", controllers.ExportPDF)

		// Recurring Transactions
		protected.POST("/recurring-transactions", controllers.CreateRecurringTransaction)
		protected.GET("/recurring-transactions", controllers.GetRecurringTransactions)
		protected.PUT("/recurring-transactions/:id", controllers.UpdateRecurringTransaction)
		protected.DELETE("/recurring-transactions/:id", controllers.DeleteRecurringTransaction)

		// Notifications
		protected.GET("/notifications", controllers.GetNotifications)
	}
}
