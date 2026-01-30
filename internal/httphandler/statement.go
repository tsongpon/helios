package httphandler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type StatementHandler struct {
	pdfService PDFService
}

func NewStatementHandler(pdfService PDFService) *StatementHandler {
	return &StatementHandler{
		pdfService: pdfService,
	}
}

func (h *StatementHandler) CreateStatement(c *echo.Context) error {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "file is required",
		})
	}

	// Validate file type
	if file.Header.Get("Content-Type") != "application/pdf" {
		// Also check file extension as fallback
		if len(file.Filename) < 4 || file.Filename[len(file.Filename)-4:] != ".pdf" {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "only PDF files are allowed",
			})
		}
	}

	// Get optional password for protected PDFs
	password := c.FormValue("password")

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to open uploaded file",
		})
	}
	defer src.Close()

	// Fix userID for now
	// TODO: Implement user authentication and authorization
	userID := "1234567890"

	// Extract text from PDF and parse transactions
	transactions, err := h.pdfService.ExtractText(c.Request().Context(), userID, src, password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to extract text from PDF: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, toTransactionResponses(transactions))
}
