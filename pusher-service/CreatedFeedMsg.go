package main

import "time"

type CreatedFeedMsg struct {
	Type        string    `json:"type"`
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewCreatedFeedMsg(id string, title string, description string, createdAt time.Time) *CreatedFeedMsg {
	return &CreatedFeedMsg{
		Type:        "created_feed",
		ID:          id,
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
	}
}
