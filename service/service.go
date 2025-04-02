package service

import "awesomeProject/repository"

type Impl struct {
	transaction *repository.TransactionRepository
	balance     *repository.UserBalanceRepository
}

func NewServiceImpl(transaction *repository.TransactionRepository, balance *repository.UserBalanceRepository) *Impl {
	return &Impl{transaction: transaction, balance: balance}
}
