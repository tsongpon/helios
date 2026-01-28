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
	llmRepository       LLMRepository
	statementRepository StatementRepository
}

// NewPDFService creates a new PDFService instance
func NewPDFService(llmRepository LLMRepository, statementRepository StatementRepository) *PDFService {
	return &PDFService{
		llmRepository:       llmRepository,
		statementRepository: statementRepository,
	}
}

// ExtractText extracts text content from a PDF file using pdftotext
// password is optional - pass empty string for non-protected PDFs
func (s *PDFService) ExtractText(ctx context.Context, file io.Reader, password string) (model.Statement, error) {
	// Create a temporary file to store the PDF
	tmpFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return model.Statement{}, err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write the uploaded content to temp file
	if _, err := io.Copy(tmpFile, file); err != nil {
		return model.Statement{}, err
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
			return model.Statement{}, fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
		}
		return model.Statement{}, err
	}

	extractedText := strings.TrimSpace(stdout.String())

	// Send extracted text to LLM repository for parsing
	statement, err := s.llmRepository.ParseStatement(extractedText)
	if err != nil {
		return model.Statement{}, fmt.Errorf("failed to parse statement: %w", err)
	}

	if err := s.statementRepository.Save(statement); err != nil {
		return model.Statement{}, fmt.Errorf("failed to save statement: %w", err)
	}

	return statement, nil
}
