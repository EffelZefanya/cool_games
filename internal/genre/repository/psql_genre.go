package repository

import (
    "context"
    "cool-games/internal/domain"
    "database/sql"
)

type psqlGenreRepository struct {
    db *sql.DB
}

func NewPsqlGenreRepository(db *sql.DB) domain.GenreRepository {
    return &psqlGenreRepository{db}
}

func (m *psqlGenreRepository) Fetch(ctx context.Context) ([]domain.Genre, error) {
    rows, err := m.db.QueryContext(ctx, "SELECT id, genre_name FROM genres")
    if err != nil { return nil, err }
    defer rows.Close()

    var res []domain.Genre
    for rows.Next() {
        var g domain.Genre
        if err := rows.Scan(&g.ID, &g.Name); err != nil { return nil, err }
        res = append(res, g)
    }
    return res, nil
}

func (m *psqlGenreRepository) Store(ctx context.Context, g *domain.Genre) error {
    return m.db.QueryRowContext(ctx, "INSERT INTO genres (genre_name) VALUES ($1) RETURNING id", g.Name).Scan(&g.ID)
}

func (m *psqlGenreRepository) Update(ctx context.Context, g *domain.Genre) error {
    _, err := m.db.ExecContext(ctx, "UPDATE genres SET genre_name = $1 WHERE id = $2", g.Name, g.ID)
    return err
}

func (m *psqlGenreRepository) Delete(ctx context.Context, id int) error {
    _, err := m.db.ExecContext(ctx, "DELETE FROM genres WHERE id = $1", id)
    return err
}