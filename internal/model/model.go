package model

type Transaction struct {
	TransactionDate string  `json:"transaction_date"`
	PostingDate     string  `json:"posting_date"`
	Description     string  `json:"description"`
	Amount          float64 `json:"amount"`
}
