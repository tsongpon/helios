package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// PDFService handles PDF text extraction
type PDFService struct {
	llmReposiroty LLMRepository
}

// NewPDFService creates a new PDFService instance
func NewPDFService(llmRepository LLMRepository) *PDFService {
	return &PDFService{
		llmReposiroty: llmRepository,
	}
}

// ExtractText extracts text content from a PDF file using pdftotext
// password is optional - pass empty string for non-protected PDFs
func (s *PDFService) ExtractText(ctx context.Context, file io.Reader, password string) (string, error) {
	// Create a temporary file to store the PDF
	tmpFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write the uploaded content to temp file
	if _, err := io.Copy(tmpFile, file); err != nil {
		return "", err
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
			return "", fmt.Errorf("%s", strings.TrimSpace(stderr.String()))
		}
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}
