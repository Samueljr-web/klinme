package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Samueljr-web/klinme-api/cleaner"
	"github.com/Samueljr-web/klinme-api/db"
	"github.com/Samueljr-web/klinme-api/models"
	"github.com/Samueljr-web/klinme-api/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	//Limit file size to 10MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)

	//Get the file from request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	//nValidate file type
	ext := filepath.Ext(fileHeader.Filename)
	if ext != ".csv" && ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV and Excel files are allowed"})
		return
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	// Parse CSV
	records, headers, err := cleaner.ParseCSV(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse CSV file"})
		return
	}

	// Run cleaning pipeline
	result, err := cleaner.RunPipeline(records, headers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clean file"})
		return
	}

	// Write cleaned records to a buffer
	var cleanedBuf bytes.Buffer
	if err := cleaner.WriteCSV(&cleanedBuf, result.Records, result.Headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write cleaned file"})
		return
	}

	// Upload raw file to Azure
	file.Seek(0, 0) // reset file pointer to beginning
	rawFileName := fmt.Sprintf("%s-%s", uuid.New().String(), fileHeader.Filename)
	rawURL, err := storage.UploadFile(
		context.Background(),
		"raw",
		rawFileName,
		file,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload raw file"})
		return
	}

	//Upload cleaned file to Azure
	cleanedFileName := fmt.Sprintf("cleaned-%s", rawFileName)
	cleanedURL, err := storage.UploadFile(
		context.Background(),
		"cleaned",
		cleanedFileName,
		&cleanedBuf,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload cleaned file"})
		return
	}

	//  Get userID from middleware
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
		models.FileStatusDone,
		time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file record"})
		return
	}

	// Save clean job to Postgres
	_, err = db.Conn.Exec(
		context.Background(),
		`INSERT INTO clean_jobs (id, user_id, file_id, status, rows_processed, rows_cleaned, created_at, completed_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		uuid.New().String(),
		userID,
		fileID,
		models.CleanJobStatusDone,
		result.RowsIn,
		result.RowsCleaned,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save clean job"})
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{
		"message":      "File cleaned successfully",
		"file_id":      fileID,
		"raw_url":      rawURL,
		"cleaned_url":  cleanedURL,
		"rows_in":      result.RowsIn,
		"rows_out":     result.RowsOut,
		"rows_cleaned": result.RowsCleaned,
	})
}
