package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/tsongpon/helios/internal/model"
)

// PDFService handles PDF text extraction and parsing
type PDFService struct {
	llmRepository         LLMRepository
	transactionRepository TransactionRepository
}

// NewPDFService creates a new PDFService instance
func NewPDFService(llmRepository LLMRepository, transactionRepository TransactionRepository) *PDFService {
	return &PDFService{
		llmRepository:         llmRepository,
		transactionRepository: transactionRepository,
	}
}

// ExtractText extracts text content from a PDF file using pdftotext
// password is optional - pass empty string for non-protected PDFs
func (s *PDFService) ExtractText(ctx context.Context, userID string, file io.Reader, password string) ([]model.Transaction, error) {
	// Create a temporary file to store the PDF
	tmpFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write the uploaded content to temp file
	if _, err := io.Copy(tmpFile, file); err != nil {
		return nil, err
	}
	tmpFile.Close()

	// Build pdftotext command arguments
	args := []string{"-layout"}
	if password != "" {
		args = append(args, "-upw", password)
	}
	args = append(args, tmpFile.Name(), "-")

	// Use pdftotext to extract text (supports Thai and other Unicode)
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("pdftotext", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return nil, fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
		}
		return nil, err
	}

	extractedText := strings.TrimSpace(stdout.String())

	// Send extracted text to LLM repository for parsing
	transactions, err := s.llmRepository.ParseStatement(extractedText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse statement: %w", err)
	}
	for i := range transactions {
		transactions[i].UserID = userID
	}

	if err := s.transactionRepository.Save(ctx, transactions); err != nil {
		return nil, fmt.Errorf("failed to save transactions: %w", err)
	}

	return transactions, nil
}
