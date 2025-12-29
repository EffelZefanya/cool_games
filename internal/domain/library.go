package domain

import (
	"context"
	"time"
)

type LibraryEntry struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	GameID     int       `json:"game_id"`
	PurchasedAt time.Time `json:"purchased_at"`
}

type LibraryRepository interface {
    AddToLibrary(ctx context.Context, userID int, gameID int) error
    GetOwnedGames(ctx context.Context, userID int) ([]Game, error)
}