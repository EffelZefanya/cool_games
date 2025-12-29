package domain

import (
	"context"
	"time"
)

type PurchaseRequest struct {
	GameID int `json:"game_id" binding:"required"`
}

type SalesReportEntry struct {
    GameID        int       `json:"game_id"`
    GameName      string    `json:"game_name"`
    PriceAtSale   float64   `json:"price_at_sale"`
    PurchasedDate time.Time `json:"purchased_date"`
    CustomerEmail string    `json:"customer_email"`
}

type OrderUsecase interface {
    BuyGame(ctx context.Context, customerID int, gameID int) error
    GetPublisherSalesReport(ctx context.Context, customerID int) ([]SalesReportEntry, error)
    AddBalance(ctx context.Context, customerID int, amount float64) error
	GetCustomerLibrary(ctx context.Context, userID int) ([]Game, error)
}

type OrderRepository interface {
    ExecutePurchase(ctx context.Context, customerID int, gameID int, price float64) error
	GetPublisherSales(ctx context.Context, publisherID int) ([]SalesReportEntry, error)
	RecordLedger(ctx context.Context, customerID int, amount float64, description string) error
}