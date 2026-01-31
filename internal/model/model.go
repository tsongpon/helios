package model

type Transaction struct {
	ID              string
	UserID          string
	CardNumber      string
	TransactionDate string
	PostingDate     string
	Description     string
	Amount          float64
	IsInstallment   bool
	InstallmentTerm string
}
