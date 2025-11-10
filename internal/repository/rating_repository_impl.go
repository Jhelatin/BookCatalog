package postgres

import (
	"awesomeProject/internal/entity"
	"context"
	"database/sql"
	"time"
)

type ratingRepository struct {
	db *sql.DB
}

func NewRatingRepository(db *sql.DB) RatingRepository {
	return &ratingRepository{db: db}
}

func (r *ratingRepository) CreateOrUpdate(ctx context.Context, rating *entity.Rating) error {
	query := `
        INSERT INTO ratings (user_id, book_id, score, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id, book_id) 
        DO UPDATE SET score = $3, updated_at = $5
        RETURNING id
    `

	rating.CreatedAt = time.Now()
	rating.UpdatedAt = time.Now()

	return r.db.QueryRowContext(ctx, query,
		rating.UserID, rating.BookID, rating.Score, rating.CreatedAt, rating.UpdatedAt,
	).Scan(&rating.ID)
}

func (r *ratingRepository) GetByUserAndBook(ctx context.Context, userID, bookID int) (*entity.Rating, error) {
	query := `
        SELECT id, user_id, book_id, score, created_at, updated_at
        FROM ratings 
        WHERE user_id = $1 AND book_id = $2
    `

	rating := &entity.Rating{}
	err := r.db.QueryRowContext(ctx, query, userID, bookID).Scan(
		&rating.ID, &rating.UserID, &rating.BookID, &rating.Score,
		&rating.CreatedAt, &rating.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return rating, err
}

func (r *ratingRepository) GetByBookID(ctx context.Context, bookID int) ([]*entity.Rating, error) {
	query := `
        SELECT id, user_id, book_id, score, created_at, updated_at
        FROM ratings 
        WHERE book_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []*entity.Rating
	for rows.Next() {
		rating := &entity.Rating{}
		err := rows.Scan(
			&rating.ID, &rating.UserID, &rating.BookID, &rating.Score,
			&rating.CreatedAt, &rating.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	return ratings, nil
}

func (r *ratingRepository) GetByUserID(ctx context.Context, userID int) ([]*entity.Rating, error) {
	query := `
        SELECT id, user_id, book_id, score, created_at, updated_at
        FROM ratings 
        WHERE user_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []*entity.Rating
	for rows.Next() {
		rating := &entity.Rating{}
		err := rows.Scan(
			&rating.ID, &rating.UserID, &rating.BookID, &rating.Score,
			&rating.CreatedAt, &rating.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	return ratings, nil
}

func (r *ratingRepository) GetAverageRating(ctx context.Context, bookID int) (float64, int, error) {
	query := `
        SELECT 
            COALESCE(AVG(score), 0) as average_rating,
            COUNT(*) as total_ratings
        FROM ratings 
        WHERE book_id = $1
    `

	var averageRating float64
	var totalRatings int

	err := r.db.QueryRowContext(ctx, query, bookID).Scan(&averageRating, &totalRatings)
	if err != nil {
		return 0, 0, err
	}

	return averageRating, totalRatings, nil
}

func (r *ratingRepository) Delete(ctx context.Context, userID, bookID int) error {
	query := `DELETE FROM ratings WHERE user_id = $1 AND book_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, bookID)
	return err
}

func (r *ratingRepository) GetBooksWithUserRatings(ctx context.Context, userID int) ([]*entity.BookWithRating, error) {
	query := `
        SELECT 
            b.id, b.title, b.author_id, b.isbn, b.description, 
            b.published_at, b.created_at, b.updated_at,
            a.id, a.first_name, a.last_name, a.biography, a.birth_date,
            a.created_at, a.updated_at,
            COALESCE(AVG(r_all.score), 0) as average_rating,
            COUNT(r_all.id) as total_ratings,
            r_user.score as user_rating
        FROM books b
        LEFT JOIN authors a ON b.author_id = a.id
        LEFT JOIN ratings r_all ON b.id = r_all.book_id
        LEFT JOIN ratings r_user ON b.id = r_user.book_id AND r_user.user_id = $1
        GROUP BY b.id, a.id, r_user.score
        ORDER BY b.title
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*entity.BookWithRating
	for rows.Next() {
		book := &entity.BookWithRating{
			Book: entity.Book{Author: &entity.Author{}},
		}
		var userRating sql.NullInt64

		err := rows.Scan(
			&book.ID, &book.Title, &book.AuthorID, &book.ISBN, &book.Description,
			&book.PublishedAt, &book.CreatedAt, &book.UpdatedAt,
			&book.Author.ID, &book.Author.FirstName, &book.Author.LastName,
			&book.Author.Biography, &book.Author.BirthDate,
			&book.Author.CreatedAt, &book.Author.UpdatedAt,
			&book.AverageRating, &book.TotalRatings, &userRating,
		)
		if err != nil {
			return nil, err
		}

		if userRating.Valid {
			userRatingValue := int(userRating.Int64)
			book.UserRating = &userRatingValue
		}

		books = append(books, book)
	}

	return books, nil
}
