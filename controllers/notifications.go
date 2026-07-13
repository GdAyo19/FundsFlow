package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/gin-gonic/gin"
)

type Notification struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

func GetNotifications(c *gin.Context) {
	userID := c.GetUint("userID")
	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()

	var notifications []Notification

	var budgets []models.Budget
	config.DB.Where("user_id = ? AND month = ? AND year = ?", userID, currentMonth, currentYear).Find(&budgets)

	for _, budget := range budgets {
		var totalSpent float64
		config.DB.Model(&models.Transaction{}).
			Where("user_id = ? AND category = ? AND type = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
				userID, budget.Category, "expense", currentMonth, currentYear).
			Select("COALESCE(SUM(amount),0)").Scan(&totalSpent)

		progress := (totalSpent / budget.Amount) * 100

		if progress >= 100 {
			notifications = append(notifications, Notification{
				Type:     "budget_exceeded",
				Message:  fmt.Sprintf("Budget exceeded for %s! You have spent $%.2f of $%.2f.", budget.Category, totalSpent, budget.Amount),
				Severity: "critical",
			})
		} else if progress >= 75 {
			notifications = append(notifications, Notification{
				Type:     "budget_warning",
				Message:  fmt.Sprintf("Budget at risk for %s. You have spent $%.2f of $%.2f (%.0f%%).", budget.Category, totalSpent, budget.Amount, progress),
				Severity: "warning",
			})
		}
	}

	var goals []models.SavingsGoal
	config.DB.Where("user_id = ?", userID).Find(&goals)

	for _, goal := range goals {
		if now.After(goal.Deadline) {
			var savedAmount float64
			config.DB.Model(&models.GoalContribution{}).
				Where("goal_id = ?", goal.ID).
				Select("COALESCE(SUM(amount),0)").Scan(&savedAmount)

			if savedAmount < goal.TargetAmount {
				notifications = append(notifications, Notification{
					Type:     "goal_overdue",
					Message:  fmt.Sprintf("Savings goal \"%s\" is overdue. Target: $%.2f, Saved: $%.2f.", goal.Title, goal.TargetAmount, savedAmount),
					Severity: "critical",
				})
			}
		} else {
			daysUntilDeadline := int(goal.Deadline.Sub(now).Hours() / 24)
			if daysUntilDeadline <= 7 && daysUntilDeadline > 0 {
				var savedAmount float64
				config.DB.Model(&models.GoalContribution{}).
					Where("goal_id = ?", goal.ID).
					Select("COALESCE(SUM(amount),0)").Scan(&savedAmount)

				if savedAmount < goal.TargetAmount {
					notifications = append(notifications, Notification{
						Type:     "goal_deadline_approaching",
						Message:  fmt.Sprintf("Savings goal \"%s\" deadline is approaching (%d days left). Target: $%.2f, Saved: $%.2f.", goal.Title, daysUntilDeadline, goal.TargetAmount, savedAmount),
						Severity: "warning",
					})
				}
			}
		}
	}

	c.JSON(http.StatusOK, notifications)
}
