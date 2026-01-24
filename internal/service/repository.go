package service

type LLMRepository interface {
	ParseStatement(statementText string) (string, error)
}
