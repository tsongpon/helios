package service

import (
	"context"
	"time"

	"github.com/tsongpon/helios/internal/model"
)

type TransactionService struct {
	transactionRepository TransactionRepository
}

func NewTransactionService(transactionRepository TransactionRepository) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
	}
}

func (s *TransactionService) GetTransactions(ctx context.Context, userID string, from, to time.Time) ([]model.Transaction, error) {
	return s.transactionRepository.GetTransactions(ctx, userID, from, to)
}
