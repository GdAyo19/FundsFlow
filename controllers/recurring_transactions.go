package controllers

import (
	"net/http"
	"time"

	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/gin-gonic/gin"
)

func CreateRecurringTransaction(c *gin.Context) {
	userID := c.GetUint("userID")

	var body struct {
		Type        string `json:"type" binding:"required,oneof=income expense"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Category    string `json:"category" binding:"required"`
		Description string `json:"description"`
		Frequency   string `json:"frequency" binding:"required,oneof=daily weekly monthly"`
		StartDate   string `json:"start_date" binding:"required"`
		EndDate     string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD."})
		return
	}

	var endDate *time.Time
	if body.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", body.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD."})
			return
		}
		endDate = &parsed
	}

	recurring := models.RecurringTransaction{
		UserID:      userID,
		Type:        body.Type,
		Amount:      body.Amount,
		Category:    body.Category,
		Description: body.Description,
		Frequency:   body.Frequency,
		NextDate:    startDate,
		EndDate:     endDate,
		IsActive:    true,
	}

	if err := config.DB.Create(&recurring).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recurring transaction"})
		return
	}

	c.JSON(http.StatusCreated, recurring)
}

func GetRecurringTransactions(c *gin.Context) {
	userID := c.GetUint("userID")

	var recurring []models.RecurringTransaction
	config.DB.Where("user_id = ?", userID).Order("next_date asc").Find(&recurring)

	c.JSON(http.StatusOK, recurring)
}

func UpdateRecurringTransaction(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var body struct {
		Type        string  `json:"type" binding:"required,oneof=income expense"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Category    string  `json:"category" binding:"required"`
		Description string  `json:"description"`
		Frequency   string  `json:"frequency" binding:"required,oneof=daily weekly monthly"`
		IsActive    *bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var recurring models.RecurringTransaction
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&recurring).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recurring transaction not found"})
		return
	}

	recurring.Type = body.Type
	recurring.Amount = body.Amount
	recurring.Category = body.Category
	recurring.Description = body.Description
	recurring.Frequency = body.Frequency
	if body.IsActive != nil {
		recurring.IsActive = *body.IsActive
	}

	if err := config.DB.Save(&recurring).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recurring transaction"})
		return
	}

	c.JSON(http.StatusOK, recurring)
}

func DeleteRecurringTransaction(c *gin.Context) {
	userID := c.GetUint("userID")
	id := c.Param("id")

	var recurring models.RecurringTransaction
	if err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&recurring).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recurring transaction not found"})
		return
	}

	if err := config.DB.Delete(&recurring).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recurring transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recurring transaction deleted successfully"})
}
