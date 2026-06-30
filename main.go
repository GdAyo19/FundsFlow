package main

import (
	"github.com/GdAyo19/FundsFlow/config"
	"github.com/GdAyo19/FundsFlow/models"
	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDataBase()

	config.DB.AutoMigrate(
		&models.User{},
	)

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Funds Flow API Running 🚀",
		})
	})

	router.Run(":8080")

}
