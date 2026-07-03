package models

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Transaction struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	Type        string    `gorm:"size:20;not null"`
	Amount      float64   `gorm:"not null"`
	Category    string    `gorm:"size:100"`
	Description string    `gorm:"size:255"`
	Date        time.Time `gorm:"not null"`
	CreatedAt   time.Time
}
