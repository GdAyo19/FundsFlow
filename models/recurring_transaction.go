package models

import "time"

type RecurringTransaction struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	Type        string    `gorm:"size:20;not null"`
	Amount      float64   `gorm:"not null"`
	Category    string    `gorm:"size:100"`
	Description string    `gorm:"size:255"`
	Frequency   string    `gorm:"size:20;not null"`
	NextDate    time.Time `gorm:"not null"`
	EndDate     *time.Time
	IsActive    bool `gorm:"default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
