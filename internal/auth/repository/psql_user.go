package repository

import (
	"context"
	"database/sql"
	"cool-games/internal/domain"
)

type psqlUserRepository struct {
	db *sql.DB
}

func NewPsqlUserRepository(db *sql.DB) domain.UserRepository {
	return &psqlUserRepository{db}
}

func (m *psqlUserRepository) Create(ctx context.Context, u *domain.User) error {
    query := `INSERT INTO users (email, hashed_password, role) VALUES ($1, $2, $3) RETURNING id, created_at`
    return m.db.QueryRowContext(ctx, query, u.Email, u.HashedPassword, u.Role).Scan(&u.ID, &u.CreatedAt)
}

func (m *psqlUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT id, email, hashed_password, role FROM users WHERE email = $1 AND deleted_at IS NULL`
	var u domain.User
	err := m.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.HashedPassword, &u.Role)
	return u, err
}

func (m *psqlUserRepository) CreatePublisher(ctx context.Context, userID int, name string) error {
    query := `INSERT INTO publishers (user_id, publisher_name) VALUES ($1, $2)`
    _, err := m.db.ExecContext(ctx, query, userID, name)
    return err
}