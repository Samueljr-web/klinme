package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Connect to  DB
    dbURL := os.Getenv("DATABASE_URL")
    conn, err := pgx.Connect(context.Background(), dbURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer conn.Close(context.Background())

    fmt.Println("Connected to database!")

    // Gin router
    r := gin.Default()

    r.GET("/api", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "welcome",
        })
    })

    r.Run(":8080")
}