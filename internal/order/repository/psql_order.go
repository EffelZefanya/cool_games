package repository

import (
	"context"
	"cool-games/internal/domain"
	"database/sql"
	"errors"
)

type psqlOrderRepository struct {
	db *sql.DB
}

func NewPsqlOrderRepository(db *sql.DB) *psqlOrderRepository {
	return &psqlOrderRepository{db: db}
}

func (r *psqlOrderRepository) GetPublisherSales(ctx context.Context, publisherID int) ([]domain.SalesReportEntry, error) {
	query := `
        SELECT g.id, g.game_name, g.price, cgl.purchase_date, u.email -- REMOVED the 'd'
        FROM customer_game_library cgl
        JOIN games g ON cgl.game_id = g.id
        JOIN customers c ON cgl.customer_id = c.id
        JOIN users u ON c.user_id = u.id
        WHERE g.publisher_id = $1
        ORDER BY cgl.purchase_date DESC`

	rows, err := r.db.QueryContext(ctx, query, publisherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var report []domain.SalesReportEntry
	for rows.Next() {
		var e domain.SalesReportEntry
		if err := rows.Scan(&e.GameID, &e.GameName, &e.PriceAtSale, &e.PurchasedDate, &e.CustomerEmail); err != nil {
			return nil, err
		}
		report = append(report, e)
	}
	return report, nil
}

func (r *psqlOrderRepository) ExecutePurchase(ctx context.Context, userID int, gameID int, price float64) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()

    var customerID int
    err = tx.QueryRowContext(ctx, "SELECT id FROM customers WHERE user_id = $1", userID).Scan(&customerID)
    if err != nil { return errors.New("customer profile not found") }

    res, err := tx.ExecContext(ctx, 
        "UPDATE customers SET current_balance = current_balance - $1 WHERE id = $2 AND current_balance >= $1", 
        price, customerID)
    if err != nil { return err }
    if rows, _ := res.RowsAffected(); rows == 0 { return errors.New("insufficient balance") }

    res, err = tx.ExecContext(ctx, 
        "UPDATE games SET stock_level = stock_level - 1 WHERE id = $1 AND stock_level > 0", 
        gameID)
    if err != nil { return err }
    if rows, _ := res.RowsAffected(); rows == 0 { return errors.New("game out of stock") }

	_, err = tx.ExecContext(ctx, `
        INSERT INTO game_quantity_history (game_id, change_amount, transaction_date) 
        VALUES ($1, -1, NOW())`, gameID)
    if err != nil { return err }

    var orderID int
    queryOrder := `
        INSERT INTO orders (customer_id, game_id, qty, total_amount, order_date) 
        VALUES ($1, $2, 1, $3, NOW()) RETURNING id`
    
    err = tx.QueryRowContext(ctx, queryOrder, customerID, gameID, price).Scan(&orderID)
    if err != nil { return err }

    _, err = tx.ExecContext(ctx, `
        INSERT INTO customer_game_library (customer_id, game_id, purchase_date) 
        VALUES ($1, $2, NOW()) 
        ON CONFLICT (customer_id, game_id) DO NOTHING`, customerID, gameID)
    if err != nil { return err }

    queryLedger := `
        INSERT INTO ledger (customer_id, order_id, amount, type, transaction_date) 
        VALUES ($1, $2, $3, 'debit', NOW())`
    
    _, err = tx.ExecContext(ctx, queryLedger, customerID, orderID, price)
    if err != nil { return err }

    return tx.Commit()
}

func (r *psqlOrderRepository) RecordLedger(ctx context.Context, userID int, amount float64, description string) error {
    query := `
        INSERT INTO ledger (customer_id, amount, type, transaction_date) 
        SELECT id, $1, 'credit', NOW() FROM customers WHERE user_id = $2`
    _, err := r.db.ExecContext(ctx, query, amount, userID)
    return err
}