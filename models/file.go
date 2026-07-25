package models

import "time"

type FileStatus string

const (
	FileStatusPending    FileStatus = "pending"
	FileStatusProcessing FileStatus = "processing"
	FileStatusDone       FileStatus = "done"
	FileStatusFailed     FileStatus = "failed"
)

type File struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	FileName  string     `json:"file_name"`
	FileSize  int64      `json:"file_size"`
	FileType  string     `json:"file_type"`
	Status    FileStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}
