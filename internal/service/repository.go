package service

import "github.com/tsongpon/helios/internal/model"

type LLMRepository interface {
	ParseStatement(statementText string) (model.Statement, error)
}

type StatementRepository interface {
	Save(statement model.Statement) error
}

type TransactionRepository interface {
	Save(transaction []model.Transaction) error
}
