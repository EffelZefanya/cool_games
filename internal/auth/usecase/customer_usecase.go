package usecase

import (
	"context"
	"cool-games/internal/domain"
	"time"
)

type customerUsecase struct {
	customerRepo   domain.CustomerRepository
	contextTimeout time.Duration
}

func NewCustomerUsecase(repo domain.CustomerRepository, timeout time.Duration) domain.CustomerUsecase {
	return &customerUsecase{
		customerRepo:   repo,
		contextTimeout: timeout,
	}
}

func (u *customerUsecase) GetProfile(ctx context.Context, userID int) (domain.Customer, error) {
    c, cancel := context.WithTimeout(ctx, u.contextTimeout)
    defer cancel()
    return u.customerRepo.GetByUserID(c, userID)
}