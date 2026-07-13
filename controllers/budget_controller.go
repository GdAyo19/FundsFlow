package controllers

import (
	"net/http"

	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/gin-gonic/gin"
)

func CreateBudget(c *gin.Context) {

	// bind the request body to the body struct
	var body models.BudgetRequest
	// bind the request body to the body struct
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	userID := c.GetUint("userID")
	// create a new budget using the data from the request body and the user ID
	budget := models.Budget{
		UserID:   userID,
		Category: body.Category,
		Amount:   body.Amount,
		Month:    body.Month,
		Year:     body.Year,
	}

	if err := config.DB.Create(&budget).Error; err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to create budget",
		})
		return
	}

	c.JSON(201, budget)
}

func GetBudgets(c *gin.Context) {
	// get the user ID from the context set by the AuthMiddleware
	userID := c.GetUint("userID")

	var budgets []models.Budget

	// retrieve all budgets for the user from the database using GORM
	result := config.DB.Where("user_id = ?", userID).Find(&budgets)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": "Failed to retrieve budgets",
		})
		return
	}

	var response []models.BudgetResponse

	for _, budget := range budgets {
		// calculate the total spent for the budget category in the specified month and year
		var totalSpent float64
		config.DB.Model(&models.Transaction{}).
			Where(
			`user_id = ? 
			AND category = ? 
			AND type = ? 
			AND EXTRACT(MONTH FROM date) = ? 
			AND EXTRACT(YEAR FROM date) = ?`, 
			userID, budget.Category, "expense", budget.Month, budget.Year).
			Select("COALESCE(SUM(amount),0)").Scan(&totalSpent)
		// calculate the remaining budget and progress percentage
		remaining := budget.Amount - totalSpent
		progress := (totalSpent / budget.Amount) * 100

		status := "On Track"

		if progress >= 100 {
			status = "Exceeded"
		} else if progress >= 75 {
			status = "At Risk"
		}

		response = append(response, models.BudgetResponse{
			ID:        budget.ID,
			Category:  budget.Category,
			Budget:    budget.Amount,
			Spent:     totalSpent,
			Remaining: remaining,
			Progress:  progress,
			Status:    status,
			Month:     budget.Month,
			Year:      budget.Year,
		})
	}

	c.JSON(http.StatusOK, response)

}

func UpdateBudget(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var body models.BudgetRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var budget models.Budget
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&budget).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	budget.Category = body.Category
	budget.Amount = body.Amount
	budget.Month = body.Month
	budget.Year = body.Year

	if err := config.DB.Save(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update budget"})
		return
	}

	c.JSON(http.StatusOK, budget)
}

func DeleteBudget(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var budget models.Budget
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&budget).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Budget not found"})
		return
	}

	if err := config.DB.Delete(&budget).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete budget"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Budget deleted successfully"})
}
