package cmd

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// addHelpResource adds a help resource to the MCP server
func addHelpResource(server *mcp.Server) {
	server.AddResource(&mcp.Resource{
		URI:         "vector://help",
		Name:        "Vector Database Help",
		Description: "Help and usage information for vector database tools",
		MIMEType:    "text/plain",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		helpText := `# Vector Database MCP Server Help

This MCP server provides tools for querying documents stored in Upstash Vector database.

## Available Tools

### vector_query
Query documents using natural language semantic search.
- query (required): Natural language query to search for
- topK (optional): Maximum number of results (default: 5)
- namespace (optional): Namespace to search within
- includeMetadata (optional): Include document metadata in results
- includeData (optional): Include document content in results

### metadata_query
Query documents using metadata filters.
- filter (optional): Metadata filter expression
- topK (optional): Maximum number of results (default: 10)
- namespace (optional): Namespace to search within
- includeMetadata (optional): Include document metadata in results
- includeData (optional): Include document content in results

### list_namespaces
List all available namespaces in the vector database.
- No input parameters required

### get_document
Retrieve a specific document by ID.
- id (required): Document ID to retrieve
- namespace (optional): Namespace containing the document
- includeMetadata (optional): Include document metadata in result
- includeData (optional): Include document content in result

## Usage Examples

1. Search for Go Fiber examples:
   vector_query(query="Go Fiber routing examples", topK=3)

2. Filter by project type:
   metadata_query(filter="project_type = 'go-fiber-recipes'")

3. List all namespaces:
   list_namespaces()

4. Get specific document:
   get_document(id="doc-123", includeData=true)

## Metadata Fields

Documents contain the following metadata fields:
- namespace: Document namespace
- full_path: Complete directory path
- recipe_name: Recipe folder name
- project_type: Top-level project category
- source_file: Original file path
- filename: Original filename
- extension: File extension
- chunk_index: Chunk number in document
- total_chunks: Total chunks in document

## Tips

- Use vector_query for semantic search with natural language
- Use metadata_query for precise filtering by known fields
- Combine tools for complex queries
- Use includeData=true to get document content
- Use includeMetadata=true to get document metadata`

		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:      req.Params.URI,
				Text:     helpText,
				MIMEType: "text/plain",
			}},
		}, nil
	})
}
