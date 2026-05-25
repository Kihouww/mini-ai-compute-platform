package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Kihouww/mini-ai-compute-platform/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (int64, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (username, password_hash, email, status)
		 VALUES (?, ?, ?, ?)`,
		user.Username,
		user.PasswordHash,
		user.Email,
		user.Status,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	var email sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, password_hash, email, status, created_at, updated_at
		 FROM users
		 WHERE username = ?
		 LIMIT 1`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&email,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	if email.Valid {
		user.Email = email.String
	}

	return &user, nil
}
