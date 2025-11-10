package postgres

import (
	"awesomeProject/internal/entity"
	"context"
	"database/sql"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
        INSERT INTO users (email, password_hash, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	return r.db.QueryRowContext(ctx, query,
		user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*entity.User, error) {
	query := `
        SELECT id, email, password_hash, role, created_at, updated_at
        FROM users WHERE id = $1
    `

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
        SELECT id, email, password_hash, role, created_at, updated_at
        FROM users WHERE email = $1
    `

	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Role,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
        UPDATE users 
        SET email = $1, password_hash = $2, role = $3, updated_at = $4
        WHERE id = $5
    `

	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.PasswordHash, user.Role, user.UpdatedAt, user.ID,
	)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
