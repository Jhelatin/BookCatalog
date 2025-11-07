package service

import (
	"awesomeProject/internal/entity"
	postgres "awesomeProject/internal/repository"
	"context"
	"time"
)

type BookService struct {
	bookRepo postgres.BookRepository
}

func NewBookService(bookRepo postgres.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (s *BookService) CreateBook(ctx context.Context, book *entity.Book) error {
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	return s.bookRepo.Create(ctx, book)
}

func (s *BookService) GetBook(ctx context.Context, id int) (*entity.Book, error) {

	return s.bookRepo.FindByID(ctx, id)
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]*entity.Book, error) {
	return s.bookRepo.FindAll(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, book *entity.Book) error {
	book.UpdatedAt = time.Now()
	return s.bookRepo.Update(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	return s.bookRepo.Delete(ctx, id)
}
