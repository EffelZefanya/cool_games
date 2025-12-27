package usecase

import (
	"context"
	"cool-games/internal/domain"
	"time"
)

type genreUsecase struct {
	genreRepo      domain.GenreRepository
	contextTimeout time.Duration
}

func NewGenreUsecase(repo domain.GenreRepository, timeout time.Duration) domain.GenreUsecase {
	return &genreUsecase{
		genreRepo:      repo,
		contextTimeout: timeout,
	}
}

func (u *genreUsecase) GetAll(ctx context.Context) ([]domain.Genre, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.genreRepo.Fetch(c)
}

func (u *genreUsecase) Create(ctx context.Context, genre *domain.Genre) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.genreRepo.Store(c, genre)
}

func (u *genreUsecase) Update(ctx context.Context, genre *domain.Genre) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.genreRepo.Update(c, genre)
}

func (u *genreUsecase) Delete(ctx context.Context, id int) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.genreRepo.Delete(c, id)
}