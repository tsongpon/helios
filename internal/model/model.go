package model

type Statement struct {
	CardNumber     string
	TotalPayment   float64
	MinimumPayment float64
	PaymentDueDate string
	CreditLine     float64
	Transactions   []Transaction
}

type Transaction struct {
	TransactionDate string
	PostingDate     string
	Description     string
	Amount          float64
	IsInstallment   bool
	InstallmentTerm string
}
