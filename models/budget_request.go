package models

type BudgetRequest struct {
	Category string  `json:"category" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Month    int     `json:"month" binding:"required,min=1,max=12"`
	Year     int     `json:"year" binding:"required"`
}	

