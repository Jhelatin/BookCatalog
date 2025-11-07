package service

import (
	"awesomeProject/internal/entity"
	postgres "awesomeProject/internal/repository"
	"context"
	"time"
)

type AuthorService struct {
	authorRepo postgres.AuthorRepository
	bookRepo   postgres.BookRepository
}

// NewAuthorService creates a new AuthorService with dependency injection
func NewAuthorService(authorRepo postgres.AuthorRepository, bookRepo postgres.BookRepository) *AuthorService {
	return &AuthorService{
		authorRepo: authorRepo,
		bookRepo:   bookRepo,
	}
}

func (s *AuthorService) CreateAuthor(ctx context.Context, author *entity.Author) error {
	author.CreatedAt = time.Now()
	author.UpdatedAt = time.Now()
	return s.authorRepo.Create(ctx, author)
}

func (s *AuthorService) GetAuthor(ctx context.Context, id int) (*entity.Author, error) {
	return s.authorRepo.FindByID(ctx, id)
}

func (s *AuthorService) GetAllAuthors(ctx context.Context) ([]*entity.Author, error) {
	return s.authorRepo.FindAll(ctx)
}

func (s *AuthorService) UpdateAuthor(ctx context.Context, author *entity.Author) error {
	author.UpdatedAt = time.Now()
	return s.authorRepo.Update(ctx, author)
}

func (s *AuthorService) DeleteAuthor(ctx context.Context, id int) error {
	return s.authorRepo.Delete(ctx, id)
}

func (s *AuthorService) GetAuthorWithBooks(ctx context.Context, id int) (*entity.Author, []*entity.Book, error) {
	author, err := s.authorRepo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	if author == nil {
		return nil, nil, nil
	}

	books, err := s.bookRepo.FindByAuthorID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return author, books, nil
}
