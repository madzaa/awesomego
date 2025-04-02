package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type UserBalanceRepository struct {
	db *sql.DB
}

func NewUserBalanceRepository(db *sql.DB) *UserBalanceRepository {
	return &UserBalanceRepository{
		db: db,
	}
}

func (u *UserBalanceRepository) GetBalance(ctx context.Context, id uint64) (uint64, error) {
	var balance uint64
	err := u.db.QueryRowContext(ctx, "SELECT balance FROM users where id = $1", id).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found id=%d: no such id", id)
		}
		return 0, fmt.Errorf("failed to get balance %w", err)
	}
	return balance, nil
}

func (u *UserBalanceRepository) UpdateBalance(ctx context.Context, id uint64, balance uint64) error {
	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, "UPDATE users SET balance = $1 WHERE id = $2", balance, id)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found: id=%d", id)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
