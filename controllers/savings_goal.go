package controllers

import (
	"net/http"
	"time"

	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/gin-gonic/gin"
)

func CreateSavingsGoal(c *gin.Context) {

	// bind the request body to the body struct
	var body models.SavingGoalRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// parse the deadline string into a time.Time object
	deadline, err := time.Parse("2006-01-02", body.Deadline)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid deadline format. Use YYYY-MM-DD.",
		})
		return
	}

	goal := models.SavingsGoal{
		UserID:       c.GetUint("userID"),
		Title:        body.Title,
		TargetAmount: body.TargetAmount,
		Deadline:     deadline,
	}

	if err := config.DB.Create(&goal).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to create savings goal",
		})
		return
	}

	c.JSON(201, goal)
}

func GetSavingsGoals(c *gin.Context) {
	userID := c.GetUint("userID")

	var goals []models.SavingsGoal

	if err := config.DB.Where("user_id = ?", userID).Find(&goals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var response []models.SavingsGoalResponse

	for _, goal := range goals {

		// calculate total saved
		var savedAmount float64
		config.DB.Model(&models.GoalContribution{}).Where("goal_id = ?", goal.ID).Select("COALESCE(SUM(amount), 0)").Scan(&savedAmount)

		// remaining amount

		remainingAmount := goal.TargetAmount - savedAmount
		if remainingAmount < 0 {
			remainingAmount = 0
		}

		// progress
		progress := 0.0
		if goal.TargetAmount > 0 {
			progress = (savedAmount / goal.TargetAmount) * 100
		}

		if progress > 100 {
			progress = 100
		}

		// days remaining

		daysRemaining := int(goal.Deadline.Sub(time.Now()).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		// Approximate months remaining
		monthsRemaining := int(goal.Deadline.Sub(time.Now()).Hours() / (24 * 30))
		if monthsRemaining < 1 {
			monthsRemaining = 1
		}

		// Monthly savings needed
		monthlyNeeded := remainingAmount / float64(monthsRemaining)

		// Goal status
		status := "On Track"

		if progress >= 100 {
			status = "Completed"
		} else if daysRemaining == 0 {
			status = "Overdue"
		}

		response = append(response, models.SavingsGoalResponse{
			ID:                   goal.ID,
			Title:                goal.Title,
			TargetAmount:         goal.TargetAmount,
			SavedAmount:          savedAmount,
			RemainingAmount:      remainingAmount,
			Progress:             progress,
			DaysRemaining:        daysRemaining,
			MonthlySavingsNeeded: monthlyNeeded,
			Status:               status,
		})

	}

	c.JSON(http.StatusOK, response)

}

func UpdateSavingsGoal(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var body models.SavingGoalRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	deadline, err := time.Parse("2006-01-02", body.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use YYYY-MM-DD."})
		return
	}

	var goal models.SavingsGoal
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Savings goal not found"})
		return
	}

	goal.Title = body.Title
	goal.TargetAmount = body.TargetAmount
	goal.Deadline = deadline

	if err := config.DB.Save(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update savings goal"})
		return
	}

	c.JSON(http.StatusOK, goal)
}

func DeleteSavingsGoal(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var goal models.SavingsGoal
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Savings goal not found"})
		return
	}

	config.DB.Where("goal_id = ?", goal.ID).Delete(&models.GoalContribution{})

	if err := config.DB.Delete(&goal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete savings goal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Savings goal deleted successfully"})
}
