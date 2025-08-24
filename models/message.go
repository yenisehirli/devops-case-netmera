package models

import "time"

// HTTP request için
type MessageRequest struct {
	Content string `json:"content" binding:"required"`
	UserID  string `json:"user_id" binding:"required"`
}

// HTTP response için
type MessageResponse struct {
	ID      string    `json:"id"`
	Content string    `json:"content"`
	UserID  string    `json:"user_id"`
	Status  string    `json:"status"`
	Time    time.Time `json:"timestamp"`
}

// MongoDB için
type MessageDocument struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Content     string    `bson:"content" json:"content"`
	UserID      string    `bson:"user_id" json:"user_id"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	ProcessedAt time.Time `bson:"processed_at,omitempty" json:"processed_at,omitempty"`
}
