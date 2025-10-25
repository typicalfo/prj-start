package vector

import (
	"context"
	"fmt"
	"github.com/typicalfo/prj-start/config"
	"github.com/typicalfo/prj-start/logger"

	"github.com/upstash/vector-go"
)

type Client struct {
	config *config.UpstashConfig
	index  *vector.Index
}

func NewClient(cfg *config.UpstashConfig) (*Client, error) {
	if !cfg.Validate() {
		return nil, fmt.Errorf("invalid Upstash configuration")
	}

	logger.LogInfo(fmt.Sprintf("Raw URL from config: '%s'", cfg.URL))
	logger.LogInfo(fmt.Sprintf("URL length: %d", len(cfg.URL)))
	logger.LogInfo(fmt.Sprintf("URL bytes: %v", []byte(cfg.URL)))

	// Initialize actual Upstash Vector client
	index := vector.NewIndex(cfg.URL, cfg.Token)

	logger.LogSuccess(fmt.Sprintf("Upstash Vector client initialized with URL: %s", cfg.URL))
	return &Client{
		config: cfg,
		index:  index,
	}, nil
}

func (c *Client) Upsert(ctx context.Context, id string, metadata map[string]string, content string, namespace string) error {
	logger.LogInfo(fmt.Sprintf("Upserting document: %s (namespace: %s)", id, namespace))

	// Convert metadata to map[string]any for Upstash SDK
	upsertMetadata := make(map[string]any)
	for k, v := range metadata {
		upsertMetadata[k] = v
	}

	// Create namespace instance
	ns := c.index.Namespace(namespace)

	// Use UpsertData to let Upstash handle embedding
	err := ns.UpsertData(vector.UpsertData{
		Id:       id,
		Data:     content,
		Metadata: upsertMetadata,
	})

	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to upsert document %s: %v", id, err))
		return err
	}

	logger.LogSuccess(fmt.Sprintf("Successfully upserted document: %s (namespace: %s)", id, namespace))
	return nil
}

func (c *Client) UpsertBatch(ctx context.Context, documents []Document, namespace string) error {
	logger.LogInfo(fmt.Sprintf("Upserting batch of %d documents (namespace: %s)", len(documents), namespace))
	logger.LogInfo(fmt.Sprintf("Namespace debug: '%s' (len=%d)", namespace, len(namespace)))

	// Convert documents to UpsertData format
	upsertData := make([]vector.UpsertData, len(documents))
	for i, doc := range documents {
		// Convert metadata to map[string]any for Upstash SDK
		upsertMetadata := make(map[string]any)
		for k, v := range doc.Metadata {
			upsertMetadata[k] = v
		}

		upsertData[i] = vector.UpsertData{
			Id:       doc.ID,
			Data:     doc.Content,
			Metadata: upsertMetadata,
		}
	}

	// Create namespace instance
	logger.LogInfo(fmt.Sprintf("Creating namespace client for: '%s'", namespace))
	ns := c.index.Namespace(namespace)

	// Use UpsertDataMany for batch processing
	logger.LogInfo(fmt.Sprintf("Calling UpsertDataMany with %d documents", len(upsertData)))
	err := ns.UpsertDataMany(upsertData)
	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to upsert batch of %d documents: %v", len(documents), err))
		return err
	}

	logger.LogSuccess(fmt.Sprintf("Successfully upserted batch of %d documents (namespace: %s)", len(documents), namespace))
	return nil
}

type Document struct {
	ID        string
	Content   string
	Metadata  map[string]string
	Namespace string
}

type QueryResult struct {
	ID       string                 `json:"id"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Data     string                 `json:"data,omitempty"`
}

func (c *Client) Query(ctx context.Context, queryVector []float32, topK int, namespace string) ([]QueryResult, error) {
	logger.LogInfo(fmt.Sprintf("Querying namespace '%s' with vector of length %d", namespace, len(queryVector)))

	// Create namespace instance
	ns := c.index.Namespace(namespace)

	// Query with vector
	results, err := ns.Query(vector.Query{
		Vector:          queryVector,
		TopK:            topK,
		IncludeMetadata: true,
	})

	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to query namespace %s: %v", namespace, err))
		return nil, err
	}

	// Convert results to our format
	queryResults := make([]QueryResult, len(results))
	for i, result := range results {
		// Convert metadata to map[string]interface{} for flexibility
		metadata := make(map[string]interface{})
		for k, v := range result.Metadata {
			metadata[k] = v
		}

		queryResults[i] = QueryResult{
			ID:       result.Id,
			Score:    float64(result.Score),
			Metadata: metadata,
		}
	}

	logger.LogSuccess(fmt.Sprintf("Query returned %d results from namespace '%s'", len(queryResults), namespace))
	return queryResults, nil
}

func (c *Client) ListNamespaces(ctx context.Context) ([]string, error) {
	logger.LogInfo("Listing all namespaces")

	namespaces, err := c.index.ListNamespaces()
	if err != nil {
		logger.LogError(fmt.Sprintf("Failed to list namespaces: %v", err))
		return nil, err
	}

	logger.LogSuccess(fmt.Sprintf("Found %d namespaces", len(namespaces)))
	return namespaces, nil
}
