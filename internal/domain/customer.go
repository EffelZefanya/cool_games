package domain

import (
	"context"
	"time"
)

type Customer struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	CustomerName   string    `json:"customer_name"`
	CurrentBalance float64   `json:"current_balance"`
	CreatedAt      time.Time `json:"created_at"`
}

type CustomerRepository interface {
	Create(ctx context.Context, customer *Customer) error
	GetByUserID(ctx context.Context, userID int) (Customer, error)
	UpdateBalance(ctx context.Context, userID int, amount float64) error
}

type CustomerUsecase interface {
	GetProfile(ctx context.Context, userID int) (Customer, error)
}