package models

type SavingGoalRequest struct {
	Title        string  `json:"title" binding:"required"`
	TargetAmount float64 `json:"target_amount" binding:"required,gt=0"`
	Deadline     string  `json:"deadline" binding:"required"`
}
