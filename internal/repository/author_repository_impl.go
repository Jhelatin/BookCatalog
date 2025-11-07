package postgres

import (
	"awesomeProject/internal/entity"
	"context"
	"database/sql"
)

type authorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) AuthorRepository {
	return &authorRepository{db: db}
}

func (r *authorRepository) Create(ctx context.Context, author *entity.Author) error {
	query := `
        INSERT INTO authors (first_name, last_name, biography, birth_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `

	return r.db.QueryRowContext(ctx, query,
		author.FirstName, author.LastName, author.Biography,
		author.BirthDate, author.CreatedAt, author.UpdatedAt,
	).Scan(&author.ID)
}

func (r *authorRepository) FindByID(ctx context.Context, id int) (*entity.Author, error) {
	query := `
        SELECT id, first_name, last_name, biography, birth_date, created_at, updated_at
        FROM authors WHERE id = $1
    `

	author := &entity.Author{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&author.ID, &author.FirstName, &author.LastName, &author.Biography,
		&author.BirthDate, &author.CreatedAt, &author.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return author, err
}

func (r *authorRepository) FindAll(ctx context.Context) ([]*entity.Author, error) {
	query := `
        SELECT id, first_name, last_name, biography, birth_date, created_at, updated_at
        FROM authors ORDER BY last_name, first_name
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []*entity.Author
	for rows.Next() {
		author := &entity.Author{}
		err := rows.Scan(
			&author.ID, &author.FirstName, &author.LastName, &author.Biography,
			&author.BirthDate, &author.CreatedAt, &author.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}

func (r *authorRepository) Update(ctx context.Context, author *entity.Author) error {
	query := `
        UPDATE authors 
        SET first_name = $1, last_name = $2, biography = $3, birth_date = $4, updated_at = $5
        WHERE id = $6
    `

	_, err := r.db.ExecContext(ctx, query,
		author.FirstName, author.LastName, author.Biography,
		author.BirthDate, author.UpdatedAt, author.ID,
	)
	return err
}

func (r *authorRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM authors WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
