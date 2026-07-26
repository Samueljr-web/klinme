package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/Samueljr-web/klinme-api/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// func CreateUser(c *gin.Context) {
// 	//sql
// 	InsertUser := `INSERT INTO users (id, clerk_user_id, email, plan, cleans_used, created_at)
//      VALUES ($1, $2, $3, $4, $5, $6)`

// 	ctx := c.Request.Context()

// 	var body struct {
// 		ClerkUserID string `json:"clerk_user_id"`
// 		Email       string `json:"email"`
// 	}

// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid request body",
// 		})
// 		return
// 	}

// 	if body.ClerkUserID == "" || body.Email == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "clerk_user_id and email are required",
// 		})
// 		return
// 	}

// 	// Insert user into db
// 	userID := uuid.New().String()
// 	_, err := db.Conn.Exec(ctx, InsertUser,
// 		userID,
// 		body.ClerkUserID,
// 		body.Email,
// 		"free",
// 		0,
// 		time.Now(),
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Failed to create user",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "User created",
// 		"user_id": userID,
// 	})

// }

func GetUser(c *gin.Context) {
	//sql query
	SelectUser := `SELECT clerk_user_id, email, plan, cleans_used, created_at FROM users WHERE id = $1`

	ctx := c.Request.Context()

	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user id is required",
		})
		return
	}

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id format",
		})
		return
	}

	var clerkUserID, email, plan string
	var cleansUsed int
	var createdAt time.Time

	err = db.Conn.QueryRow(ctx, SelectUser, id).Scan(&clerkUserID, &email, &plan, &cleansUsed, &createdAt)

	if errors.Is(err, pgx.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "database error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            id,
		"clerk_user_id": clerkUserID,
		"email":         email,
		"plan":          plan,
		"cleans_used":   cleansUsed,
		"created_at":    createdAt,
	})

}
