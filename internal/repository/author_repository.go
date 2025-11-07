package postgres

import (
	"awesomeProject/internal/entity"
	"context"
)

type AuthorRepository interface {
	Create(ctx context.Context, author *entity.Author) error
	FindByID(ctx context.Context, id int) (*entity.Author, error)
	FindAll(ctx context.Context) ([]*entity.Author, error)
	Update(ctx context.Context, author *entity.Author) error
	Delete(ctx context.Context, id int) error
}
