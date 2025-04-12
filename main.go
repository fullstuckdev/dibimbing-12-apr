package main

import (
	"webroutes/config"
	"webroutes/models"
	"webroutes/routes"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading ENV")
	}

	r := gin.Default()
	db := config.ConnectDatabase()

	db.AutoMigrate(&models.User{})

	routes.SetupRoutes(r, db)


	r.Run(":8080")

}