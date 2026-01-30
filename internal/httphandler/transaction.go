package httphandler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type TransactionHandler struct {
	transactionService TransactionService
}

func NewTransactionHandler(transactionService TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) GetTransactions(c *echo.Context) error {
	startDate := c.QueryParam("start")
	endDate := c.QueryParam("end")

	if startDate == "" || endDate == "" {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "start and end query parameters are required (format: YYYY-MM-DD)",
		})
	}

	from, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid start date format, expected YYYY-MM-DD",
		})
	}

	to, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid end date format, expected YYYY-MM-DD",
		})
	}

	// Fix userID for now
	userID := "1234567890"

	transactions, err := h.transactionService.GetTransactions(c.Request().Context(), userID, from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "failed to get transactions: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, toTransactionResponses(transactions))
}
