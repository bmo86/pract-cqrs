package models

import "time"

type Feed struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	CreatedAt   time.Time `json:"cread_at"`
}
