# Helios

A PDF statement parsing service that extracts text from PDF bank statements and uses Google Gemini LLM to parse transactions into structured data via a REST API.

## Features

- Extract text from PDF files
- Parse bank statement transactions using Google Gemini LLM
- Support for password-protected PDFs
- Unicode character support (including Thai language)
- Simple REST API interface
- Docker support for easy deployment

## Prerequisites

- Go 1.25.1 or higher
- Poppler utilities (pdftotext) with Thai language support
- Google Gemini API key
- Docker & Docker Compose (optional)

## Installation

### Local Setup

```bash
# Clone the repository
git clone https://github.com/tsongpon/helios.git
cd helios

# Download dependencies
go mod download

# Build the application
go build -o helios ./cmd/main.go
```

### Docker Setup

```bash
# Build and run with Docker Compose
docker-compose up --build

# Or build manually
docker build -t helios .
docker run -p 1323:1323 -e GEMINI_API_KEY=your_api_key helios
```

## Configuration

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| GEMINI_API_KEY | Yes | Google Gemini API key for LLM parsing |

### Getting a Gemini API Key

1. Go to [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Sign in with your Google account
3. Click "Create API Key"
4. Copy the generated API key

## Usage

Start the server:

```bash
# Local - set the API key and run
export GEMINI_API_KEY=your_api_key
./helios

# Docker
docker-compose up
```

The server runs on port **1323**.

## API Endpoints

### Health Check

```
GET /ping
```

Response: `pong`

### Parse Bank Statement

```
POST /statements
Content-Type: multipart/form-data
```

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| file | file | Yes | PDF bank statement file to upload |
| password | string | No | Password for protected PDFs |

**Response:**
```json
{
  "transactions": [
    {
      "transaction_date": "2024-01-15",
      "posting_date": "2024-01-15",
      "description": "TRANSFER TO SAVINGS",
      "amount": -500.00
    },
    {
      "transaction_date": "2024-01-16",
      "posting_date": "2024-01-16",
      "description": "SALARY DEPOSIT",
      "amount": 3000.00
    }
  ]
}
```

**Examples:**

```bash
# Parse bank statement
curl -X POST -F "file=@statement.pdf" http://localhost:1323/statements

# Password-protected PDF
curl -X POST -F "file=@statement.pdf" -F "password=secret" http://localhost:1323/statements
```

## Project Structure

```
helios/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── httphandler/         # HTTP request handlers
│   ├── model/               # Data models (Transaction)
│   ├── service/             # Business logic (PDF extraction)
│   └── repository/          # Gemini LLM integration
├── Dockerfile               # Docker build configuration
├── docker-compose.yml       # Docker Compose orchestration
├── go.mod                   # Go module definition
└── go.sum                   # Dependencies lock file
```

## Tech Stack

- **Go 1.25.1** - Primary language
- **Echo v5** - Web framework
- **pdftotext** (Poppler) - PDF text extraction
- **Google Gemini API** - LLM for transaction parsing
