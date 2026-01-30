package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tsongpon/helios/internal/model"
)

type mockTransactionRepository struct {
	transactions []model.Transaction
	err          error
	savedTxns    []model.Transaction
}

func (m *mockTransactionRepository) Save(ctx context.Context, transactions []model.Transaction) error {
	m.savedTxns = transactions
	return m.err
}

func (m *mockTransactionRepository) GetTransactions(ctx context.Context, userID string, from, to time.Time) ([]model.Transaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.transactions, nil
}

func TestTransactionService_GetTransactions(t *testing.T) {
	ctx := context.Background()
	from := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	t.Run("returns transactions successfully", func(t *testing.T) {
		expectedTxns := []model.Transaction{
			{
				UserID:          "user123",
				CardNumber:      "1234-XXXX-XXXX-5678",
				TransactionDate: "2024-12-15",
				PostingDate:     "2024-12-16",
				Description:     "AMAZON",
				Amount:          100.50,
				IsInstallment:   false,
				InstallmentTerm: "",
			},
			{
				UserID:          "user123",
				CardNumber:      "1234-XXXX-XXXX-5678",
				TransactionDate: "2024-12-20",
				PostingDate:     "2024-12-21",
				Description:     "LAZADA",
				Amount:          250.00,
				IsInstallment:   true,
				InstallmentTerm: "01/06",
			},
		}

		mockRepo := &mockTransactionRepository{
			transactions: expectedTxns,
		}

		svc := NewTransactionService(mockRepo)
		transactions, err := svc.GetTransactions(ctx, "user123", from, to)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(transactions) != len(expectedTxns) {
			t.Fatalf("expected %d transactions, got %d", len(expectedTxns), len(transactions))
		}

		for i, txn := range transactions {
			if txn.Description != expectedTxns[i].Description {
				t.Errorf("transaction %d: expected description %q, got %q", i, expectedTxns[i].Description, txn.Description)
			}
			if txn.Amount != expectedTxns[i].Amount {
				t.Errorf("transaction %d: expected amount %f, got %f", i, expectedTxns[i].Amount, txn.Amount)
			}
		}
	})

	t.Run("returns empty list when no transactions found", func(t *testing.T) {
		mockRepo := &mockTransactionRepository{
			transactions: []model.Transaction{},
		}

		svc := NewTransactionService(mockRepo)
		transactions, err := svc.GetTransactions(ctx, "user123", from, to)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(transactions) != 0 {
			t.Fatalf("expected 0 transactions, got %d", len(transactions))
		}
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		mockRepo := &mockTransactionRepository{
			err: errors.New("database connection failed"),
		}

		svc := NewTransactionService(mockRepo)
		_, err := svc.GetTransactions(ctx, "user123", from, to)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if err.Error() != "database connection failed" {
			t.Errorf("expected error message %q, got %q", "database connection failed", err.Error())
		}
	})
}
