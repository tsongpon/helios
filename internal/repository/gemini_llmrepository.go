package repository

type GeminiLLMRepository struct {
	apiKey string
}

func NewGeminiLLMRepository(apiKey string) *GeminiLLMRepository {
	return &GeminiLLMRepository{
		apiKey: apiKey,
	}
}

func (r *GeminiLLMRepository) ParseStatement(statementText string) (string, error) {
	// Implementation goes here
	return "", nil
}
