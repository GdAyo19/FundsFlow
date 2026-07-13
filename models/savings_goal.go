package models

import "time"

type SavingsGoal struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"not null"`
	Title        string    `gorm:"size:150;not null"`
	TargetAmount  float64   `gorm:"not null"`
	Deadline      time.Time `gorm:"not null"`
	Created_at    time.Time
	Updated_at    time.Time
}
