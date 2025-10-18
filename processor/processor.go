package processor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/typicalfo/prj-start/config"
	"github.com/typicalfo/prj-start/document"
	"github.com/typicalfo/prj-start/logger"
	"github.com/typicalfo/prj-start/vector"
)

// ProcessFolder processes all documents in a folder and upserts them to Upstash Vector
func ProcessFolder(ctx context.Context, cfg *config.Config, folderPath string) error {
	logger.LogInfo("Starting document processing")
	logger.LogInfo("Configuration loaded successfully")
	logger.LogInfo(fmt.Sprintf("Default namespace: %s", cfg.DefaultNamespace))
	logger.LogInfo(fmt.Sprintf("Batch size: %d", cfg.BatchSize))

	// Initialize Upstash Vector client
	client, err := vector.NewClient(&cfg.Upstash)
	if err != nil {
		return fmt.Errorf("failed to create Upstash client: %w", err)
	}

	// Read all documents from the folder
	logger.LogInfo(fmt.Sprintf("Scanning folder: %s", folderPath))
	reader := document.NewReader(folderPath)
	documents, err := reader.ReadAllDocuments()
	if err != nil {
		return fmt.Errorf("failed to read documents: %w", err)
	}

	if len(documents) == 0 {
		logger.LogWarning("No documents found to process")
		return nil
	}

	logger.LogSuccess(fmt.Sprintf("Found %d documents to process", len(documents)))

	// Create upserter
	upserter := vector.NewUpserter(client, cfg.BatchSize)

	// Process all documents
	startTime := time.Now()
	err = upserter.UpsertAllDocuments(ctx, documents)
	if err != nil {
		return fmt.Errorf("failed to upsert documents: %w", err)
	}

	duration := time.Since(startTime)
	logger.LogSuccess(fmt.Sprintf("Processing completed in %v", duration))

	// List available namespaces
	namespaces, err := client.ListNamespaces(ctx)
	if err == nil {
		logger.LogInfo("Available namespaces:")
		for _, ns := range namespaces {
			logger.LogInfo(fmt.Sprintf("  - %s", ns))
		}
	}

	return nil
}

// ValidateFolder checks if the folder exists and is readable
func ValidateFolder(folderPath string) error {
	// Check if folder exists
	info, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("folder '%s' does not exist", folderPath)
	}
	if err != nil {
		return fmt.Errorf("failed to access folder '%s': %w", folderPath, err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("'%s' is not a directory", folderPath)
	}

	// Check if it's readable
	file, err := os.Open(folderPath)
	if err != nil {
		return fmt.Errorf("folder '%s' is not readable: %w", folderPath, err)
	}
	file.Close()

	// Convert to absolute path
	absPath, err := filepath.Abs(folderPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	logger.LogInfo(fmt.Sprintf("Processing folder: %s", absPath))
	return nil
}
