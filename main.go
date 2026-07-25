package main

import (
	"log"
	"net/http"

	"github.com/Samueljr-web/klinme-api/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Connect to  DB
      db.Connect()
	  defer db.Close()

    // Gin router
    r := gin.Default()

    r.GET("/api", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "welcome",
        })
    })

    r.Run(":8080")
}