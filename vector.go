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
		logger.LogError("Upstash configuration is invalid. Please set UPSTASH_VECTOR_REST_URL and UPSTASH_VECTOR_REST_TOKEN environment variables.")
		os.Exit(1)
	}

	// Log configuration (without sensitive data)
	logger.LogInfo(fmt.Sprintf("Configuration loaded - Batch size: %d, Timeout: %d minutes, Log level: %s", cfg.BatchSize, cfg.ProcessingTimeout, cfg.LogLevel))
	if cfg.HasMCPConfig() {
		logger.LogInfo("MCP configuration found for querying")
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

	// Initialize upserter with configured batch size
	upserter := vector.NewUpserter(vectorClient, cfg.BatchSize)

	// Create context with configured timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ProcessingTimeout)*time.Minute)
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
