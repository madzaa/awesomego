package service

import (
	"awesomeProject/types"
	"context"
	"fmt"
)

func (s *Impl) GetBalance(ctx context.Context, id uint64) (uint64, error) {
	balance, err := s.balance.GetBalance(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("could not get balance: %w", err)
	}
	return balance, err
}

func (s *Impl) UpdateBalance(ctx context.Context, userId uint64, amount uint64, state string) error {
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
	if err := s.balance.UpdateBalance(ctx, userId, newBalance); err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	return nil
}

func (s *Impl) ValidateBalance(ctx context.Context, userId uint64, amount uint64, state string) error {
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
