package service

import "github.com/tsongpon/helios/internal/model"

type LLMRepository interface {
	ParseStatement(statementText string) ([]model.Transaction, error)
}
