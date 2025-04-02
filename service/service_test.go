package service

import (
	"awesomeProject/models"
	"awesomeProject/types"
	"context"
	"fmt"
	"testing"
)

type transactionRepository interface {
	CreateTransaction(ctx context.Context, transaction models.Transaction) error
	ExistsTransaction(ctx context.Context, id string) (bool, error)
}

type balanceRepository interface {
	GetBalance(ctx context.Context, id uint64) (uint64, error)
	UpdateBalance(ctx context.Context, id uint64, balance uint64) error
}

type mockRepo struct {
	exists    bool
	balance   uint64
	createErr error
}

func (m *mockRepo) ExistsTransaction(ctx context.Context, id string) (bool, error) {
	return m.exists, nil
}

func (m *mockRepo) CreateTransaction(ctx context.Context, transaction models.Transaction) error {
	return m.createErr
}

func (m *mockRepo) GetBalance(ctx context.Context, id uint64) (uint64, error) {
	return m.balance, nil
}

func (m *mockRepo) UpdateBalance(ctx context.Context, id uint64, balance uint64) error {
	return nil
}

type testImpl struct {
	transaction transactionRepository
	balance     balanceRepository
}

func newTestImpl(tx transactionRepository, bal balanceRepository) *testImpl {
	return &testImpl{
		transaction: tx,
		balance:     bal,
	}
}

func (s *testImpl) ProcessTransaction(ctx context.Context, id string, userid uint64, amount uint64, state string, source string) error {

	if err := validateInputs(state, source); err != nil {
		return err
	}

	exists, err := s.transaction.ExistsTransaction(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check transaction: %w", err)
	}
	if exists {
		return nil
	}

	if state == types.StateLose {
		balance, err := s.balance.GetBalance(ctx, userid)
		if err != nil {
			return fmt.Errorf("failed to get balance: %w", err)
		}
		if balance < amount {
			return fmt.Errorf("insufficient balance: current=%d, required=%d", balance, amount)
		}
	}

	transaction := models.Transaction{
		ID:         id,
		UserID:     userid,
		Amount:     amount,
		State:      state,
		SourceType: source,
	}

	if err := s.transaction.CreateTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return s.UpdateBalance(ctx, userid, amount, state)
}

func (s *testImpl) GetBalance(ctx context.Context, id uint64) (uint64, error) {
	return s.balance.GetBalance(ctx, id)
}

func (s *testImpl) UpdateBalance(ctx context.Context, userId uint64, amount uint64, state string) error {

	if err := s.ValidateBalance(ctx, userId, amount, state); err != nil {
		return err
	}
	currentBalance, err := s.GetBalance(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	var newBalance uint64
	switch state {
	case types.StateWin:
		newBalance = currentBalance + amount
	case types.StateLose:
		newBalance = currentBalance - amount
	default:
		return fmt.Errorf("invalid state: %s", state)
	}
	return s.balance.UpdateBalance(ctx, userId, newBalance)
}

func (s *testImpl) ValidateBalance(ctx context.Context, userId uint64, amount uint64, state string) error {

	if state == types.StateLose {
		currentBalance, err := s.GetBalance(ctx, userId)
		if err != nil {
			return fmt.Errorf("failed to get current balance: %w", err)
		}

		if currentBalance < amount {
			return fmt.Errorf("insufficient balance: current=%d, required=%d", currentBalance, amount)
		}
	}
	return nil
}

func TestProcessTransaction(t *testing.T) {
	tests := []struct {
		name      string
		repo      mockRepo
		amount    uint64
		state     string
		wantError bool
	}{
		{
			name:   "successful win",
			repo:   mockRepo{balance: 1000},
			amount: 100,
			state:  types.StateWin,
		},
		{
			name:      "insufficient balance",
			repo:      mockRepo{balance: 50},
			amount:    100,
			state:     types.StateLose,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestImpl(&tt.repo, &tt.repo)
			err := svc.ProcessTransaction(context.Background(), "tx1", 1, tt.amount, tt.state, string(types.SourceTypeGame))
			if (err != nil) != tt.wantError {
				t.Errorf("ProcessTransaction() error = %v, wantError = %v", err, tt.wantError)
			}
		})
	}
}
