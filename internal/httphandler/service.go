package httphandler

import (
	"context"
	"io"
)

type PDFService interface {
	ExtractText(ctx context.Context, file io.Reader, password string) (string, error)
}
