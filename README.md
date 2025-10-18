# Vector Document Ingestion System

A Go-based tool for intelligently chunking and ingesting development documentation into Upstash Vector database for enhanced search and retrieval.

## Overview

This project processes development documents from the `dev-docs/` folder, intelligently chunks them based on content type, and upserts them to an Upstash Vector database. The system automatically generates embeddings and handles metadata extraction for improved search capabilities.

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
```bash
export UPSTASH_VECTOR_URL="your-upstash-vector-url"
export UPSTASH_VECTOR_TOKEN="your-upstash-vector-token"
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

The application uses environment variables for configuration:

- `UPSTASH_VECTOR_URL`: Your Upstash Vector database URL
- `UPSTASH_VECTOR_TOKEN`: Your Upstash Vector authentication token

### Example .env file

```bash
UPSTASH_VECTOR_URL=https://your-vector-url.upstash.io
UPSTASH_VECTOR_TOKEN=your-verification-token
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

## License

This project is part of a larger Go backend development initiative.

## Contributing

Follow the development guidelines in AGENTS.md and use the provided Makefile targets for consistent development workflows.