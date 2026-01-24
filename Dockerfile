# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Enable auto-download of required Go toolchain
ENV GOTOOLCHAIN=auto

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /helios ./cmd/main.go

# Runtime stage
FROM alpine:3.21

WORKDIR /app

# Install poppler-utils for pdftotext (includes Thai language support)
RUN apk add --no-cache poppler-utils

# Copy binary from builder
COPY --from=builder /helios .

# Expose port
EXPOSE 1323

# Run the application
CMD ["./helios"]
