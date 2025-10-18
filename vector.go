package main

import (
	"context"
	"fmt"
	"os"
	"prj-start/config"
	"prj-start/document"
	"prj-start/logger"
	"prj-start/vector"
	"time"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	logger.LogInfo("Starting vector upsert process")

	// Load configuration
	cfg := config.LoadUpstashConfig()
	if !cfg.Validate() {
		logger.LogError("Upstash configuration is invalid. Please set UPSTASH_VECTOR_URL and UPSTASH_VECTOR_TOKEN environment variables.")
		os.Exit(1)
	}

	// Initialize vector client
	vectorClient, err := vector.NewClient(cfg)
	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to create vector client: %v", err))
		os.Exit(1)
	}

	// Initialize document reader
	reader := document.NewReader("dev-docs")

	// Read all documents
	logger.LogInfo("Reading documents from dev-docs folder...")
	documents, err := reader.ReadAllDocuments()
	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to read documents: %v", err))
		os.Exit(1)
	}

	if len(documents) == 0 {
		logger.LogWarning("No documents found to process")
		return
	}

	logger.LogSuccess(fmt.Sprintf("Found %d documents to process", len(documents)))

	// Initialize upserter
	upserter := vector.NewUpserter(vectorClient, 10) // Batch size of 10

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Process documents
	startTime := time.Now()
	err = upserter.UpsertAllDocuments(ctx, documents)
	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to upsert documents: %v", err))
		os.Exit(1)
	}

	duration := time.Since(startTime)
	logger.LogSuccess(fmt.Sprintf("Vector upsert completed successfully in %v", duration))
}
