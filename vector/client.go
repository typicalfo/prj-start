package vector

import (
	"context"
	"fmt"
	"prj-start/config"
	"prj-start/logger"
)

// Mock client for now - will implement actual Upstash Vector client
type Client struct {
	config *config.UpstashConfig
}

func NewClient(cfg *config.UpstashConfig) (*Client, error) {
	if !cfg.Validate() {
		return nil, fmt.Errorf("invalid Upstash configuration")
	}

	logger.LogSuccess("Upstash Vector client initialized (mock)")
	return &Client{
		config: cfg,
	}, nil
}

func (c *Client) Upsert(ctx context.Context, id string, metadata map[string]string, content string) error {
	logger.LogInfo(fmt.Sprintf("Upserting document: %s", id))

	// TODO: Implement actual Upstash Vector upsert
	// For now, just log the operation
	logger.LogInfo(fmt.Sprintf("Content preview: %s...", content[:min(50, len(content))]))
	logger.LogSuccess(fmt.Sprintf("Successfully upserted document: %s", id))
	return nil
}

func (c *Client) UpsertBatch(ctx context.Context, documents []Document) error {
	logger.LogInfo(fmt.Sprintf("Upserting batch of %d documents", len(documents)))

	for i, doc := range documents {
		err := c.Upsert(ctx, doc.ID, doc.Metadata, doc.Content)
		if err != nil {
			logger.LogError(fmt.Sprintf("Failed to upsert document %d in batch: %v", i, err))
			return err
		}
	}

	logger.LogSuccess(fmt.Sprintf("Successfully upserted batch of %d documents", len(documents)))
	return nil
}

type Document struct {
	ID       string
	Content  string
	Metadata map[string]string
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
