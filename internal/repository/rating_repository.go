package postgres

import (
	"awesomeProject/internal/entity"
	"context"
)

type RatingRepository interface {
	CreateOrUpdate(ctx context.Context, rating *entity.Rating) error
	GetByUserAndBook(ctx context.Context, userID, bookID int) (*entity.Rating, error)
	GetByBookID(ctx context.Context, bookID int) ([]*entity.Rating, error)
	GetByUserID(ctx context.Context, userID int) ([]*entity.Rating, error)
	GetAverageRating(ctx context.Context, bookID int) (float64, int, error)
	Delete(ctx context.Context, userID, bookID int) error
	GetBooksWithUserRatings(ctx context.Context, userID int) ([]*entity.BookWithRating, error)
}
