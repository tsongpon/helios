package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tsongpon/helios/internal/model"
)

type GeminiLLMRepository struct {
	apiKey string
}

func NewGeminiLLMRepository(apiKey string) *GeminiLLMRepository {
	return &GeminiLLMRepository{
		apiKey: apiKey,
	}
}

// Gemini API request/response structures
type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (r *GeminiLLMRepository) ParseStatement(statementText string) (model.Statement, error) {
	prompt := fmt.Sprintf(`Parse the following bank statement text and extract statement information and all transactions.

First, output the statement header information in the following format on the FIRST line:
HEADER|card_number|total_payment|minimum_payment|payment_due_date|credit_line

Rules for header:
- card_number: The credit card number (may be partially masked, e.g., "1234-56XX-XXXX-7890")
- total_payment: Total payment amount as a number
- minimum_payment: Minimum payment amount as a number
- payment_due_date: Payment due date in format YYYY-MM-DD
- credit_line: Credit line/limit as a number
- If any field is not found, use empty string for text fields and 0 for numeric fields

Then, output each transaction on a separate line in pipe-delimited format:
transaction_date|posting_date|description|amount|is_installment|installment_term

Rules:
- transaction_date: The date the transaction occurred (format: YYYY-MM-DD)
- posting_date: The date the transaction was posted (format: YYYY-MM-DD), use transaction_date if not available
- description: The transaction description (remove extra whitespace)
- amount: The transaction amount as a number. Use NEGATIVE values for credits/refunds/payments (marked with "CR" suffix or with "-" suffix). Use POSITIVE values for purchases/charges.
  - Example: "31,751.00 CR" should be -31751.00 (credit/payment)
  - Example: "14.20-" should be -14.20 (payment/credit)
  - Example: "1,070.00" should be 1070.00 (purchase/charge)
- is_installment: "true" if this is an installment transaction, "false" otherwise
- installment_term: For installment transactions, the term indicator (e.g., "009/010" means 9th payment of 10 total). Empty string for non-installment transactions.

For installment transactions:
- They may appear in a separate "Installment" section OR inline with the description
- The installment term format is like "009/010" or "04/06" or "10/10" indicating current term / total terms
- Inline format: The term may appear right after the merchant name, e.g., "ZOOM CAMERA-WEST GATE 10/10" or "2C2P *LAZADA 04/06"
  - In this case, extract the term (e.g., "10/10", "04/06") as installment_term
  - Use the rightmost amount as the transaction amount
- Separate section format: if a line shows "13,281.00  009/010  6,640.50", use 6,640.50 as the amount
- ANY transaction with a term pattern like "NN/NN" (digits/digits) should be marked as is_installment=true

Determining the year for transactions:
- If transaction dates only show day/month (e.g., "17/12" or "17/12/"), look for the PAYMENT DATE in the statement header to determine the year
- The PAYMENT DATE is usually in format "DD/MM/YY" (e.g., "06/02/25" means 2025)
- If transaction month is greater than payment month, the transaction year is the previous year
- Example: Payment date is 06/02/25, transaction date 17/12 means 2024-12-17; transaction date 05/01 means 2025-01-05

Output ONLY the HEADER line followed by the pipe-delimited transaction lines, no other headers or extra text.

Bank Statement Text:
%s`, statementText)

	req := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return model.Statement{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", r.apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return model.Statement{}, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return model.Statement{}, fmt.Errorf("failed to send request to Gemini API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Statement{}, fmt.Errorf("Gemini API returned status %d", resp.StatusCode)
	}

	var geminiResp geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return model.Statement{}, fmt.Errorf("failed to decode Gemini response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return model.Statement{}, fmt.Errorf("no response from Gemini API")
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	return parsePipeDelimitedResponse(responseText)
}

func parsePipeDelimitedResponse(text string) (model.Statement, error) {
	var statement model.Statement
	lines := strings.Split(strings.TrimSpace(text), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")

		// Parse header line
		if len(parts) == 6 && strings.TrimSpace(parts[0]) == "HEADER" {
			statement.CardNumber = strings.TrimSpace(parts[1])
			if v, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64); err == nil {
				statement.TotalPayment = v
			}
			if v, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64); err == nil {
				statement.MinimumPayment = v
			}
			statement.PaymentDueDate = strings.TrimSpace(parts[4])
			if v, err := strconv.ParseFloat(strings.TrimSpace(parts[5]), 64); err == nil {
				statement.CreditLine = v
			}
			continue
		}

		// Parse transaction line
		if len(parts) != 6 {
			continue // Skip malformed lines
		}

		amount, err := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
		if err != nil {
			continue // Skip lines with invalid amounts
		}

		isInstallment := strings.TrimSpace(parts[4]) == "true"
		installmentTerm := strings.TrimSpace(parts[5])

		transaction := model.Transaction{
			TransactionDate: strings.TrimSpace(parts[0]),
			PostingDate:     strings.TrimSpace(parts[1]),
			Description:     strings.TrimSpace(parts[2]),
			Amount:          amount,
			IsInstallment:   isInstallment,
			InstallmentTerm: installmentTerm,
		}
		statement.Transactions = append(statement.Transactions, transaction)
	}

	return statement, nil
}
