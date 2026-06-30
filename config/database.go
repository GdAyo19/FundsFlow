package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// creating a data base connection 
func ConnectDataBase() {
	dns := "host=localhost user=fundsflow password=password dbname=fundsflow port=5433 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		panic("DataBase Connection failed")
	}

	fmt.Println("DataBase Connection established")
	DB = db
}
