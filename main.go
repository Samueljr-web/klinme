package main

import (
	"log"

	"github.com/Samueljr-web/klinme-api/db"
	"github.com/Samueljr-web/klinme-api/routes"
	"github.com/Samueljr-web/klinme-api/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to  DB
	db.Connect()
	defer db.Close()

	//Connect to azure
	storage.Connect()

	r := gin.Default()
	routes.SetupRoutes(r)
	r.Run(":8080")
}
