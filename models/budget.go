package models

import "time"

type Budget struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`
	Category  string  `gorm:"size:100;not null"`
	Amount    float64 `gorm:"not null"`
	Month     int     `gorm:"not null"`
	Year      int     `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
