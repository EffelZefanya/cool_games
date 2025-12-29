package repository

import (
	"context"
	"database/sql"
	"cool-games/internal/domain"
)

type psqlCustomerRepository struct {
	db *sql.DB
}

func NewPsqlCustomerRepository(db *sql.DB) domain.CustomerRepository {
	return &psqlCustomerRepository{db}
}

func (m *psqlCustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
	query := `INSERT INTO customers (user_id, customer_name, current_balance) VALUES ($1, $2, $3)`
	_, err := m.db.ExecContext(ctx, query, c.UserID, c.CustomerName, c.CurrentBalance)
	return err
}

func (m *psqlCustomerRepository) GetByUserID(ctx context.Context, userID int) (domain.Customer, error) {
	query := `SELECT id, user_id, customer_name, current_balance, created_at FROM customers WHERE user_id = $1`
	var c domain.Customer
	err := m.db.QueryRowContext(ctx, query, userID).Scan(&c.ID, &c.UserID, &c.CustomerName, &c.CurrentBalance, &c.CreatedAt)
	if err != nil {
		return domain.Customer{}, err
	}
	return c, nil
}

func (m *psqlCustomerRepository) UpdateBalance(ctx context.Context, userID int, amount float64) error {
	query := `UPDATE customers SET current_balance = current_balance + $1 WHERE user_id = $2`
	
	result, err := m.db.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return domain.ErrGameNotFound
	}
	return nil
}