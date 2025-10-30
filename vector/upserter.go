package vector

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/typicalfo/prj-start/document"
	"github.com/typicalfo/prj-start/logger"
	"path/filepath"
	"strings"
)

type Upserter struct {
	client    *Client
	batchSize int
}

func NewUpserter(client *Client, batchSize int) *Upserter {
	if batchSize <= 0 {
		batchSize = 10 // default batch size
	}
	return &Upserter{
		client:    client,
		batchSize: batchSize,
	}
}

func (u *Upserter) UpsertAllDocuments(ctx context.Context, documents []document.FileInfo) error {
	logger.LogInfo(fmt.Sprintf("Starting upsert of %d documents", len(documents)))

	totalChunks := 0
	processedChunks := 0
	failedDocuments := 0

	// First pass: count total chunks for progress tracking
	for _, doc := range documents {
		chunker := document.NewChunker(1000)
		chunks, err := chunker.ChunkDocument(doc)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error chunking document %s: %v", doc.RelativePath, err))
			failedDocuments++
			continue
		}
		totalChunks += len(chunks)
	}

	logger.LogInfo(fmt.Sprintf("Total chunks to process: %d", totalChunks))

	// Second pass: process documents by namespace
	namespaces := make(map[string][]Document)

	for i, doc := range documents {
		logger.LogProgress(i+1, len(documents), fmt.Sprintf("Processing document: %s", doc.RelativePath))

		// Extract namespace from file path (exclude dev-docs)
		namespace := u.extractNamespace(doc.RelativePath)

		chunker := document.NewChunker(1000)
		chunks, err := chunker.ChunkDocument(doc)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error chunking document %s: %v", doc.RelativePath, err))
			failedDocuments++
			continue
		}

		// Convert chunks to documents
		for _, chunk := range chunks {
			docID := u.generateDocumentID(doc.RelativePath, chunk.Index)

			// Prepare metadata
			metadata := make(map[string]string)
			for k, v := range chunk.Metadata {
				metadata[k] = v
			}
			metadata["chunk_index"] = fmt.Sprintf("%d", chunk.Index)
			metadata["source_file"] = doc.RelativePath
			metadata["file_size"] = fmt.Sprintf("%d", doc.Size)

			// Add recipe/project information
			fullPath := u.extractFullPath(doc.RelativePath)
			metadata["namespace"] = namespace
			metadata["full_path"] = fullPath
			metadata["recipe_name"] = u.extractRecipeName(doc.RelativePath)
			metadata["project_type"] = u.extractProjectType(doc.RelativePath)

			namespaces[namespace] = append(namespaces[namespace], Document{
				ID:        docID,
				Content:   chunk.Content,
				Metadata:  metadata,
				Namespace: namespace,
			})
		}
	}

	// Process batches by namespace
	for namespace, docs := range namespaces {
		logger.LogInfo(fmt.Sprintf("Processing %d documents for namespace: %s", len(docs), namespace))

		batch := make([]Document, 0, u.batchSize)
		for _, doc := range docs {
			batch = append(batch, doc)

			// Process batch when it reaches the size limit
			if len(batch) >= u.batchSize {
				err := u.client.UpsertBatch(ctx, batch, namespace)
				if err != nil {
					logger.LogError(fmt.Sprintf("Error upserting batch for namespace %s: %v", namespace, err))
					return err
				}
				processedChunks += len(batch)
				logger.LogProgress(processedChunks, totalChunks, "Chunks processed")
				batch = batch[:0] // Clear batch
			}
		}

		// Process remaining items in batch for this namespace
		if len(batch) > 0 {
			err := u.client.UpsertBatch(ctx, batch, namespace)
			if err != nil {
				logger.LogError(fmt.Sprintf("Error upserting final batch for namespace %s: %v", namespace, err))
				return err
			}
			processedChunks += len(batch)
			logger.LogProgress(processedChunks, totalChunks, "Final chunks processed")
		}
	}

	logger.LogSuccess(fmt.Sprintf("Upsert completed! Processed %d chunks from %d documents across %d namespaces", processedChunks, len(documents)-failedDocuments, len(namespaces)))
	if failedDocuments > 0 {
		logger.LogWarning(fmt.Sprintf("Failed to process %d documents", failedDocuments))
	}

	return nil
}

func (u *Upserter) extractNamespace(relativePath string) string {
	// Extract namespace as recipe name (last directory in path)
	// Upstash Vector doesn't support nested namespaces with slashes
	// Format: "go-fiber-recipes/404-handler/main.go" -> "404-handler"
	// Format: "clean-architecture/api/handlers/book_handler.go" -> "handlers"
	// Format: "clean-code/app/server/domain/books.go" -> "server"

	parts := strings.Split(relativePath, string(filepath.Separator))
	if len(parts) >= 2 {
		// Get the immediate parent directory of the file
		parentDir := parts[len(parts)-2]

		// Sanitize namespace: remove leading dots and special characters
		// Upstash namespaces must start with alphanumeric character
		if strings.HasPrefix(parentDir, ".") {
			// Remove leading dots and replace with "hidden-" prefix
			sanitized := strings.TrimLeft(parentDir, ".")
			if sanitized == "" {
				return "hidden"
			}
			return "hidden-" + sanitized
		}

		// Skip common system directories that shouldn't be namespaces
		systemDirs := map[string]bool{
			"git":          true,
			"vscode":       true,
			"idea":         true,
			"node_modules": true,
			"vendor":       true,
			"dist":         true,
			"build":        true,
		}

		if systemDirs[parentDir] {
			return "system"
		}

		return parentDir
	}
	return "default"
}

func (u *Upserter) extractRecipeName(relativePath string) string {
	// Extract the recipe name (last directory in path)
	// Format: "go-fiber-recipes/404-handler/main.go" -> "404-handler"
	// Format: "clean-architecture/api/handlers/book_handler.go" -> "handlers"

	parts := strings.Split(relativePath, string(filepath.Separator))
	if len(parts) >= 2 {
		// Get the last directory (excluding filename)
		lastDir := parts[len(parts)-2]
		return lastDir
	}
	return "root"
}

func (u *Upserter) extractProjectType(relativePath string) string {
	// Extract top-level project type
	// Format: "go-fiber-recipes/404-handler/main.go" -> "go-fiber-recipes"
	// Format: "clean-architecture/api/handlers/book_handler.go" -> "clean-architecture"

	parts := strings.Split(relativePath, string(filepath.Separator))
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return "unknown"
}

func (u *Upserter) extractFullPath(relativePath string) string {
	// Extract full directory path (excluding filename)
	// Format: "go-fiber-recipes/404-handler/main.go" -> "go-fiber-recipes/404-handler"
	// Format: "clean-code/app/server/domain/books.go" -> "clean-code/app/server"

	dirPath := filepath.Dir(relativePath)
	if dirPath == "." || dirPath == "" {
		return "root"
	}
	return dirPath
}

func (u *Upserter) generateDocumentID(filePath string, chunkIndex int) string {
	// Create a unique ID based on file path and chunk index
	hash := md5.Sum([]byte(fmt.Sprintf("%s:%d", filePath, chunkIndex)))
	return fmt.Sprintf("doc_%x", hash[:8])
}

func (u *Upserter) UpsertDocument(ctx context.Context, doc document.FileInfo) error {
	logger.LogInfo(fmt.Sprintf("Upserting single document: %s", doc.RelativePath))

	// Extract namespace from file path
	namespace := u.extractNamespace(doc.RelativePath)

	chunker := document.NewChunker(1000)
	chunks, err := chunker.ChunkDocument(doc)
	if err != nil {
		return fmt.Errorf("error chunking document %s: %w", doc.RelativePath, err)
	}

	documents := make([]Document, len(chunks))
	for i, chunk := range chunks {
		docID := u.generateDocumentID(doc.RelativePath, chunk.Index)

		metadata := make(map[string]string)
		for k, v := range chunk.Metadata {
			metadata[k] = v
		}
		metadata["chunk_index"] = fmt.Sprintf("%d", chunk.Index)
		metadata["source_file"] = doc.RelativePath
		metadata["file_size"] = fmt.Sprintf("%d", doc.Size)

		// Add recipe/project information
		fullPath := u.extractFullPath(doc.RelativePath)
		metadata["namespace"] = namespace
		metadata["full_path"] = fullPath
		metadata["recipe_name"] = u.extractRecipeName(doc.RelativePath)
		metadata["project_type"] = u.extractProjectType(doc.RelativePath)

		documents[i] = Document{
			ID:        docID,
			Content:   chunk.Content,
			Metadata:  metadata,
			Namespace: namespace,
		}
	}

	return u.client.UpsertBatch(ctx, documents, namespace)
}

func (u *Upserter) ValidateDocument(doc document.FileInfo) error {
	if strings.TrimSpace(doc.Content) == "" {
		return fmt.Errorf("document content is empty: %s", doc.RelativePath)
	}

	if doc.Size > 10*1024*1024 { // 10MB limit
		return fmt.Errorf("document too large: %s (%d bytes)", doc.RelativePath, doc.Size)
	}

	if doc.Topic == "" {
		return fmt.Errorf("document has no topic: %s", doc.RelativePath)
	}

	return nil
}
