package controllers

import (
	"net/http"
	"strconv"

	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/gin-gonic/gin"
)

func AddContribution(c *gin.Context) {

	var body models.GoalContributionRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.GetUint("userID")

	// 
	goalID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": " Invalid goal ID",
		})
		return
	}

	var goal models.SavingsGoal

	if err := config.DB.Where("id = ? AND user_id = ? ", goalID, userID).First(&goal).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "goal not found",
		})
		return
	}

	contribution := models.GoalContribution{
		GoalID: uint(goalID),
		UserID: userID,
		Amount: body.Amount,
		Note:   body.Note,
	}

	if err := config.DB.Create(&contribution).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create contribution",
		})
		return
	}

	c.JSON(http.StatusOK, contribution)
}

func DeleteContribution(c *gin.Context) {
	userID := c.GetUint("userID")
	contributionID := c.Param("contributionId")

	var contribution models.GoalContribution
	if err := config.DB.Where("id = ? AND user_id = ?", contributionID, userID).First(&contribution).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contribution not found"})
		return
	}

	if err := config.DB.Delete(&contribution).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contribution"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contribution deleted successfully"})
}


