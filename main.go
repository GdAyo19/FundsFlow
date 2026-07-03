package main

import (
	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/GdAyo19/FundsFlow/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDataBase()

	config.DB.AutoMigrate(
		&models.User{},
		&models.Transaction{},
	)

	router := gin.Default()

	routes.SetUpRoutes(router)

	router.Run(":8080")

}
