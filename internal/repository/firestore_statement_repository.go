package repository

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/tsongpon/helios/internal/model"
)

type FirestoreStatementRepository struct {
	client *firestore.Client
}

func NewFirestoreStatementRepository(client *firestore.Client) *FirestoreStatementRepository {
	return &FirestoreStatementRepository{
		client: client,
	}
}

func (r *FirestoreStatementRepository) Save(statement model.Statement) error {
	ctx := context.Background()

	transactions := make([]map[string]any, len(statement.Transactions))
	for i, t := range statement.Transactions {
		transactions[i] = map[string]any{
			"transaction_date": t.TransactionDate,
			"posting_date":     t.PostingDate,
			"description":      t.Description,
			"amount":           t.Amount,
			"is_installment":   t.IsInstallment,
			"installment_term": t.InstallmentTerm,
		}
	}

	doc := map[string]any{
		"card_number":      statement.CardNumber,
		"total_payment":    statement.TotalPayment,
		"minimum_payment":  statement.MinimumPayment,
		"payment_due_date": statement.PaymentDueDate,
		"credit_line":      statement.CreditLine,
		"transactions":     transactions,
	}

	_, err := r.client.Collection("statements").Doc(statement.CardNumber).Set(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to save statement: %w", err)
	}

	return nil
}
