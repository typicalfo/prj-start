# Vector Document Ingestion System

A Go-based tool for intelligently chunking and ingesting development documentation into Upstash Vector database for enhanced search and retrieval.

## Overview

This project processes development documents from `dev-docs/` folder, intelligently chunks them based on content type, and upserts them to an Upstash Vector database. The system automatically generates embeddings and handles metadata extraction for improved search capabilities.

The project uses a dual approach:
- **Ingestion**: Go-based code for chunking and upserting documents
- **Querying**: Upstash MCP server for natural language queries against indexed data

## Quickstart

### Installation

**Option 1: Go Install (Recommended)**
```bash
go install github.com/typicalfo/prj-start@latest
```

**Option 2: Build from Source**
```bash
git clone https://github.com/typicalfo/prj-start.git
cd prj-start
make build
```

### Setup

1. **Initialize configuration** (interactive setup):
   ```bash
   prj-start init
   ```
   
   This will guide you through setting up Upstash Vector credentials and save them to:
   - `~/.config/prj-start/config.yaml` (Linux/macOS)
   - `%LOCALAPPDATA%/prj-start/config.yaml` (Windows)

2. **Organize your documents** in any folder structure:
   ```
   my-docs/
   ├── go-fiber-recipes/
   │   ├── 404-handler/
   │   ├── authentication/
   │   └── database/
   ├── clean-architecture/
   └── your-project/
   ```

3. **Process your documents**:
   ```bash
   # Process current directory
   prj-start
   
   # Process specific folder
   prj-start --folder /path/to/docs
   
   # Verbose output
   prj-start --folder ./docs --verbose
   ```

4. **Query your data** using natural language with [Upstash MCP Server](#upstash-mcp-server-for-querying)

For detailed setup instructions, see:
- [Configuration](#configuration) - Complete environment setup
- [Upstash MCP Server](#upstash-mcp-server-for-querying) - Natural language querying
- [Agent Instructions](./AGENT_MCP_INSTRUCTIONS.md) - Development agent guidance

## Features

- **Intelligent Document Chunking**: Content-aware chunking strategies for different file types
  - Go files: Chunked by functions, structs, and interfaces
  - Markdown: Chunked by headers and sections
  - SQL: Chunked by statements
  - Config files: Chunked by logical sections
  - Text files: Paragraph-based chunking

- **Metadata Extraction**: Automatically extracts and includes:
  - Topic (folder name) for categorization
  - File path and extension
  - Chunk position within document
  - Content type identification

- **Robust Processing**: 
  - Batch processing for efficiency
  - Error handling and recovery
  - Progress tracking with colored logging
  - Skip binary files and large files automatically

- **Developer-Friendly**:
  - Colored logging for clear visibility
  - Comprehensive Makefile for common operations
  - Environment-based configuration
  - Mock mode for testing

## Project Structure

```
├── vector.go                 # Main application entry point
├── config/
│   └── upstash_config.go     # Configuration management
├── logger/
│   └── colored_logger.go     # Colored logging utilities
├── document/
│   ├── reader.go            # Document reading functionality
│   └── chunker.go           # Intelligent content chunking
├── vector/
│   ├── client.go            # Upstash Vector client
│   └── upserter.go          # Batch upsert operations
├── dev-docs/                # Source documents to process
├── Makefile                 # Build and development commands
└── README.md               # This file
```

## Installation

### Prerequisites

- Go 1.24+ installed
- Upstash Vector database credentials

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd prj-start
```

2. Install dependencies:
```bash
make deps
```

3. Set up environment variables:
   
   The application uses `godotenv` to automatically load environment variables from a `.env` file. Copy the example file and update it with your credentials:
   
   ```bash
   cp .env.example .env
   # Edit .env with your Upstash Vector credentials
   ```
   
   Or set them manually:
   ```bash
   export UPSTASH_VECTOR_REST_URL="your-upstash-vector-url"
   export UPSTASH_VECTOR_REST_TOKEN="your-upstash-vector-token"
   ```

## Usage

### Build and Run

```bash
# Build the application
make build

# Run with default settings
./main

# Or run directly with Go
make run
```

### Development

```bash
# Run with hot reload (requires Air)
make dev

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean
```

### Available Commands

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Build and run
make dev           # Run with Air hot reload
make test          # Run tests
make clean         # Clean build artifacts
make deps          # Install dependencies
make fmt           # Format code
make kill          # Kill running processes
```

## Configuration

The application uses `godotenv` to automatically load environment variables from a `.env` file.

### Required Variables

- `UPSTASH_VECTOR_REST_URL`: Your Upstash Vector database URL
- `UPSTASH_VECTOR_REST_TOKEN`: Your Upstash Vector authentication token

### Optional Variables

- `UPSTASH_VECTOR_INDEX_URL`: Your Upstash Vector index URL
- `UPSTASH_EMAIL`: Email for Upstash MCP server (for querying)
- `UPSTASH_API_KEY`: API key for Upstash MCP server (for querying)
- `BATCH_SIZE`: Number of documents to process in each batch (default: 10)
- `PROCESSING_TIMEOUT_MINUTES`: Timeout for document processing (default: 30)
- `LOG_LEVEL`: Logging level - debug, info, warn, error (default: info)

### Environment File

Create a `.env` file in project root:

```bash
# Upstash Vector Configuration (Required)
# Get these from your Upstash Console (https://console.upstash.com)
UPSTASH_VECTOR_REST_URL=https://your-vector-url.upstash.io
UPSTASH_VECTOR_REST_TOKEN=your-verification-token
UPSTASH_VECTOR_INDEX_URL=https://your-index-url.upstash.io

# Upstash MCP Server Configuration (Optional)
# Get these from Account > Management API > Create API key in Upstash Console
UPSTASH_EMAIL=your-email@example.com
UPSTASH_API_KEY=your-upstash-api-key

# Application Settings (Optional)
# Batch size for document upsert operations (default: 10)
BATCH_SIZE=10

# Timeout for document processing in minutes (default: 30)
PROCESSING_TIMEOUT_MINUTES=30

# Log level: debug, info, warn, error (default: info)
LOG_LEVEL=info
```

The application will automatically load these variables when started.

### Upstash MCP Server for Querying

For querying the indexed data using natural language, configure the Upstash MCP server:

1. **Get API Key**: Go to `Account > Management API > Create API key` in Upstash Console
2. **Configure MCP**: Add to your MCP configuration file:

```json
{
  "mcpServers": {
    "upstash": {
      "command": "npx",
      "args": [
        "-y",
        "@upstash/mcp-server",
        "run",
        "<UPSTASH_EMAIL>",
        "<UPSTASH_API_KEY>"
      ]
    }
  }
}
```

For detailed MCP setup instructions, see [dev-docs/upstash/MCP.md](./dev-docs/upstash/MCP.md).

### Example .env file

```bash
# Upstash Vector Configuration
UPSTASH_VECTOR_REST_URL=https://your-vector-url.upstash.io
UPSTASH_VECTOR_REST_TOKEN=your-verification-token
UPSTASH_VECTOR_INDEX_URL=https://your-index-url.upstash.io
```

## Document Processing

### Supported File Types

- **Go files** (`.go`): Chunked by functions, structs, interfaces, variables, and constants
- **Markdown** (`.md`): Chunked by headers and sections
- **SQL** (`.sql`): Chunked by individual statements
- **Configuration** (`.json`, `.yaml`, `.yml`, `.toml`): Chunked by logical sections
- **HTML** (`.html`): Chunked by major structural elements
- **Text files**: Paragraph-based chunking with size limits

### Metadata Schema

Each chunk includes the following metadata:

```json
{
  "filename": "go-fiber-recipes/clean-architecture/main.go",
  "topic": "clean-architecture",
  "extension": ".go",
  "chunk_index": "2",
  "total_chunks": "5",
  "chunk_type": "go_construct",
  "source_file": "go-fiber-recipes/clean-architecture/main.go",
  "file_size": "1024"
}
```

### Processing Flow

1. **Document Discovery**: Recursively scan `dev-docs/` folder
2. **File Filtering**: Skip binary files and files > 1MB
3. **Content Reading**: Read text content with proper encoding
4. **Intelligent Chunking**: Apply content-specific chunking strategies
5. **Metadata Extraction**: Extract file and chunk metadata
6. **Batch Upsert**: Send chunks to Upstash Vector in batches
7. **Progress Tracking**: Real-time progress with colored logging

## Example Output

### Ingestion Process
```
[INFO] Starting vector upsert process
[SUCCESS] Upstash Vector client initialized
[INFO] Reading documents from dev-docs folder...
[SUCCESS] Successfully read 106 documents
[INFO] Starting upsert of 106 documents
[PROGRESS] 1/106 (0.9%) - Processing document: go-fiber-recipes/404-handler/README.md
[SUCCESS] Created 11 chunks for go-fiber-recipes/404-handler/README.md
[INFO] Upserting batch of 10 documents
[SUCCESS] Successfully upserted batch of 10 documents
[SUCCESS] Vector upsert completed successfully in 2.345s
```

### Query Examples with MCP
Once data is ingested, you can query using natural language:
- "Show me all Go files related to clean architecture"
- "Find examples of Fiber middleware usage"
- "What database connection patterns are used in the recipes?"
- "Show me error handling patterns in Go Fiber"

## Architecture and Workflow

### Data Flow
1. **Document Discovery**: Scan `dev-docs/` folder for supported file types
2. **Intelligent Chunking**: Apply content-specific chunking strategies
3. **Metadata Extraction**: Extract file paths, topics, and content types
4. **Vector Ingestion**: Upsert chunks with embeddings to Upstash Vector
5. **Natural Language Querying**: Use Upstash MCP server for queries

### Separation of Concerns
- **Ingestion Code** (this project): Handles document processing and upserting
- **Query Interface** (Upstash MCP): Provides natural language access to indexed data

### Benefits of This Approach
- **Specialized Tools**: Each component uses the best tool for its job
- **Scalable Querying**: MCP server handles complex natural language queries
- **Maintainable Code**: Clear separation between ingestion and querying
- **Flexible Access**: Query from any MCP-compatible client (Cursor, Claude, etc.)

## Development Guidelines

See [AGENTS.md](./AGENTS.md) for detailed development guidelines and coding standards.

## Implementation Details

The implementation follows clean architecture principles:

- **Configuration Layer**: Environment-based configuration management
- **Document Layer**: File reading and intelligent chunking
- **Vector Layer**: Upstash client and batch operations
- **Logging Layer**: Colored logging for visibility

For detailed implementation plans and status, see:
- [IMPLEMENTATION_PLAN.md](./IMPLEMENTATION_PLAN.md)
- [CURRENT_STATUS.md](./CURRENT_STATUS.md)

## Additional Documentation

- [Upstash Go SDK Documentation](./UPSTASH-GO-SDK.md) - Complete API reference for the Go client
- [Upstash Upsert API](./UPSTASH-UPSERT.md) - API documentation for upsert operations
- [Upstash MCP Server](./dev-docs/upstash/MCP.md) - Natural language querying setup guide
- [Agent MCP Instructions](./AGENT_MCP_INSTRUCTIONS.md) - Comprehensive guide for development agents using MCP
- [Quick Start MCP](./QUICK_START_MCP.md) - Fast-start guide for agents to begin using the vector database

## License

This project is part of a larger Go backend development initiative.

## Contributing

Follow the development guidelines in AGENTS.md and use the provided Makefile targets for consistent development workflows.