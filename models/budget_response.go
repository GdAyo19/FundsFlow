package models

type BudgetResponse struct {
	ID        uint    `json:"id"`
	Category  string  `json:"category"`
	Budget    float64 `json:"budget"`
	Spent     float64 `json:"spent"`
	Remaining float64 `json:"remaining"`
	Progress  float64 `json:"progress"`
	Status    string  `json:"status"`
	Month     int     `json:"month"`
	Year      int     `json:"year"`
}
