package postgres

import (
	"awesomeProject/internal/entity"
	"context"
	"database/sql"
)

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, book *entity.Book) error {
	query := `
        INSERT INTO books (title, author, isbn, published_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	return r.db.QueryRowContext(ctx, query,
		book.Title, book.Author, book.ISBN, book.PublishedAt, book.CreatedAt, book.UpdatedAt,
	).Scan(&book.ID)
}

func (r *bookRepository) FindByID(ctx context.Context, id int) (*entity.Book, error) {
	query := `
        SELECT id, title, author, isbn, published_at, created_at, updated_at
        FROM books WHERE id = $1
    `

	book := &entity.Book{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&book.ID, &book.Title, &book.Author, &book.ISBN,
		&book.PublishedAt, &book.CreatedAt, &book.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return book, err
}

func (r *bookRepository) FindAll(ctx context.Context) ([]*entity.Book, error) {
	query := `
        SELECT id, title, author, isbn, published_at, created_at, updated_at
        FROM books ORDER BY id
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*entity.Book
	for rows.Next() {
		book := &entity.Book{}
		err := rows.Scan(
			&book.ID, &book.Title, &book.Author, &book.ISBN,
			&book.PublishedAt, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (r *bookRepository) Update(ctx context.Context, book *entity.Book) error {
	query := `
        UPDATE books 
        SET title = $1, author = $2, isbn = $3, published_at = $4, updated_at = $5
        WHERE id = $6
    `

	_, err := r.db.ExecContext(ctx, query,
		book.Title, book.Author, book.ISBN, book.PublishedAt, book.UpdatedAt, book.ID,
	)
	return err
}

func (r *bookRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
