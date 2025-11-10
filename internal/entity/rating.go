package entity

import (
	"time"
)

type Rating struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	BookID    int       `json:"book_id"`
	Score     int       `json:"score"` // 0-5
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RatingRequest struct {
	BookID int `json:"book_id" validate:"required,min=1"`
	Score  int `json:"score" validate:"required,min=0,max=5"`
}

type BookWithRating struct {
	Book
	AverageRating float64 `json:"average_rating"`
	UserRating    *int    `json:"user_rating,omitempty"`
	TotalRatings  int     `json:"total_ratings"`
}
