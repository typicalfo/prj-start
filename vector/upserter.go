package vector

import (
	"context"
	"crypto/md5"
	"fmt"
	"prj-start/document"
	"prj-start/logger"
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

	// Second pass: process documents
	batch := make([]Document, 0, u.batchSize)

	for i, doc := range documents {
		logger.LogProgress(i+1, len(documents), fmt.Sprintf("Processing document: %s", doc.RelativePath))

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

			batch = append(batch, Document{
				ID:       docID,
				Content:  chunk.Content,
				Metadata: metadata,
			})

			// Process batch when it reaches the size limit
			if len(batch) >= u.batchSize {
				err := u.client.UpsertBatch(ctx, batch)
				if err != nil {
					logger.LogError(fmt.Sprintf("Error upserting batch: %v", err))
					return err
				}
				processedChunks += len(batch)
				logger.LogProgress(processedChunks, totalChunks, "Chunks processed")
				batch = batch[:0] // Clear batch
			}
		}
	}

	// Process remaining items in batch
	if len(batch) > 0 {
		err := u.client.UpsertBatch(ctx, batch)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error upserting final batch: %v", err))
			return err
		}
		processedChunks += len(batch)
		logger.LogProgress(processedChunks, totalChunks, "Final chunks processed")
	}

	logger.LogSuccess(fmt.Sprintf("Upsert completed! Processed %d chunks from %d documents", processedChunks, len(documents)-failedDocuments))
	if failedDocuments > 0 {
		logger.LogWarning(fmt.Sprintf("Failed to process %d documents", failedDocuments))
	}

	return nil
}

func (u *Upserter) generateDocumentID(filePath string, chunkIndex int) string {
	// Create a unique ID based on file path and chunk index
	hash := md5.Sum([]byte(fmt.Sprintf("%s:%d", filePath, chunkIndex)))
	return fmt.Sprintf("doc_%x", hash[:8])
}

func (u *Upserter) UpsertDocument(ctx context.Context, doc document.FileInfo) error {
	logger.LogInfo(fmt.Sprintf("Upserting single document: %s", doc.RelativePath))

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

		documents[i] = Document{
			ID:       docID,
			Content:  chunk.Content,
			Metadata: metadata,
		}
	}

	return u.client.UpsertBatch(ctx, documents)
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
