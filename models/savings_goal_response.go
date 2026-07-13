package models

type SavingsGoalResponse struct {
	ID                   uint    `json:"id"`
	Title                string  `json:"title"`
	TargetAmount         float64 `json:"target_amount"`
	SavedAmount          float64 `json:"saved_amount"`
	RemainingAmount      float64 `json:"remaining_amount"`
	Progress             float64 `json:"progress"`
	DaysRemaining        int     `json:"days_remaining"`
	MonthlySavingsNeeded float64 `json:"monthly_savings_needed"`
	Status               string  `json:"status"`
}


