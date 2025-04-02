package service

import (
	"context"
)

type Service interface {
	ProcessTransaction(ctx context.Context, id string, userid uint64, amount uint64, state string, source string) error
	GetBalance(ctx context.Context, id uint64) (uint64, error)
	UpdateBalance(ctx context.Context, userId uint64, amount uint64, state string) error
	ValidateBalance(ctx context.Context, userId uint64, amount uint64, state string) error
}
