package main

import (
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tsongpon/helios/internal/httphandler"
	"github.com/tsongpon/helios/internal/repository"
	"github.com/tsongpon/helios/internal/service"
)

func main() {
	llmAPIKey := os.Getenv("GEMINI_API_KEY")
	llmRepository := repository.NewGeminiLLMRepository(llmAPIKey)

	pdfService := service.NewPDFService(llmRepository)

	pingHandler := httphandler.NewPingHandler()
	statementHandler := httphandler.NewStatementHandler(pdfService)

	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/ping", pingHandler.Ping)
	e.POST("/statements", statementHandler.CreateStatement)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
