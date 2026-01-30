package httphandler

import "github.com/tsongpon/helios/internal/model"

type TransactionResponse struct {
	CardNumber      string  `json:"card_number"`
	UserID          string  `json:"user_id"`
	TransactionDate string  `json:"transaction_date"`
	PostingDate     string  `json:"posting_date"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	IsInstallment   bool    `json:"is_installment"`
	InstallmentTerm string  `json:"installment_term"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func toTransactionResponses(transactions []model.Transaction) []TransactionResponse {
	responses := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		responses[i] = TransactionResponse{
			CardNumber:      t.CardNumber,
			UserID:          t.UserID,
			TransactionDate: t.TransactionDate,
			PostingDate:     t.PostingDate,
			Description:     t.Description,
			Amount:          t.Amount,
			IsInstallment:   t.IsInstallment,
			InstallmentTerm: t.InstallmentTerm,
		}
	}
	return responses
}
