package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Samueljr-web/klinme-api/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClerkWebhookEvent struct {
	Type string `json:"type"`
	Data struct {
		ID             string `json:"id"`
		EmailAddresses []struct {
			EmailAddress string `json:"email_address"`
		} `json:"email_addresses"`
	} `json:"data"`
}

func ClerkWebhook(c *gin.Context) {
	ctx := context.Background()
	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Parse the event
	var event ClerkWebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse event"})
		return
	}

	// Handle user.created event
	if event.Type == "user.created" {
		clerkUserID := event.Data.ID
		email := event.Data.EmailAddresses[0].EmailAddress

		// Check if user already exists
		var exists bool
		err := db.Conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE clerk_user_id = $1)`, clerkUserID).Scan(&exists)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		if exists {
			c.JSON(http.StatusOK, gin.H{"message": "User already exists"})
			return
		}

		// Create user in Postgres
		_, err = db.Conn.Exec(ctx,
			`INSERT INTO users (id, clerk_user_id, email, plan, cleans_used, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			uuid.New().String(),
			clerkUserID,
			email,
			"free",
			0,
			time.Now(),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		fmt.Printf("User created: %s\n", email)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook received"})
}
