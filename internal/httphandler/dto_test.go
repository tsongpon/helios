package httphandler

import (
	"testing"

	"github.com/tsongpon/helios/internal/model"
)

func TestToTransactionResponses(t *testing.T) {
	t.Run("converts transactions to responses correctly", func(t *testing.T) {
		transactions := []model.Transaction{
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

		responses := toTransactionResponses(transactions)

		if len(responses) != 2 {
			t.Fatalf("expected 2 responses, got %d", len(responses))
		}

		// Check first transaction
		if responses[0].UserID != "user123" {
			t.Errorf("expected UserID user123, got %s", responses[0].UserID)
		}
		if responses[0].CardNumber != "1234-XXXX-XXXX-5678" {
			t.Errorf("expected CardNumber 1234-XXXX-XXXX-5678, got %s", responses[0].CardNumber)
		}
		if responses[0].TransactionDate != "2024-12-15" {
			t.Errorf("expected TransactionDate 2024-12-15, got %s", responses[0].TransactionDate)
		}
		if responses[0].PostingDate != "2024-12-16" {
			t.Errorf("expected PostingDate 2024-12-16, got %s", responses[0].PostingDate)
		}
		if responses[0].Description != "AMAZON" {
			t.Errorf("expected Description AMAZON, got %s", responses[0].Description)
		}
		if responses[0].Amount != 100.50 {
			t.Errorf("expected Amount 100.50, got %f", responses[0].Amount)
		}
		if responses[0].IsInstallment != false {
			t.Errorf("expected IsInstallment false, got %v", responses[0].IsInstallment)
		}
		if responses[0].InstallmentTerm != "" {
			t.Errorf("expected InstallmentTerm empty, got %s", responses[0].InstallmentTerm)
		}

		// Check second transaction (installment)
		if responses[1].Description != "LAZADA" {
			t.Errorf("expected Description LAZADA, got %s", responses[1].Description)
		}
		if responses[1].IsInstallment != true {
			t.Errorf("expected IsInstallment true, got %v", responses[1].IsInstallment)
		}
		if responses[1].InstallmentTerm != "01/06" {
			t.Errorf("expected InstallmentTerm 01/06, got %s", responses[1].InstallmentTerm)
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		transactions := []model.Transaction{}

		responses := toTransactionResponses(transactions)

		if len(responses) != 0 {
			t.Errorf("expected 0 responses, got %d", len(responses))
		}
	})

	t.Run("handles nil-like empty values", func(t *testing.T) {
		transactions := []model.Transaction{
			{
				UserID:          "",
				CardNumber:      "",
				TransactionDate: "",
				PostingDate:     "",
				Description:     "",
				Amount:          0,
				IsInstallment:   false,
				InstallmentTerm: "",
			},
		}

		responses := toTransactionResponses(transactions)

		if len(responses) != 1 {
			t.Fatalf("expected 1 response, got %d", len(responses))
		}

		if responses[0].UserID != "" {
			t.Errorf("expected empty UserID, got %s", responses[0].UserID)
		}
		if responses[0].Amount != 0 {
			t.Errorf("expected Amount 0, got %f", responses[0].Amount)
		}
	})

	t.Run("handles negative amounts", func(t *testing.T) {
		transactions := []model.Transaction{
			{
				Description: "REFUND",
				Amount:      -50.00,
			},
		}

		responses := toTransactionResponses(transactions)

		if responses[0].Amount != -50.00 {
			t.Errorf("expected Amount -50.00, got %f", responses[0].Amount)
		}
	})
}
