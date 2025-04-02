package repository

import (
	"awesomeProject/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (t *TransactionRepository) CreateTransaction(ctx context.Context, transaction models.Transaction) error {
	result, err := t.db.ExecContext(ctx,
		"INSERT INTO transactions (id,user_id,amount,state,source_type,created_at) VALUES ($1, $2, $3, $4, $5, $6)", transaction.ID, transaction.UserID, transaction.Amount, transaction.State, transaction.SourceType, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create transaction %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("could not create transactiod id=%s", transaction.ID)
	}
	return nil
}

func (t *TransactionRepository) GetTransaction(ctx context.Context, id uint64) (*models.Transaction, error) {
	transaction := &models.Transaction{}
	var createdAt time.Time

	err := t.db.QueryRowContext(ctx, "SELECT * FROM transactions WHERE id=$1", id).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.State,
		&transaction.SourceType,
		&createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found id=%d, %w", id, err)
		}
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	return transaction, nil
}

func (t *TransactionRepository) ExistsTransaction(ctx context.Context, id string) (bool, error) {
	var exists bool
	err := t.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM transactions WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check for transaction %w", err)
	}
	return exists, nil
}
