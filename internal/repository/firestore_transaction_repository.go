package repository

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/tsongpon/helios/internal/model"
)

type FirestoreTransactionRepository struct {
	client *firestore.Client
}

func NewFirestoreTransactionRepository(client *firestore.Client) *FirestoreTransactionRepository {
	return &FirestoreTransactionRepository{
		client: client,
	}
}

func (r *FirestoreTransactionRepository) Save(ctx context.Context, transactions []model.Transaction) error {
	batch := r.client.Batch()
	collection := r.client.Collection("transactions")

	for _, t := range transactions {
		doc := map[string]any{
			"card_number":      t.CardNumber,
			"user_id":          t.UserID,
			"transaction_date": t.TransactionDate,
			"posting_date":     t.PostingDate,
			"description":      t.Description,
			"amount":           t.Amount,
			"is_installment":   t.IsInstallment,
			"installment_term": t.InstallmentTerm,
		}
		docRef := collection.NewDoc()
		batch.Set(docRef, doc)
	}

	_, err := batch.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to save transactions: %w", err)
	}

	return nil
}

func (r *FirestoreTransactionRepository) GetTransactions(ctx context.Context, userID string, from, to time.Time) ([]model.Transaction, error) {
	fromStr := from.Format("2006-01-02")
	toStr := to.Format("2006-01-02")

	docs, err := r.client.Collection("transactions").
		Where("user_id", "==", userID).
		Where("transaction_date", ">=", fromStr).
		Where("transaction_date", "<=", toStr).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	transactions := make([]model.Transaction, 0, len(docs))
	for _, doc := range docs {
		data := doc.Data()
		t := model.Transaction{
			UserID:          stringVal(data, "user_id"),
			CardNumber:      stringVal(data, "card_number"),
			TransactionDate: stringVal(data, "transaction_date"),
			PostingDate:     stringVal(data, "posting_date"),
			Description:     stringVal(data, "description"),
			Amount:          floatVal(data, "amount"),
			IsInstallment:   boolVal(data, "is_installment"),
			InstallmentTerm: stringVal(data, "installment_term"),
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func stringVal(data map[string]any, key string) string {
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}

func floatVal(data map[string]any, key string) float64 {
	if v, ok := data[key].(float64); ok {
		return v
	}
	return 0
}

func boolVal(data map[string]any, key string) bool {
	if v, ok := data[key].(bool); ok {
		return v
	}
	return false
}
