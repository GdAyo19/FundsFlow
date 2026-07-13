package main

import (
	"log"

	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/GdAyo19/FundsFlow/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	// establish a connection to the database using GORM
	config.ConnectDataBase()

	// automatically migrate the database schema for the User, Transaction, and Budget models
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Transaction{},
		&models.Budget{},
		&models.SavingsGoal{},
		&models.GoalContribution{},
	); err != nil {
		log.Fatal("Failed to migrate database schema: ", err)
	}

	// create a new Gin router instance
	router := gin.Default()

	// set up the routes for the API endpoints using the SetUpRoutes function from the routes package
	routes.SetUpRoutes(router)

	router.Run(":8080")

}
