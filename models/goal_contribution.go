package models

import "time"

type GoalContribution struct {
	ID        uint    `gorm:"primaryKey"`
	GoalID    uint    `gorm:"not null"`
	UserID    uint    `gorm:"not null"`
	Amount    float64 `gorm:"not null"`
	Note      string  `gorm:"size:255"`
	CreatedAt time.Time
}
