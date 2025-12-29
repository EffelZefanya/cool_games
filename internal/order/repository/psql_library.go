package repository

import (
	"context"
	"database/sql"
	"cool-games/internal/domain"
)

type psqlLibraryRepository struct {
	db *sql.DB
}

func NewPsqlLibraryRepository(db *sql.DB) domain.LibraryRepository {
	return &psqlLibraryRepository{db: db}
}

func (r *psqlLibraryRepository) AddToLibrary(ctx context.Context, userID int, gameID int) error {
	query := `
		INSERT INTO customer_game_library (customer_id, game_id)
		SELECT id, $2 FROM customers WHERE user_id = $1`
	
	_, err := r.db.ExecContext(ctx, query, userID, gameID)
	return err
}

func (r *psqlLibraryRepository) GetOwnedGames(ctx context.Context, userID int) ([]domain.Game, error) {
	query := `
		SELECT g.id, g.publisher_id, g.developer_id, g.game_name, g.price, g.stock_level
		FROM games g
		INNER JOIN customer_game_library cgl ON g.id = cgl.game_id
		INNER JOIN customers c ON cgl.customer_id = c.id
		WHERE c.user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []domain.Game
	for rows.Next() {
		var g domain.Game
		err := rows.Scan(&g.ID, &g.PublisherID, &g.DeveloperID, &g.Name, &g.Price, &g.StockLevel)
		if err != nil {
			return nil, err
		}
		games = append(games, g)
	}

	return games, nil
}