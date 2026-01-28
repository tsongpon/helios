package httphandler

import "github.com/tsongpon/helios/internal/model"

type StatementResponse struct {
	CardNumber     string                `json:"card_number"`
	TotalPayment   float64               `json:"total_payment"`
	MinimumPayment float64               `json:"minimum_payment"`
	PaymentDueDate string                `json:"payment_due_date"`
	CreditLine     float64               `json:"credit_line"`
	Transactions   []TransactionResponse `json:"transactions"`
}

type TransactionResponse struct {
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

func toStatementResponse(s model.Statement) StatementResponse {
	transactions := make([]TransactionResponse, len(s.Transactions))
	for i, t := range s.Transactions {
		transactions[i] = TransactionResponse{
			TransactionDate: t.TransactionDate,
			PostingDate:     t.PostingDate,
			Description:     t.Description,
			Amount:          t.Amount,
			IsInstallment:   t.IsInstallment,
			InstallmentTerm: t.InstallmentTerm,
		}
	}
	return StatementResponse{
		CardNumber:     s.CardNumber,
		TotalPayment:   s.TotalPayment,
		MinimumPayment: s.MinimumPayment,
		PaymentDueDate: s.PaymentDueDate,
		CreditLine:     s.CreditLine,
		Transactions:   transactions,
	}
}
