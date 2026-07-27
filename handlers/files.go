package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Samueljr-web/klinme-api/db"
	"github.com/Samueljr-web/klinme-api/models"
	"github.com/Samueljr-web/klinme-api/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	// Limit file size to 10MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	// Get the file from request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No file uploaded",
		})
		return
	}

	// Validate file type
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".csv" && ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only CSV and Excel files are allowed",
		})
		return
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open file",
		})
		return
	}
	defer file.Close()

	// Generate unique file name
	uniqueFileName := fmt.Sprintf("%s-%s", uuid.New().String(), fileHeader.Filename)

	// Upload to Azure raw container
	containerName := "raw"
	fileURL, err := storage.UploadFile(
		context.Background(),
		containerName,
		uniqueFileName,
		file,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload file to storage",
		})
		return
	}

	// Get userID from auth middleware
	userID := c.GetString("userID")

	// Save file metadata to Postgres
	fileID := uuid.New().String()
	_, err = db.Conn.Exec(
		context.Background(),
		`INSERT INTO files (id, user_id, file_name, file_size, file_type, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		fileID,
		userID,
		fileHeader.Filename,
		fileHeader.Size,
		ext,
		models.FileStatusPending,
		time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file record",
		})
		return
	}

	// Create a CleanJob record
	_, err = db.Conn.Exec(
		context.Background(),
		`INSERT INTO clean_jobs (id, user_id, file_id, status, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		uuid.New().String(),
		userID,
		fileID,
		models.CleanJobStatusPending,
		time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create clean job",
		})
		return
	}

	//  Return success
	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_id":  fileID,
		"file_url": fileURL,
	})
}
