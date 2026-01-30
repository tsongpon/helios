package httphandler

import (
	"context"
	"io"
	"time"

	"github.com/tsongpon/helios/internal/model"
)

type PDFService interface {
	ExtractText(ctx context.Context, userID string, file io.Reader, password string) ([]model.Transaction, error)
}

type TransactionService interface {
	GetTransactions(ctx context.Context, userID string, from, to time.Time) ([]model.Transaction, error)
}
