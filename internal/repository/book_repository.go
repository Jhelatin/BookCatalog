package postgres

import (
	"awesomeProject/internal/entity"
	"context"
)

type BookRepository interface {
	Create(ctx context.Context, book *entity.Book) error
	FindByID(ctx context.Context, id int) (*entity.Book, error)
	FindAll(ctx context.Context) ([]*entity.Book, error)
	Update(ctx context.Context, book *entity.Book) error
	Delete(ctx context.Context, id int) error
}
