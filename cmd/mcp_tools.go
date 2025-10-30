package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/upstash/vector-go"
)

// VectorQueryInput represents input for vector query tool
type VectorQueryInput struct {
	Query           string `json:"query" jsonschema:"natural language query to search for"`
	TopK            int    `json:"topK,omitempty" jsonschema:"maximum number of results to return"`
	Namespace       string `json:"namespace,omitempty" jsonschema:"namespace to search within"`
	IncludeMetadata bool   `json:"includeMetadata,omitempty" jsonschema:"include document metadata in results"`
	IncludeData     bool   `json:"includeData,omitempty" jsonschema:"include document content in results"`
}

// VectorQueryOutput represents output for vector query tool
type VectorQueryOutput struct {
	Results []QueryResult `json:"results"`
	Query   string        `json:"query"`
	Count   int           `json:"count"`
}

// QueryResult represents a single query result
type QueryResult struct {
	ID       string                 `json:"id"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Data     string                 `json:"data,omitempty"`
}

// MetadataQueryInput represents input for metadata query tool
type MetadataQueryInput struct {
	Filter          string `json:"filter,omitempty" jsonschema:"metadata filter expression"`
	TopK            int    `json:"topK,omitempty" jsonschema:"maximum number of results to return"`
	Namespace       string `json:"namespace,omitempty" jsonschema:"namespace to search within"`
	IncludeMetadata bool   `json:"includeMetadata,omitempty" jsonschema:"include document metadata in results"`
	IncludeData     bool   `json:"includeData,omitempty" jsonschema:"include document content in results"`
}

// ListNamespacesInput represents input for list namespaces tool
type ListNamespacesInput struct {
	// No input required
}

// ListNamespacesOutput represents output for list namespaces tool
type ListNamespacesOutput struct {
	Namespaces []string `json:"namespaces"`
	Count      int      `json:"count"`
}

// GetDocumentInput represents input for get document tool
type GetDocumentInput struct {
	ID              string `json:"id" jsonschema:"document ID to retrieve"`
	Namespace       string `json:"namespace,omitempty" jsonschema:"namespace containing the document"`
	IncludeMetadata bool   `json:"includeMetadata,omitempty" jsonschema:"include document metadata in result"`
	IncludeData     bool   `json:"includeData,omitempty" jsonschema:"include document content in result"`
}

// GetDocumentOutput represents output for get document tool
type GetDocumentOutput struct {
	ID       string                 `json:"id"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Data     string                 `json:"data,omitempty"`
}

// addVectorQueryTool adds the vector query tool to the MCP server
func addVectorQueryTool(server *mcp.Server, vectorClient *vector.Index) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "vector_query",
		Description: "Query documents using natural language semantic search",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input VectorQueryInput) (
		*mcp.CallToolResult,
		VectorQueryOutput,
		error,
	) {
		// Set defaults
		if input.TopK <= 0 {
			input.TopK = 5
		}

		// Perform query using raw data (automatic embedding)
		var results []vector.VectorScore
		var err error

		if input.Namespace != "" {
			namespaceClient := vectorClient.Namespace(input.Namespace)
			results, err = namespaceClient.QueryData(vector.QueryData{
				Data:            input.Query,
				TopK:            input.TopK,
				IncludeMetadata: input.IncludeMetadata,
				IncludeData:     input.IncludeData,
			})
		} else {
			results, err = vectorClient.QueryData(vector.QueryData{
				Data:            input.Query,
				TopK:            input.TopK,
				IncludeMetadata: input.IncludeMetadata,
				IncludeData:     input.IncludeData,
			})
		}

		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Query failed: %v", err)}},
				IsError: true,
			}, VectorQueryOutput{}, err
		}

		// Convert results
		queryResults := make([]QueryResult, len(results))
		for i, result := range results {
			queryResults[i] = QueryResult{
				ID:       result.Id,
				Score:    float64(result.Score),
				Metadata: result.Metadata,
				Data:     result.Data,
			}
		}

		output := VectorQueryOutput{
			Results: queryResults,
			Query:   input.Query,
			Count:   len(queryResults),
		}

		if debug {
			log.Printf("Vector query: %s returned %d results", input.Query, len(queryResults))
		}

		return nil, output, nil
	})
}

// addMetadataQueryTool adds the metadata filter query tool to the MCP server
func addMetadataQueryTool(server *mcp.Server, vectorClient *vector.Index) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "metadata_query",
		Description: "Query documents using metadata filters",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input MetadataQueryInput) (
		*mcp.CallToolResult,
		VectorQueryOutput,
		error,
	) {
		// Set defaults
		if input.TopK <= 0 {
			input.TopK = 10
		}

		// Perform query with metadata filter
		var results []vector.VectorScore
		var err error

		if input.Namespace != "" {
			namespaceClient := vectorClient.Namespace(input.Namespace)
			if input.Filter != "" {
				results, err = namespaceClient.Query(vector.Query{
					Vector:          []float32{0.0, 0.0}, // Dummy vector
					TopK:            input.TopK,
					IncludeMetadata: input.IncludeMetadata,
					IncludeData:     input.IncludeData,
					Filter:          input.Filter,
				})
			} else {
				results, err = namespaceClient.Query(vector.Query{
					Vector:          []float32{0.0, 0.0}, // Dummy vector
					TopK:            input.TopK,
					IncludeMetadata: input.IncludeMetadata,
					IncludeData:     input.IncludeData,
				})
			}
		} else {
			if input.Filter != "" {
				results, err = vectorClient.Query(vector.Query{
					Vector:          []float32{0.0, 0.0}, // Dummy vector
					TopK:            input.TopK,
					IncludeMetadata: input.IncludeMetadata,
					IncludeData:     input.IncludeData,
					Filter:          input.Filter,
				})
			} else {
				results, err = vectorClient.Query(vector.Query{
					Vector:          []float32{0.0, 0.0}, // Dummy vector
					TopK:            input.TopK,
					IncludeMetadata: input.IncludeMetadata,
					IncludeData:     input.IncludeData,
				})
			}
		}

		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Metadata query failed: %v", err)}},
				IsError: true,
			}, VectorQueryOutput{}, err
		}

		// Convert results
		queryResults := make([]QueryResult, len(results))
		for i, result := range results {
			queryResults[i] = QueryResult{
				ID:       result.Id,
				Score:    float64(result.Score),
				Metadata: result.Metadata,
				Data:     result.Data,
			}
		}

		output := VectorQueryOutput{
			Results: queryResults,
			Query:   fmt.Sprintf("metadata filter: %s", input.Filter),
			Count:   len(queryResults),
		}

		if debug {
			log.Printf("Metadata query: %s returned %d results", input.Filter, len(queryResults))
		}

		return nil, output, nil
	})
}

// addListNamespacesTool adds the list namespaces tool to the MCP server
func addListNamespacesTool(server *mcp.Server, vectorClient *vector.Index) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_namespaces",
		Description: "List all available namespaces in the vector database",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ListNamespacesInput) (
		*mcp.CallToolResult,
		ListNamespacesOutput,
		error,
	) {
		namespaces, err := vectorClient.ListNamespaces()
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to list namespaces: %v", err)}},
				IsError: true,
			}, ListNamespacesOutput{}, err
		}

		output := ListNamespacesOutput{
			Namespaces: namespaces,
			Count:      len(namespaces),
		}

		if debug {
			log.Printf("Listed %d namespaces", len(namespaces))
		}

		return nil, output, nil
	})
}

// addGetDocumentTool adds the get document tool to the MCP server
func addGetDocumentTool(server *mcp.Server, vectorClient *vector.Index) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_document",
		Description: "Retrieve a specific document by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input GetDocumentInput) (
		*mcp.CallToolResult,
		GetDocumentOutput,
		error,
	) {
		// Fetch the document
		var vectors []vector.Vector
		var err error

		if input.Namespace != "" {
			namespaceClient := vectorClient.Namespace(input.Namespace)
			vectors, err = namespaceClient.Fetch(vector.Fetch{
				Ids:             []string{input.ID},
				IncludeMetadata: input.IncludeMetadata,
				IncludeData:     input.IncludeData,
			})
		} else {
			vectors, err = vectorClient.Fetch(vector.Fetch{
				Ids:             []string{input.ID},
				IncludeMetadata: input.IncludeMetadata,
				IncludeData:     input.IncludeData,
			})
		}

		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to fetch document: %v", err)}},
				IsError: true,
			}, GetDocumentOutput{}, err
		}

		if len(vectors) == 0 {
			return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Document with ID '%s' not found", input.ID)}},
				IsError: true,
			}, GetDocumentOutput{}, fmt.Errorf("document not found")
		}

		vector := vectors[0]
		output := GetDocumentOutput{
			ID:       vector.Id,
			Metadata: vector.Metadata,
			Data:     vector.Data,
		}

		if debug {
			log.Printf("Retrieved document: %s", input.ID)
		}

		return nil, output, nil
	})
}
