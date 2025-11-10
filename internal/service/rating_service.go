package service

import (
	"awesomeProject/internal/entity"
	postgres "awesomeProject/internal/repository"
	"context"
	"errors"
)

type RatingService struct {
	ratingRepo postgres.RatingRepository
	bookRepo   postgres.BookRepository
}

func NewRatingService(ratingRepo postgres.RatingRepository, bookRepo postgres.BookRepository) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
		bookRepo:   bookRepo,
	}
}

func (s *RatingService) RateBook(ctx context.Context, userID int, req *entity.RatingRequest) error {
	book, err := s.bookRepo.FindByID(ctx, req.BookID)
	if err != nil {
		return err
	}
	if book == nil {
		return errors.New("book not found")
	}

	if req.Score < 0 || req.Score > 5 {
		return errors.New("score must be between 0 and 5")
	}

	rating := &entity.Rating{
		UserID: userID,
		BookID: req.BookID,
		Score:  req.Score,
	}

	return s.ratingRepo.CreateOrUpdate(ctx, rating)
}

func (s *RatingService) GetUserRating(ctx context.Context, userID, bookID int) (*entity.Rating, error) {
	return s.ratingRepo.GetByUserAndBook(ctx, userID, bookID)
}

func (s *RatingService) GetBookRatings(ctx context.Context, bookID int) ([]*entity.Rating, error) {
	return s.ratingRepo.GetByBookID(ctx, bookID)
}

func (s *RatingService) GetUserRatings(ctx context.Context, userID int) ([]*entity.Rating, error) {
	return s.ratingRepo.GetByUserID(ctx, userID)
}

func (s *RatingService) GetBookAverageRating(ctx context.Context, bookID int) (float64, int, error) {
	return s.ratingRepo.GetAverageRating(ctx, bookID)
}

func (s *RatingService) RemoveRating(ctx context.Context, userID, bookID int) error {
	return s.ratingRepo.Delete(ctx, userID, bookID)
}

func (s *RatingService) GetBooksWithUserRatings(ctx context.Context, userID int) ([]*entity.BookWithRating, error) {
	return s.ratingRepo.GetBooksWithUserRatings(ctx, userID)
}

func (s *RatingService) EnrichBookWithRatings(ctx context.Context, book *entity.Book, userID int) (*entity.BookResponse, error) {
	averageRating, totalRatings, err := s.GetBookAverageRating(ctx, book.ID)
	if err != nil {
		return nil, err
	}

	response := &entity.BookResponse{
		Book:          *book,
		AverageRating: averageRating,
		TotalRatings:  totalRatings,
	}

	if userID > 0 {
		userRating, err := s.GetUserRating(ctx, userID, book.ID)
		if err != nil {
			return nil, err
		}
		if userRating != nil {
			response.UserRating = &userRating.Score
		}
	}

	return response, nil
}

func (s *RatingService) EnrichBooksWithRatings(ctx context.Context, books []*entity.Book, userID int) ([]*entity.BookResponse, error) {
	var enrichedBooks []*entity.BookResponse

	for _, book := range books {
		enrichedBook, err := s.EnrichBookWithRatings(ctx, book, userID)
		if err != nil {
			return nil, err
		}
		enrichedBooks = append(enrichedBooks, enrichedBook)
	}

	return enrichedBooks, nil
}
