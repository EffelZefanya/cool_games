package repository

import (
	"context"
	"cool-games/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type psqlGameRepository struct {
	db *sql.DB
}

func NewPsqlGameRepository(db *sql.DB) domain.GameRepository {
	return &psqlGameRepository{db}
}

func (m *psqlGameRepository) getGenresForGame(ctx context.Context, gameID int) ([]domain.Genre, error) {
    query := `
        SELECT g.id, g.genre_name 
        FROM genres g 
        JOIN game_genres gg ON g.id = gg.genre_id 
        WHERE gg.game_id = $1`
    
    rows, err := m.db.QueryContext(ctx, query, gameID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    genres := []domain.Genre{}
    for rows.Next() {
        var gen domain.Genre
        if err := rows.Scan(&gen.ID, &gen.Name); err != nil {
            return nil, err
        }
        genres = append(genres, gen)
    }
    return genres, nil
}

func (m *psqlGameRepository) Store(ctx context.Context, g *domain.Game) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO games (publisher_id, developer_id, game_name, price, stock_level) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	if err := tx.QueryRowContext(ctx, query, g.PublisherID, g.DeveloperID, g.Name, g.Price, g.StockLevel).Scan(&g.ID); err != nil {
		return err
	}

	for _, gen := range g.Genres {
		_, err := tx.ExecContext(ctx, "INSERT INTO game_genres (game_id, genre_id) VALUES ($1, $2)", g.ID, gen.ID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (m *psqlGameRepository) Fetch(ctx context.Context, search string, minPrice, maxPrice float64) ([]domain.Game, error) {
	query := `SELECT id, publisher_id, developer_id, game_name, price, stock_level FROM games WHERE deleted_at IS NULL`
	args := []interface{}{}
	argCount := 1

	if search != "" {
		query += fmt.Sprintf(" AND game_name ILIKE $%d", argCount)
		args = append(args, "%"+search+"%")
		argCount++
	}
	if minPrice > 0 {
		query += fmt.Sprintf(" AND price >= $%d", argCount)
		args = append(args, minPrice)
		argCount++
	}
	if maxPrice > 0 {
		query += fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, maxPrice)
		argCount++
	}

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Game
    for rows.Next() {
        var g domain.Game
        err := rows.Scan(&g.ID, &g.PublisherID, &g.DeveloperID, &g.Name, &g.Price, &g.StockLevel)
        if err != nil { return nil, err }
        
        g.Genres, _ = m.getGenresForGame(ctx, g.ID)
        res = append(res, g)
    }
    return res, nil
}

func (m *psqlGameRepository) GetByID(ctx context.Context, id int) (domain.Game, error) {
	query := `SELECT id, publisher_id, developer_id, game_name, price, stock_level 
              FROM games WHERE id = $1 AND deleted_at IS NULL`
	var g domain.Game
	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&g.ID, &g.PublisherID, &g.DeveloperID, &g.Name, &g.Price, &g.StockLevel,
	)
	if err != nil {
		return domain.Game{}, domain.ErrGameNotFound
	}

	g.Genres, _ = m.getGenresForGame(ctx, g.ID)
    return g, nil
}

func (m *psqlGameRepository) Update(ctx context.Context, g *domain.Game) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE games SET developer_id=$1, game_name=$2, price=$3, stock_level=$4, updated_at=NOW() WHERE id=$5`
	_, err = tx.ExecContext(ctx, query, g.DeveloperID, g.Name, g.Price, g.StockLevel, g.ID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM game_genres WHERE game_id = $1", g.ID)
	if err != nil {
		return err
	}

	for _, gen := range g.Genres {
		_, err := tx.ExecContext(ctx, "INSERT INTO game_genres (game_id, genre_id) VALUES ($1, $2)", g.ID, gen.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (m *psqlGameRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE games SET deleted_at = $1 WHERE id = $2`
	_, err := m.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (m *psqlGameRepository) UpdateStock(ctx context.Context, gameID int, change int) error {
    tx, err := m.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    query := `UPDATE games SET stock_level = stock_level + $1 WHERE id = $2 AND stock_level + $1 >= 0`
    res, err := tx.ExecContext(ctx, query, change, gameID)
    if err != nil {
        return err
    }
    
    rows, _ := res.RowsAffected()
    if rows == 0 {
        return errors.New("could not update stock (insufficient stock or game not found)")
    }

    historyQuery := `INSERT INTO game_quantity_history (game_id, change_amount, transaction_date) VALUES ($1, $2, NOW())`
    if _, err := tx.ExecContext(ctx, historyQuery, gameID, change); err != nil {
        return err
    }

    return tx.Commit()
}

func (m *psqlGameRepository) FetchByPublisher(ctx context.Context, publisherID int) ([]domain.Game, error) {
	query := `SELECT id, publisher_id, developer_id, game_name, price, stock_level 
              FROM games WHERE publisher_id = $1 AND deleted_at IS NULL`
	rows, err := m.db.QueryContext(ctx, query, publisherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Game
	for rows.Next() {
		var g domain.Game
		rows.Scan(&g.ID, &g.PublisherID, &g.DeveloperID, &g.Name, &g.Price, &g.StockLevel)
		g.Genres, _ = m.getGenresForGame(ctx, g.ID)
		res = append(res, g)
	}
	return res, nil
}

func (m *psqlGameRepository) GetPublisherIDByUserID(ctx context.Context, userID int) (int, error) {
	var id int
	query := `SELECT id FROM publishers WHERE user_id = $1`
	err := m.db.QueryRowContext(ctx, query, userID).Scan(&id)
	return id, err
}