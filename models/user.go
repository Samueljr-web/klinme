package models

import "time"

type User struct{
	ID string `json:"id"`
	ClerkUserID string `json:"clerk_user_id"`
	Email string `json:"email"`
	Plan string `json:"plan"`
	CleansUsed int `json:"cleans_used"`
	CreatedAt time.Time `json:"created_at"`
}