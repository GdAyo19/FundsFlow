package controllers

import (
	"net/http"
	"time"

	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/gin-gonic/gin"
)

type DashboardResponse struct {
	TotalIncome     float64                `json:"total_income"`
	TotalExpenses   float64                `json:"total_expenses"`
	CurrentBalance  float64                `json:"current_balance"`
	SavingsProgress []SavingsGoalProgress  `json:"savings_progress"`
	BudgetAlerts    []BudgetAlert          `json:"budget_alerts"`
}

type SavingsGoalProgress struct {
	ID            uint    `json:"id"`
	Title         string  `json:"title"`
	Progress      float64 `json:"progress"`
	SavedAmount   float64 `json:"saved_amount"`
	TargetAmount  float64 `json:"target_amount"`
	Status        string  `json:"status"`
}

type BudgetAlert struct {
	ID        uint    `json:"id"`
	Category  string  `json:"category"`
	Budget    float64 `json:"budget"`
	Spent     float64 `json:"spent"`
	Progress  float64 `json:"progress"`
	Status    string  `json:"status"`
}

func Dashboard(c *gin.Context) {
	userID := c.GetUint("userID")
	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()

	var totalIncome float64
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
			userID, "income", currentMonth, currentYear).
		Select("COALESCE(SUM(amount),0)").Scan(&totalIncome)

	var totalExpenses float64
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
			userID, "expense", currentMonth, currentYear).
		Select("COALESCE(SUM(amount),0)").Scan(&totalExpenses)

	var allTimeIncome float64
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "income").
		Select("COALESCE(SUM(amount),0)").Scan(&allTimeIncome)

	var allTimeExpenses float64
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ?", userID, "expense").
		Select("COALESCE(SUM(amount),0)").Scan(&allTimeExpenses)

	currentBalance := allTimeIncome - allTimeExpenses

	var goals []models.SavingsGoal
	config.DB.Where("user_id = ?", userID).Find(&goals)

	var savingsProgress []SavingsGoalProgress
	for _, goal := range goals {
		var savedAmount float64
		config.DB.Model(&models.GoalContribution{}).
			Where("goal_id = ?", goal.ID).
			Select("COALESCE(SUM(amount),0)").Scan(&savedAmount)

		progress := 0.0
		if goal.TargetAmount > 0 {
			progress = (savedAmount / goal.TargetAmount) * 100
		}
		if progress > 100 {
			progress = 100
		}

		status := "On Track"
		if progress >= 100 {
			status = "Completed"
		} else if time.Now().After(goal.Deadline) {
			status = "Overdue"
		}

		savingsProgress = append(savingsProgress, SavingsGoalProgress{
			ID:           goal.ID,
			Title:        goal.Title,
			Progress:     progress,
			SavedAmount:  savedAmount,
			TargetAmount: goal.TargetAmount,
			Status:       status,
		})
	}

	var budgets []models.Budget
	config.DB.Where("user_id = ? AND month = ? AND year = ?", userID, currentMonth, currentYear).Find(&budgets)

	var budgetAlerts []BudgetAlert
	for _, budget := range budgets {
		var totalSpent float64
		config.DB.Model(&models.Transaction{}).
			Where("user_id = ? AND category = ? AND type = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
				userID, budget.Category, "expense", currentMonth, currentYear).
			Select("COALESCE(SUM(amount),0)").Scan(&totalSpent)

		progress := (totalSpent / budget.Amount) * 100
		status := "On Track"
		if progress >= 100 {
			status = "Exceeded"
		} else if progress >= 75 {
			status = "At Risk"
		}

		if progress >= 75 {
			budgetAlerts = append(budgetAlerts, BudgetAlert{
				ID:       budget.ID,
				Category: budget.Category,
				Budget:   budget.Amount,
				Spent:    totalSpent,
				Progress: progress,
				Status:   status,
			})
		}
	}

	c.JSON(http.StatusOK, DashboardResponse{
		TotalIncome:     totalIncome,
		TotalExpenses:   totalExpenses,
		CurrentBalance:  currentBalance,
		SavingsProgress: savingsProgress,
		BudgetAlerts:    budgetAlerts,
	})
}
