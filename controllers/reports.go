package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/gin-gonic/gin"
)

type MonthlyReportResponse struct {
	Month        int              `json:"month"`
	Year         int              `json:"year"`
	TotalIncome  float64          `json:"total_income"`
	TotalExpense float64          `json:"total_expense"`
	NetSavings   float64          `json:"net_savings"`
	Categories   []CategoryReport `json:"categories"`
}

type CategoryReport struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
	Count    int64   `json:"count"`
}

func MonthlyReport(c *gin.Context) {
	userID := c.GetUint("userID")
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	if monthStr != "" {
		fmt.Sscanf(monthStr, "%d", &month)
	}
	if yearStr != "" {
		fmt.Sscanf(yearStr, "%d", &year)
	}

	var totalIncome, totalExpense float64
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
			userID, "income", month, year).
		Select("COALESCE(SUM(amount),0)").Scan(&totalIncome)

	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
			userID, "expense", month, year).
		Select("COALESCE(SUM(amount),0)").Scan(&totalExpense)

	var expenseCategories []CategoryReport
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
			userID, "expense", month, year).
		Select("category, COALESCE(SUM(amount),0) as total, COUNT(*) as count").
		Group("category").
		Scan(&expenseCategories)

	c.JSON(http.StatusOK, MonthlyReportResponse{
		Month:        month,
		Year:         year,
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		NetSavings:   totalIncome - totalExpense,
		Categories:   expenseCategories,
	})
}

func CategoryReportHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	category := c.Query("category")
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	if monthStr != "" {
		fmt.Sscanf(monthStr, "%d", &month)
	}
	if yearStr != "" {
		fmt.Sscanf(yearStr, "%d", &year)
	}

	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category query parameter is required"})
		return
	}

	var transactions []models.Transaction
	config.DB.Where("user_id = ? AND category = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
		userID, category, month, year).Find(&transactions)

	var total float64
	for _, t := range transactions {
		total += t.Amount
	}

	c.JSON(http.StatusOK, gin.H{
		"category":     category,
		"month":        month,
		"year":         year,
		"total":        total,
		"count":        len(transactions),
		"transactions": transactions,
	})
}

func ExportCSV(c *gin.Context) {
	userID := c.GetUint("userID")

	var transactions []models.Transaction
	config.DB.Where("user_id = ?", userID).Order("date desc").Find(&transactions)

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=transactions.csv")

	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"ID", "Type", "Amount", "Category", "Description", "Date", "CreatedAt"})

	for _, t := range transactions {
		writer.Write([]string{
			fmt.Sprintf("%d", t.ID),
			t.Type,
			fmt.Sprintf("%.2f", t.Amount),
			t.Category,
			t.Description,
			t.Date.Format("2006-01-02"),
			t.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	writer.Flush()
}

func ExportPDF(c *gin.Context) {
	userID := c.GetUint("userID")

	var transactions []models.Transaction
	config.DB.Where("user_id = ?", userID).Order("date desc").Find(&transactions)

	html := `<html><head><style>
		table { width: 100%%; border-collapse: collapse; }
		th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
		th { background-color: #4CAF50; color: white; }
	</style></head><body>
	<h1>Transaction Report</h1>
	<table>
		<tr><th>ID</th><th>Type</th><th>Amount</th><th>Category</th><th>Description</th><th>Date</th></tr>`

	for _, t := range transactions {
		html += fmt.Sprintf(
			"<tr><td>%d</td><td>%s</td><td>%.2f</td><td>%s</td><td>%s</td><td>%s</td></tr>",
			t.ID, t.Type, t.Amount, t.Category, t.Description, t.Date.Format("2006-01-02"),
		)
	}

	html += `</table></body></html>`

	c.Header("Content-Type", "text/html")
	c.Header("Content-Disposition", "attachment; filename=transactions.html")
	c.String(http.StatusOK, html)
}
