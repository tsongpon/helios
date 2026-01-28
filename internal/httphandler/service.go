package httphandler

import (
	"context"
	"io"

	"github.com/tsongpon/helios/internal/model"
)

type PDFService interface {
	ExtractText(ctx context.Context, file io.Reader, password string) (model.Statement, error)
}
