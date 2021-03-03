package models

import "time"

type Task struct {
	ID          string    `gorm:"primary_key;" json:"id"`
	Content     string    `gorm:"size:255;not null;" json:"content"`
	UserID      string    `gorm:"size:255;not null;" json:"user_id"`
	CreatedDate time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_date"`
}
