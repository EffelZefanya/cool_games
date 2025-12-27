package usecase

import (
	"context"
	"cool-games/internal/domain"
	"errors"
	"time"
)

type orderUsecase struct {
	gameRepo     domain.GameRepository
	customerRepo domain.CustomerRepository
	orderRepo domain.OrderRepository
	libraryRepo  domain.LibraryRepository
	timeout      time.Duration
}


func NewOrderUsecase(
    g domain.GameRepository, 
    c domain.CustomerRepository, 
    o domain.OrderRepository, 
    l domain.LibraryRepository,
    t time.Duration,
) domain.OrderUsecase {
    return &orderUsecase{
        gameRepo:     g,
        customerRepo: c,
        orderRepo:    o,
        libraryRepo:  l,
        timeout:      t,
    }
}

func (u *orderUsecase) BuyGame(ctx context.Context, customerID int, gameID int) error {
	c, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	game, err := u.gameRepo.GetByID(c, gameID)
	if err != nil { return err }
	if game.StockLevel <= 0 { return errors.New("insufficient stock") }

	customer, err := u.customerRepo.GetByUserID(c, customerID)
	if err != nil { return err }
	if customer.CurrentBalance < game.Price { return errors.New("insufficient balance") }

	if u.libraryRepo != nil {
		ownedGames, _ := u.libraryRepo.GetOwnedGames(c, customerID)
		for _, g := range ownedGames {
			if g.ID == gameID {
				return errors.New("you already own this game")
			}
		}
	}

	return u.orderRepo.ExecutePurchase(c, customerID, gameID, game.Price)
}

func (u *orderUsecase) AddBalance(ctx context.Context, userID int, amount float64) error {
    c, cancel := context.WithTimeout(ctx, u.timeout)
    defer cancel()

    err := u.customerRepo.UpdateBalance(c, userID, amount)
    if err != nil { return err }

    return u.orderRepo.RecordLedger(c, userID, amount, "Top-up")
}

func (u *orderUsecase) GetPublisherSalesReport(ctx context.Context, userID int) ([]domain.SalesReportEntry, error) {
    c, cancel := context.WithTimeout(ctx, u.timeout)
    defer cancel()

    pubID, err := u.gameRepo.GetPublisherIDByUserID(c, userID)
    if err != nil {
        return nil, errors.New("publisher profile not found")
    }

    return u.orderRepo.GetPublisherSales(c, pubID)
}

func (u *orderUsecase) GetCustomerLibrary(ctx context.Context, userID int) ([]domain.Game, error) {
    c, cancel := context.WithTimeout(ctx, u.timeout)
    defer cancel()

    return u.libraryRepo.GetOwnedGames(c, userID)
}