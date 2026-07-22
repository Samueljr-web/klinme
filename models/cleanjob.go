package models

import "time"

type CleanJobStatus string

const (
	CleanJobStatusPending    CleanJobStatus = "pending"
	CleanJobStatusProcessing CleanJobStatus = "processing"
	CleanJobStatusDone       CleanJobStatus = "done"
	CleanJobStatusFailed     CleanJobStatus = "failed"
)

type CleanJob struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	FileID        string         `json:"file_id"`
	Status        CleanJobStatus `json:"status"`
	RowsProcessed int            `json:"rows_processed"`
	RowsCleaned   int            `json:"rows_cleaned"`
	ErrorMessage  string         `json:"error_message"`
	CreatedAt     time.Time      `json:"created_at"`
	CompletedAt   time.Time      `json:"completed_at"`
}
