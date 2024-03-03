package domain

import (
	"time"
)

type File struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	OriginalName string    `json:"original_name"`
	Location     string    `json:"location"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
	UserId       int       `json:"user_id"`
}
