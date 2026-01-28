package model

type Statement struct {
	CardNumber     string        `json:"card_number"`
	TotalPayment   float64       `json:"total_payment"`
	MinimumPayment float64       `json:"minimum_payment"`
	PaymentDueDate string        `json:"payment_due_date"`
	CreditLine     float64       `json:"credit_line"`
	Transactions   []Transaction `json:"transactions"`
}

type Transaction struct {
	TransactionDate string  `json:"transaction_date"`
	PostingDate     string  `json:"posting_date"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
	IsInstallment   bool    `json:"is_installment"`
	InstallmentTerm string  `json:"installment_term"`
}
