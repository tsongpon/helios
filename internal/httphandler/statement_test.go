package httphandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/tsongpon/helios/internal/model"
)

type mockPDFService struct {
	transactions []model.Transaction
	err          error
}

func (m *mockPDFService) ExtractText(ctx context.Context, userID string, file io.Reader, password string) ([]model.Transaction, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.transactions, nil
}

func TestStatementHandler_CreateStatement(t *testing.T) {
	t.Run("returns transactions successfully", func(t *testing.T) {
		mockService := &mockPDFService{
			transactions: []model.Transaction{
				{
					UserID:          "1234567890",
					CardNumber:      "1234-XXXX-XXXX-5678",
					TransactionDate: "2024-12-15",
					PostingDate:     "2024-12-16",
					Description:     "AMAZON",
					Amount:          100.50,
					IsInstallment:   false,
				},
			},
		}

		handler := NewStatementHandler(mockService)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.pdf")
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		part.Write([]byte("fake pdf content"))
		writer.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/statements", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handler.CreateStatement(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var response []TransactionResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(response) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(response))
		}

		if response[0].Description != "AMAZON" {
			t.Errorf("expected description AMAZON, got %s", response[0].Description)
		}
	})

	t.Run("returns error when file is missing", func(t *testing.T) {
		mockService := &mockPDFService{}
		handler := NewStatementHandler(mockService)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/statements", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateStatement(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		var response ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error != "file is required" {
			t.Errorf("unexpected error message: %s", response.Error)
		}
	})

	t.Run("returns error when file is not PDF", func(t *testing.T) {
		mockService := &mockPDFService{}
		handler := NewStatementHandler(mockService)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.txt")
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		part.Write([]byte("fake text content"))
		writer.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/statements", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handler.CreateStatement(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
		}

		var response ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error != "only PDF files are allowed" {
			t.Errorf("unexpected error message: %s", response.Error)
		}
	})

	t.Run("returns error when PDF service fails", func(t *testing.T) {
		mockService := &mockPDFService{
			err: errors.New("failed to parse PDF"),
		}
		handler := NewStatementHandler(mockService)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.pdf")
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		part.Write([]byte("fake pdf content"))
		writer.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/statements", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handler.CreateStatement(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}

		var response ErrorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.Error != "failed to extract text from PDF: failed to parse PDF" {
			t.Errorf("unexpected error message: %s", response.Error)
		}
	})

	t.Run("accepts PDF with application/pdf content type", func(t *testing.T) {
		mockService := &mockPDFService{
			transactions: []model.Transaction{},
		}
		handler := NewStatementHandler(mockService)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		h := make(map[string][]string)
		h["Content-Disposition"] = []string{`form-data; name="file"; filename="test.pdf"`}
		h["Content-Type"] = []string{"application/pdf"}
		part, err := writer.CreatePart(h)
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		part.Write([]byte("fake pdf content"))
		writer.Close()

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/statements", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handler.CreateStatement(c)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})
}
