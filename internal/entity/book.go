package entity

import (
	"time"
)

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	AuthorID    int       `json:"author_id"`
	Author      *Author   `json:"author,omitempty"` // Для JOIN операций
	ISBN        string    `json:"isbn"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BookResponse struct {
	Book
	AverageRating float64 `json:"average_rating,omitempty"`
	UserRating    *int    `json:"user_rating,omitempty"`
	TotalRatings  int     `json:"total_ratings,omitempty"`
}
