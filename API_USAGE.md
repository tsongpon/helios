# Statement Extraction API - Usage Guide

## Overview

This API extracts credit card statement data from uploaded PDF files using AI-powered parsing (Google Gemini LLM via LangChain4j).

## Architecture

1. **Controller** (`StatementController`) - Receives PDF uploads, extracts text using Apache PDFBox
2. **Service** (`StatementService`) - Orchestrates the parsing flow
3. **Repository** (`GeminiLLMRepository`) - Uses LangChain4j to call Gemini LLM for parsing statement text to structured data
4. **Model** (`Statement`, `Transaction`) - Data classes representing the parsed statement

## Setup

### 1. Configure Google Cloud Credentials

Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to point to your GCP service account key:

```bash
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/your-service-account-key.json"
```

### 2. Configure Application Properties

Set the following environment variables or update `application.yaml`:

```bash
export GEMINI_PROJECT_ID="your-gcp-project-id"
export GEMINI_LOCATION="us-central1"  # Optional, defaults to us-central1
export GEMINI_MODEL_NAME="gemini-1.5-flash"  # Optional, defaults to gemini-1.5-flash
```

### 3. Build the Application

```bash
mvn clean install
```

### 4. Run the Application

```bash
mvn spring-boot:run
```

The API will be available at `http://localhost:8080`

## API Endpoints

### 1. Extract Statement (Full Parsing)

**Endpoint:** `POST /v1/statements/extract`

**Description:** Uploads a PDF statement, extracts text, and parses it into a structured `Statement` object.

**Request:**
- Content-Type: `multipart/form-data`
- Parameters:
  - `file` (required): PDF file
  - `password` (optional): PDF password if encrypted

**Response:**
- Content-Type: `application/json`
- Body: `Statement` object with all parsed data

**Example:**

```bash
curl -X POST http://localhost:8080/v1/statements/extract \
  -F "file=@/path/to/statement.pdf"
```

**Sample Response:**

```json
{
  "cardNumber": "5468 48XX XXXX 8032",
  "statementDate": "2025-08-05",
  "creditLine": 600000.00,
  "totalPaymentDue": 21367.39,
  "outstandingBalance": 38344.01,
  "minimumPaymentDue": 10942.35,
  "availableCredit": 560965.99,
  "availableCreditLimit": 600000.00,
  "paymentDueDate": "2025-08-25",
  "lastPaymentDate": "2025-07-28",
  "transactions": [
    {
      "transactionDate": "2025-07-05",
      "postingDate": "2025-07-06",
      "description": "OMISE.CO BANGKOK",
      "amount": 690.00,
      "isInstallment": false,
      "installmentCurrent": null,
      "installmentTotal": null,
      "installmentPlanAmount": null
    },
    {
      "transactionDate": "2025-07-06",
      "postingDate": "2025-07-07",
      "description": "WWW.GRAB.COM BANGKOK THA",
      "amount": 251.00,
      "isInstallment": false,
      "installmentCurrent": null,
      "installmentTotal": null,
      "installmentPlanAmount": null
    }
  ]
}
```

### 2. Extract Text Only

**Endpoint:** `POST /v1/statements/extract-text`

**Description:** Uploads a PDF and returns the extracted raw text without parsing.

**Request:**
- Content-Type: `multipart/form-data`
- Parameters:
  - `file` (required): PDF file
  - `password` (optional): PDF password if encrypted

**Response:**
- Content-Type: `text/plain`
- Body: Extracted text from PDF

**Example:**

```bash
curl -X POST http://localhost:8080/v1/statements/extract-text \
  -F "file=@/path/to/statement.pdf"
```

## Error Responses

All errors return a JSON error response:

```json
{
  "message": "Error description"
}
```

**Common Error Scenarios:**
- 400 Bad Request: No file uploaded, invalid file type, or empty statement text
- 500 Internal Server Error: PDF processing error, LLM parsing failure

## Data Model

### Statement

```java
public record Statement(
    String cardNumber,              // Partially masked card number
    LocalDate statementDate,        // Statement generation date
    Double creditLine,              // Total credit limit
    Double totalPaymentDue,         // Total amount due
    Double outstandingBalance,      // Current outstanding balance
    Double minimumPaymentDue,       // Minimum payment required
    Double availableCredit,         // Available credit amount
    Double availableCreditLimit,    // Available credit limit
    LocalDate paymentDueDate,       // Payment deadline
    LocalDate lastPaymentDate,      // Last payment received date
    List<Transaction> transactions  // List of transactions
)
```

### Transaction

```java
public record Transaction(
    LocalDate transactionDate,      // Transaction date
    LocalDate postingDate,          // Posting date
    String description,             // Transaction description
    Double amount,                  // Amount (positive for charges, negative for payments)
    Boolean isInstallment,          // Whether this is an installment
    Integer installmentCurrent,     // Current installment number
    Integer installmentTotal,       // Total installments
    Double installmentPlanAmount    // Total installment plan amount
)
```

## Testing with Example Statement

A sample statement PDF is included at:
```
src/main/resources/statement.pdf
```

Test the API:

```bash
curl -X POST http://localhost:8080/v1/statements/extract \
  -F "file=@src/main/resources/statement.pdf" \
  | jq .
```

## Configuration Reference

### application.yaml

```yaml
spring:
  application:
    name: helios
  servlet:
    multipart:
      max-file-size: 10MB
      max-request-size: 10MB

gemini:
  project-id: ${GEMINI_PROJECT_ID:your-gcp-project-id}
  location: ${GEMINI_LOCATION:us-central1}
  model-name: ${GEMINI_MODEL_NAME:gemini-1.5-flash}
```

## Dependencies

- Spring Boot 4.0.1
- Apache PDFBox 3.0.3
- LangChain4j 0.36.2
- LangChain4j Vertex AI Gemini 0.36.2
- Google Cloud AI Platform 3.58.0
- Gson 2.11.0

## Troubleshooting

### Authentication Error

If you get authentication errors, ensure:
1. `GOOGLE_APPLICATION_CREDENTIALS` is set correctly
2. The service account has Vertex AI user permissions
3. Vertex AI API is enabled in your GCP project

### PDF Processing Error

If PDF extraction fails:
1. Check that the file is a valid PDF
2. If encrypted, provide the password parameter
3. Check file size is under 10MB

### LLM Parsing Error

If parsing fails:
1. Check Gemini API quota limits
2. Verify the statement format matches expected structure
3. Review LLM response in error logs
