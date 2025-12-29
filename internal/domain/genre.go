package domain

import "context"

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"genre_name"`
}

type GenreRepository interface {
    Fetch(ctx context.Context) ([]Genre, error)
    Store(ctx context.Context, genre *Genre) error
    Update(ctx context.Context, genre *Genre) error
    Delete(ctx context.Context, id int) error
}

type GenreUsecase interface {
    GetAll(ctx context.Context) ([]Genre, error)
    Create(ctx context.Context, genre *Genre) error
    Update(ctx context.Context, genre *Genre) error
    Delete(ctx context.Context, id int) error
}