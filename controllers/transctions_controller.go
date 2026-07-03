package controllers

import (
	"net/http"
	"time"

	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {

	var body struct {
		Type        string  `json:"type"`
		Amount      float64 `json:"amount"`
		Category    string  `json:"category"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body request",
		})
		return
	}

	// get the user ID from the context set by the AuthMiddleware
	userID := c.GetUint("userID")
	// create a new transaction using the data from the request body and the user ID
	transaction := models.Transaction{
		UserID:      userID,
		Type:        body.Type,
		Amount:      body.Amount,
		Category:    body.Category,
		Description: body.Description,
		Date:        time.Now(),
	}
	// save the transaction to the database using GORM
	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create transaction",
		})
		return
	}

	c.JSON(http.StatusCreated, transaction)

}

func GetTransactions(c *gin.Context) {
	// get the user ID from the context set by the AuthMiddleware
	userID := c.GetUint("userID")

	var transactions []models.Transaction

	// retrieve all transactions for the user from the database using GORM
	result := config.DB.Where("user_id = ?", userID).Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve transactions",
		})
		return
	}
	// return the transactions as a JSON response
	c.JSON(http.StatusOK, transactions)
}

func UpdateTransaction(c *gin.Context) {
	var body models.UpdateTransactionRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid body request",
		})
		return
	}

	userID := c.GetUint("userID")

	id := c.Param("id")

	var transaction models.Transaction

	// retrieve the transaction from the database using GORM
	result := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&transaction)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Transaction not found",
		})
		return
	}

	if transaction.UserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "You are not authorized to update this transaction",
		})
		return
	}

	// update the transaction fields with the data from the request body
	transaction.Type = body.Type
	transaction.Amount = body.Amount
	transaction.Category = body.Category
	transaction.Description = body.Description

	if err := config.DB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update transaction",
		})
		return
	}
	// return the updated transaction as a JSON response
	c.JSON(http.StatusOK, transaction)
}

func DeleteTransaction(c *gin.Context) {
	// get the transaction ID from the URL parameter and the user ID from the context
	id := c.Param("id")
	userID := c.GetUint("userID")

	var transaction models.Transaction

	// retrieve the transaction from the database using GORM
	err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Transaction not found",
		})
		return
	}

	if err := config.DB.Delete(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete transaction",
		})
		return
	}

	// return a success message as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction deleted successfully",
	})

}
