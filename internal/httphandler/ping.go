package httphandler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) Ping(c *echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
