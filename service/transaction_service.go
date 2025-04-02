package service

import (
	"awesomeProject/models"
	"awesomeProject/types"
	"context"
	"fmt"
)

func (s *Impl) ProcessTransaction(ctx context.Context, id string, userid uint64, amount uint64, state string, source string) error {

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

func validateInputs(state, source string) error {
	if state != types.StateWin && state != types.StateLose {
		return fmt.Errorf("invalid transaction state: %s", state)
	}

	switch types.SourceType(source) {
	case types.SourceTypeGame, types.SourceTypeServer, types.SourceTypePayment:
		return nil
	default:
		return fmt.Errorf("invalid source type: %s", source)
	}
}
