package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUnauthorizedAction = errors.New("you are not authorized to modify this resource")
	ErrGameNotFound       = errors.New("game not found")
)

type Game struct {
    ID          int       `json:"id"`
    PublisherID int       `json:"publisher_id"`
    DeveloperID int       `json:"developer_id" binding:"required"`
    Name        string    `json:"game_name" binding:"required"`
    Price       float64   `json:"price" binding:"required"`
    StockLevel  int       `json:"stock_level"`
    Genres      []Genre   `json:"genres"`
    ReleaseDate time.Time `json:"release_date"`
}

type RestockRequest struct {
	Amount int `json:"amount" binding:"required,gt=0"`
}

type GameRepository interface {
	Fetch(ctx context.Context, search string, minPrice, maxPrice float64) ([]Game, error)
	GetByID(ctx context.Context, id int) (Game, error)
	Store(ctx context.Context, game *Game) error
	Update(ctx context.Context, game *Game) error
	Delete(ctx context.Context, id int) error
	UpdateStock(ctx context.Context, gameID int, change int) error
	FetchByPublisher(ctx context.Context, publisherID int) ([]Game, error)
	GetPublisherIDByUserID(ctx context.Context, userID int) (int, error)
}

type GameUsecase interface {
    GetAll(ctx context.Context, search string, minPrice, maxPrice float64) ([]Game, error) 
    GetByID(ctx context.Context, id int) (Game, error)
    GetByPublisher(ctx context.Context, publisherID int) ([]Game, error)
    Create(ctx context.Context, game *Game, requesterID int) error
    Update(ctx context.Context, id int, game *Game, requesterID int, role string) error
    Delete(ctx context.Context, id int, requesterID int, role string) error 
    Restock(ctx context.Context, gameID int, requesterID int, amount int) error
}