package service

import (
	"context"
	"time"

	"github.com/tsongpon/helios/internal/model"
)

type LLMRepository interface {
	ParseStatement(statementText string) ([]model.Transaction, error)
}

type TransactionRepository interface {
	Save(ctx context.Context, transactions []model.Transaction) error
	GetTransactions(ctx context.Context, userID string, from, to time.Time) ([]model.Transaction, error)
}
