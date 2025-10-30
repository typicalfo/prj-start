# MCP Server Implementation Plan

## Goal
Create an MCP server that allows local agents with stdio to query the configured Upstash Vector database containing ingested documents.

## Requirements Analysis

### MCP Server Requirements
- **Transport**: Use stdio transport for local agent communication
- **Protocol**: Implement Model Context Protocol (MCP) v0.5.0
- **Tools**: Provide query tools for vector database access
- **Resources**: Optional static resources for help/documentation
- **Error Handling**: Proper error responses and logging

### Upstash Vector Integration Requirements
- **Client**: Use existing Upstash Vector Go SDK
- **Configuration**: Load from existing config system
- **Namespaces**: Support namespace-based querying
- **Query Types**: Support both vector and metadata queries
- **Embeddings**: Support raw data queries with automatic embedding

## Implementation Plan

### Phase 1: Dependencies and Setup
1. **Add MCP SDK Dependency**
   ```bash
   go get github.com/modelcontextprotocol/go-sdk v0.5.0
   ```

2. **Update go.mod** to include both MCP SDK and Upstash Vector SDK

3. **Create MCP Server Structure**
   - New file: `cmd/mcp_server.go` (or update existing `cmd/mcp.go`)
   - Separate MCP logic from command handling

### Phase 2: Core MCP Server Implementation

#### 2.1 Server Initialization
```go
func createMCPServer(cfg *config.Config) *mcp.Server {
    server := mcp.NewServer(&mcp.Implementation{
        Name:    "prj-start-vector-db",
        Version: "1.0.0",
    }, &mcp.ServerOptions{
        InitializedHandler: func(ctx context.Context, req *mcp.InitializedRequest) {
            log.Println("MCP client initialized")
        },
    })
    
    return server
}
```

#### 2.2 Upstash Vector Client Integration
```go
func createVectorClient(cfg *config.Config) *vector.Index {
    if cfg.Upstash.URL != "" && cfg.Upstash.Token != "" {
        return vector.NewIndex(cfg.Upstash.URL, cfg.Upstash.Token)
    }
    return vector.NewIndexFromEnv()
}
```

### Phase 3: Tool Implementation

#### 3.1 Vector Query Tool
**Purpose**: Query vectors using semantic similarity

**Input Schema**:
```go
type VectorQueryInput struct {
    Query      string `json:"query" jsonschema:"natural language query to search for"`
    TopK       int    `json:"topK" jsonschema:"maximum number of results to return (default: 5)"`
    Namespace  string `json:"namespace" jsonschema:"namespace to search within (optional)"`
    IncludeMetadata bool `json:"includeMetadata" jsonschema:"include document metadata in results"`
    IncludeData     bool `json:"includeData" jsonschema:"include document content in results"`
}
```

**Output Schema**:
```go
type VectorQueryOutput struct {
    Results []QueryResult `json:"results"`
    Query   string       `json:"query"`
    Count   int          `json:"count"`
}

type QueryResult struct {
    ID       string                 `json:"id"`
    Score    float64                `json:"score"`
    Metadata map[string]interface{}  `json:"metadata,omitempty"`
    Data     string                 `json:"data,omitempty"`
}
```

#### 3.2 Metadata Filter Query Tool
**Purpose**: Query documents using metadata filters

**Input Schema**:
```go
type MetadataQueryInput struct {
    Filter          string `json:"filter" jsonschema:"metadata filter expression"`
    TopK            int    `json:"topK" jsonschema:"maximum number of results to return (default: 10)"`
    Namespace        string `json:"namespace" jsonschema:"namespace to search within (optional)"`
    IncludeMetadata bool   `json:"includeMetadata" jsonschema:"include document metadata in results"`
    IncludeData     bool   `json:"includeData" jsonschema:"include document content in results"`
}
```

#### 3.3 List Namespaces Tool
**Purpose**: List all available namespaces

**Input Schema**:
```go
type ListNamespacesInput struct {
    // No input required
}
```

**Output Schema**:
```go
type ListNamespacesOutput struct {
    Namespaces []string `json:"namespaces"`
    Count      int      `json:"count"`
}
```

#### 3.4 Get Document Tool
**Purpose**: Retrieve specific document by ID

**Input Schema**:
```go
type GetDocumentInput struct {
    ID              string `json:"id" jsonschema:"document ID to retrieve"`
    Namespace       string `json:"namespace" jsonschema:"namespace containing the document (optional)"`
    IncludeMetadata bool   `json:"includeMetadata" jsonschema:"include document metadata in result"`
    IncludeData     bool   `json:"includeData" jsonschema:"include document content in result"`
}
```

### Phase 4: Resource Implementation

#### 4.1 Help Resource
Static resource providing usage information:
```go
server.AddResource(&mcp.Resource{
    URI:         "vector://help",
    Name:        "Vector Database Help",
    Description: "Help and usage information for vector database tools",
    MIMEType:    "text/plain",
}, helpResourceHandler)
```

#### 4.2 Schema Resource
Resource providing JSON schemas for tools:
```go
server.AddResource(&mcp.Resource{
    URI:         "vector://schemas",
    Name:        "Tool Schemas",
    Description: "JSON schemas for vector database tools",
    MIMEType:    "application/json",
}, schemaResourceHandler)
```

### Phase 5: Command Integration

#### 5.1 Update cmd/mcp.go
Replace placeholder with actual MCP server implementation:

```go
func runMCP(cmd *cobra.Command, args []string) error {
    // Load configuration
    cfg, err := config.LoadConfig(cfgFile)
    if err != nil {
        return fmt.Errorf("failed to load configuration: %w", err)
    }

    // Check Upstash configuration
    if !cfg.HasUpstashConfig() {
        return fmt.Errorf("Upstash configuration incomplete. Run 'prj-start init' first")
    }

    // Create and configure MCP server
    server := createMCPServer(cfg)
    vectorClient := createVectorClient(cfg)
    
    // Add tools
    addVectorQueryTool(server, vectorClient)
    addMetadataQueryTool(server, vectorClient)
    addListNamespacesTool(server, vectorClient)
    addGetDocumentTool(server, vectorClient)
    
    // Add resources
    addHelpResource(server)
    addSchemaResource(server)

    // Start server with stdio transport
    log.Println("Starting MCP server...")
    if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
        return fmt.Errorf("MCP server failed: %w", err)
    }

    return nil
}
```

### Phase 6: Error Handling and Logging

#### 6.1 Error Handling Strategy
- **Configuration Errors**: Clear messages about missing config
- **Vector DB Errors**: Proper error propagation to MCP client
- **Query Errors**: Meaningful error messages for invalid queries
- **Network Errors**: Timeout and retry logic where appropriate

#### 6.2 Logging Strategy
- **Server Level**: MCP server lifecycle events
- **Tool Level**: Query execution and results
- **Debug Level**: Detailed information when verbose flag is set

### Phase 7: Testing

#### 7.1 Unit Tests
- Test each tool handler with mock vector client
- Test error scenarios
- Test configuration loading

#### 7.2 Integration Tests
- Test with real Upstash Vector instance
- Test MCP protocol communication
- Test with sample ingested data

#### 7.3 End-to-End Tests
- Test with actual MCP client
- Test query workflows
- Test error handling

## File Structure

### New Files
```
cmd/
├── mcp.go              # Updated with actual implementation
├── mcp_server.go       # New: MCP server logic
├── mcp_tools.go        # New: Tool implementations
└── mcp_resources.go     # New: Resource handlers

mcp/                    # New package for MCP-specific code
├── server.go           # Server creation and configuration
├── tools.go            # Tool definitions and handlers
├── resources.go        # Resource handlers
└── types.go           # MCP-specific types
```

### Modified Files
```
go.mod                  # Add MCP SDK dependency
cmd/mcp.go             # Replace placeholder with implementation
```

## Configuration Requirements

### Existing Config Support
The MCP server will use the existing configuration system:
- `UPSTASH_VECTOR_REST_URL`
- `UPSTASH_VECTOR_REST_TOKEN`
- `UPSTASH_VECTOR_INDEX_URL` (optional)

### New Config Options (Optional)
- `MCP_LOG_LEVEL` - MCP-specific logging level
- `MCP_DEFAULT_TOP_K` - Default result count
- `MCP_DEFAULT_NAMESPACE` - Default namespace for queries

## Usage Examples

### Agent Integration
```bash
# Start MCP server
prj-start mcp

# Agent can now call tools:
# - vector_query: "Find examples of Go Fiber routing"
# - metadata_query: "filter by project_type='go-fiber-recipes'"
# - list_namespaces: "Show all available namespaces"
# - get_document: "Retrieve document with ID 'doc-123'"
```

### Configuration
```bash
# Set environment variables
export UPSTASH_VECTOR_REST_URL="https://your-vector.upstash.io"
export UPSTASH_VECTOR_REST_TOKEN="your-token"

# Or use config file
prj-start init
prj-start mcp
```

## Success Criteria

1. **Functional MCP Server**: Server starts and communicates via stdio
2. **Tool Availability**: All 4 tools work correctly
3. **Query Functionality**: Can query ingested documents
4. **Error Handling**: Graceful error handling and clear messages
5. **Configuration**: Uses existing config system
6. **Documentation**: Clear help and usage information
7. **Testing**: Comprehensive test coverage

## Next Steps After Implementation

1. **Performance Optimization**: Caching, connection pooling
2. **Advanced Features**: Hybrid queries, resumable queries
3. **Monitoring**: Metrics and health checks
4. **Security**: Access control, rate limiting
5. **Additional Tools**: Document updates, deletions

This plan provides a comprehensive roadmap for implementing a production-ready MCP server that integrates seamlessly with the existing document ingestion system.