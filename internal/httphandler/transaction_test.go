package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/tsongpon/helios/internal/model"
)

type mockTransactionService struct {
	transactions []model.Transaction
	err          error
}

func (m *mockTransactionService) GetTransactions(ctx context.Context, userID string, from, to time.Time) ([]model.Transaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.transactions, nil
}

func TestTransactionHandler_GetTransactions(t *testing.T) {
	t.Run("returns transactions successfully", func(t *testing.T) {
		mockService := &mockTransactionService{
			transactions: []model.Transaction{
				{
					UserID:          "1234567890",
					CardNumber:      "1234-XXXX-XXXX-5678",
					TransactionDate: "2024-12-15",
					PostingDate:     "2024-12-16",
					Description:     "AMAZON",
					Amount:          100.50,
					IsInstallment:   false,
				},
			},
		}

		handler := NewTransactionHandler(mockService)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/transactions?start=2024-12-01&end=2024-12-31", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetTransactions(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var response []TransactionResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(response) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(response))
		}

		if response[0].Description != "AMAZON" {
			t.Errorf("expected description AMAZON, got %s", response[0].Description)
		}
	})

	t.Run("returns error when start date is missing", func(t *testing.T) {
		mockService := &mockTransactionService{}
		handler := NewTransactionHandler(mockService)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/transactions?end=2024-12-31", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetTransactions(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		var response ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error != "start and end query parameters are required (format: YYYY-MM-DD)" {
			t.Errorf("unexpected error message: %s", response.Error)
		}
	})

	t.Run("returns error when end date is missing", func(t *testing.T) {
		mockService := &mockTransactionService{}
		handler := NewTransactionHandler(mockService)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/transactions?start=2024-12-01", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetTransactions(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("returns error for invalid start date format", func(t *testing.T) {
		mockService := &mockTransactionService{}
		handler := NewTransactionHandler(mockService)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/transactions?start=12-01-2024&end=2024-12-31", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetTransactions(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		var response ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error != "invalid start date format, expected YYYY-MM-DD" {
			t.Errorf("unexpected error message: %s", response.Error)
		}
	})

	t.Run("returns error for invalid end date format", func(t *testing.T) {
		mockService := &mockTransactionService{}
		handler := NewTransactionHandler(mockService)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/transactions?start=2024-12-01&end=31-12-2024", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetTransactions(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		var response ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error != "invalid end date format, expected YYYY-MM-DD" {
			t.Errorf("unexpected error message: %s", response.Error)
		}
	})

	t.Run("returns error when service fails", func(t *testing.T) {
		mockService := &mockTransactionService{
			err: errors.New("database error"),
		}
		handler := NewTransactionHandler(mockService)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/transactions?start=2024-12-01&end=2024-12-31", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetTransactions(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}

		var response ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error != "failed to get transactions: database error" {
			t.Errorf("unexpected error message: %s", response.Error)
		}
	})

	t.Run("returns empty list when no transactions found", func(t *testing.T) {
		mockService := &mockTransactionService{
			transactions: []model.Transaction{},
		}
		handler := NewTransactionHandler(mockService)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/transactions?start=2024-12-01&end=2024-12-31", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetTransactions(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var response []TransactionResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(response) != 0 {
			t.Errorf("expected 0 transactions, got %d", len(response))
		}
	})
}
