package usecase

import (
	"context"
	"cool-games/internal/domain"
	"errors"
	"time"
)

type gameUsecase struct {
	gameRepo       domain.GameRepository
	contextTimeout time.Duration
}

func NewGameUsecase(g domain.GameRepository, timeout time.Duration) domain.GameUsecase {
	return &gameUsecase{gameRepo: g, contextTimeout: timeout}
}

func (u *gameUsecase) GetAll(ctx context.Context, search string, minPrice, maxPrice float64) ([]domain.Game, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.gameRepo.Fetch(c, search, minPrice, maxPrice)
}

func (u *gameUsecase) GetByID(ctx context.Context, id int) (domain.Game, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	return u.gameRepo.GetByID(c, id)
}

func (u *gameUsecase) Create(ctx context.Context, g *domain.Game, requesterID int) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	pubID, err := u.gameRepo.GetPublisherIDByUserID(c, requesterID)
	if err != nil {
		return errors.New("publisher profile not found")
	}

	g.PublisherID = pubID
	return u.gameRepo.Store(c, g)
}

func (u *gameUsecase) Update(ctx context.Context, id int, g *domain.Game, requesterID int, role string) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existing, err := u.gameRepo.GetByID(c, id)
	if err != nil {
		return err
	}

	if role == "publisher" {
		pubID, err := u.gameRepo.GetPublisherIDByUserID(c, requesterID)
		if err != nil || existing.PublisherID != pubID {
			return domain.ErrUnauthorizedAction
		}
	}

	g.ID = id
	g.PublisherID = existing.PublisherID
	return u.gameRepo.Update(c, g)
}

func (u *gameUsecase) Delete(ctx context.Context, id int, requesterID int, role string) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existing, err := u.gameRepo.GetByID(c, id)
	if err != nil {
		return err
	}

	if role == "publisher" {
		pubID, _ := u.gameRepo.GetPublisherIDByUserID(c, requesterID)
		if existing.PublisherID != pubID {
			return domain.ErrUnauthorizedAction
		}
	}

	return u.gameRepo.Delete(c, id)
}

func (u *gameUsecase) GetByPublisher(ctx context.Context, requesterID int) ([]domain.Game, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	pubID, err := u.gameRepo.GetPublisherIDByUserID(c, requesterID)
	if err != nil {
		return nil, err
	}
	return u.gameRepo.FetchByPublisher(c, pubID)
}

func (u *gameUsecase) Restock(ctx context.Context, gameID int, requesterID int, amount int) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	existing, err := u.gameRepo.GetByID(c, gameID)
	if err != nil {
		return err
	}

	pubID, _ := u.gameRepo.GetPublisherIDByUserID(c, requesterID)
	if existing.PublisherID != pubID {
		return domain.ErrUnauthorizedAction
	}

	return u.gameRepo.UpdateStock(c, gameID, amount)
}