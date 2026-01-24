# Helios

A PDF text extraction service that converts PDF documents into plain text format via a REST API.

## Features

- Extract text from PDF files
- Support for password-protected PDFs
- Unicode character support (including Thai language)
- Simple REST API interface
- Docker support for easy deployment

## Prerequisites

- Go 1.25.1 or higher
- Poppler utilities (pdftotext) with Thai language support
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
docker run -p 1323:1323 helios
```

## Usage

Start the server:

```bash
# Local
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

### Extract Text from PDF

```
POST /statements
Content-Type: multipart/form-data
```

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| file | file | Yes | PDF file to upload |
| password | string | No | Password for protected PDFs |

**Response:**
```json
{
  "text": "extracted text content..."
}
```

**Examples:**

```bash
# Simple PDF extraction
curl -X POST -F "file=@document.pdf" http://localhost:1323/statements

# Password-protected PDF
curl -X POST -F "file=@document.pdf" -F "password=secret" http://localhost:1323/statements
```

## Project Structure

```
helios/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── httphandler/         # HTTP request handlers
│   ├── model/               # Data models
│   ├── service/             # Business logic
│   └── repository/          # Data layer (placeholder)
├── Dockerfile               # Docker build configuration
├── docker-compose.yml       # Docker Compose orchestration
├── go.mod                   # Go module definition
└── go.sum                   # Dependencies lock file
```

## Tech Stack

- **Go 1.25.1** - Primary language
- **Echo v5** - Web framework
- **pdftotext** (Poppler) - PDF text extraction
