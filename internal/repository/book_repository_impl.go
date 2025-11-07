package postgres

import (
	"awesomeProject/internal/entity"
	"context"
	"database/sql"
)

type bookRepository struct {
	db *sql.DB
}

// NewBookRepository creates a new instance of BookRepository
func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, book *entity.Book) error {
	query := `
        INSERT INTO books (title, author_id, isbn, description, published_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `

	return r.db.QueryRowContext(ctx, query,
		book.Title, book.AuthorID, book.ISBN, book.Description,
		book.PublishedAt, book.CreatedAt, book.UpdatedAt,
	).Scan(&book.ID)
}

func (r *bookRepository) FindByID(ctx context.Context, id int) (*entity.Book, error) {
	query := `
        SELECT b.id, b.title, b.author_id, b.isbn, b.description, b.published_at, 
               b.created_at, b.updated_at,
               a.id, a.first_name, a.last_name, a.biography, a.birth_date,
               a.created_at, a.updated_at
        FROM books b
        LEFT JOIN authors a ON b.author_id = a.id
        WHERE b.id = $1
    `

	book := &entity.Book{Author: &entity.Author{}}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&book.ID, &book.Title, &book.AuthorID, &book.ISBN, &book.Description,
		&book.PublishedAt, &book.CreatedAt, &book.UpdatedAt,
		&book.Author.ID, &book.Author.FirstName, &book.Author.LastName,
		&book.Author.Biography, &book.Author.BirthDate,
		&book.Author.CreatedAt, &book.Author.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return book, err
}

func (r *bookRepository) FindAll(ctx context.Context) ([]*entity.Book, error) {
	query := `
        SELECT b.id, b.title, b.author_id, b.isbn, b.description, b.published_at, 
               b.created_at, b.updated_at,
               a.id, a.first_name, a.last_name, a.biography, a.birth_date,
               a.created_at, a.updated_at
        FROM books b
        LEFT JOIN authors a ON b.author_id = a.id
        ORDER BY b.title
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*entity.Book
	for rows.Next() {
		book := &entity.Book{Author: &entity.Author{}}
		err := rows.Scan(
			&book.ID, &book.Title, &book.AuthorID, &book.ISBN, &book.Description,
			&book.PublishedAt, &book.CreatedAt, &book.UpdatedAt,
			&book.Author.ID, &book.Author.FirstName, &book.Author.LastName,
			&book.Author.Biography, &book.Author.BirthDate,
			&book.Author.CreatedAt, &book.Author.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (r *bookRepository) FindByAuthorID(ctx context.Context, authorID int) ([]*entity.Book, error) {
	query := `
        SELECT b.id, b.title, b.author_id, b.isbn, b.description, b.published_at, 
               b.created_at, b.updated_at,
               a.id, a.first_name, a.last_name, a.biography, a.birth_date,
               a.created_at, a.updated_at
        FROM books b
        LEFT JOIN authors a ON b.author_id = a.id
        WHERE b.author_id = $1
        ORDER BY b.title
    `

	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*entity.Book
	for rows.Next() {
		book := &entity.Book{Author: &entity.Author{}}
		err := rows.Scan(
			&book.ID, &book.Title, &book.AuthorID, &book.ISBN, &book.Description,
			&book.PublishedAt, &book.CreatedAt, &book.UpdatedAt,
			&book.Author.ID, &book.Author.FirstName, &book.Author.LastName,
			&book.Author.Biography, &book.Author.BirthDate,
			&book.Author.CreatedAt, &book.Author.UpdatedAt,
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
        SET title = $1, author_id = $2, isbn = $3, description = $4, 
            published_at = $5, updated_at = $6
        WHERE id = $7
    `

	_, err := r.db.ExecContext(ctx, query,
		book.Title, book.AuthorID, book.ISBN, book.Description,
		book.PublishedAt, book.UpdatedAt, book.ID,
	)
	return err
}

func (r *bookRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
