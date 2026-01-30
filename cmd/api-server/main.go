package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tsongpon/helios/internal/httphandler"
	"github.com/tsongpon/helios/internal/repository"
	"github.com/tsongpon/helios/internal/service"
)

func main() {
	godotenv.Load()

	gcpProjectID := os.Getenv("GCP_PROJECT_ID")
	databaseID := os.Getenv("GCP_FIRESTORE_DATABASE_ID")
	ctx := context.Background()
	firestoreClient, err := firestore.NewClientWithDatabase(ctx, gcpProjectID, databaseID)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	llmAPIKey := os.Getenv("GEMINI_API_KEY")
	llmRepository := repository.NewGeminiLLMRepository(llmAPIKey)
	transactionRepository := repository.NewFirestoreTransactionRepository(firestoreClient)

	pdfService := service.NewPDFService(llmRepository, transactionRepository)
	transactionService := service.NewTransactionService(transactionRepository)

	pingHandler := httphandler.NewPingHandler()
	statementHandler := httphandler.NewStatementHandler(pdfService)
	transactionHandler := httphandler.NewTransactionHandler(transactionService)

	e := echo.New()
	e.Use(middleware.RequestLogger())

	e.GET("/ping", pingHandler.Ping)
	e.POST("/statements", statementHandler.CreateStatement)
	e.GET("/transactions", transactionHandler.GetTransactions)

	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
