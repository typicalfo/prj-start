# Vector Upsert Implementation Plan

## Overview
Implement vector upsert functionality for ingesting all documents in the dev-docs folder into Upstash Vector database. The system will handle document chunking, metadata extraction, and batch upserts with proper error handling and logging.

## Requirements from vector.go
- Ingest all documents in dev-docs folder
- Handle subfolders grouped by topic (include folder name in metadata)
- Create intelligent chunks based on content (not fixed line counts)
- Upstash Vector DB will generate its own embeddings
- Only upload raw data, not pre-generated vectors

## Implementation Phases

### Phase 1: Setup and Configuration
1. **Analyze dev-docs structure**
   - Scan dev-docs folder for subfolders and file types
   - Identify document formats (.go, .md, .sql, .json, etc.)
   - Map folder structure to topics/categories

2. **Upstash Vector client setup**
   - Initialize Upstash Vector client
   - Configure connection parameters
   - Set up authentication

### Phase 2: Document Processing
3. **Document reading functionality**
   - Recursive file walker for dev-docs folder
   - File type detection and appropriate readers
   - Handle different encodings and binary files

4. **Intelligent chunking**
   - Go code: chunk by functions/structs/interfaces
   - Markdown: chunk by headers/sections
   - SQL: chunk by statements/tables
   - Config files: chunk by logical sections
   - Default: paragraph-based chunking

### Phase 3: Metadata and Upload
5. **Metadata extraction**
   - Folder name as topic/category
   - File name and extension
   - File path relative to dev-docs
   - Chunk position within document
   - Content type detection

6. **Vector upsert functionality**
   - Batch processing for efficiency
   - Error handling and retry logic
   - Progress tracking
   - Duplicate detection/handling

### Phase 4: Enhancement and Testing
7. **Colored logging**
   - Integration with logrus or similar
   - Different colors for different operations
   - Progress indicators
   - Error highlighting

8. **Testing and validation**
   - Test with sample documents
   - Validate chunking quality
   - Test metadata accuracy
   - Performance testing

## Technical Architecture

### Core Components
```
main.go
├── config/
│   └── upstash_config.go
├── document/
│   ├── reader.go
│   ├── chunker.go
│   └── metadata.go
├── vector/
│   ├── client.go
│   └── upserter.go
└── logger/
    └── colored_logger.go
```

### Data Flow
1. Scan dev-docs folder → Identify files
2. Read each file → Extract content
3. Analyze content type → Apply appropriate chunking
4. Extract metadata → Combine with chunks
5. Batch prepare → Upsert to Upstash Vector
6. Log progress → Handle errors

## Key Design Decisions

### Chunking Strategy
- **Go files**: By function/struct/interface boundaries
- **Markdown**: By header hierarchy (H1, H2, etc.)
- **SQL**: By statement terminators
- **JSON/YAML**: By object/array boundaries
- **Text**: By paragraphs with semantic overlap

### Metadata Schema
```json
{
  "topic": "folder_name",
  "filename": "example.go",
  "filepath": "dev-docs/go-fiber-recipes/main.go",
  "extension": ".go",
  "chunk_index": 1,
  "total_chunks": 5,
  "content_type": "go_source",
  "created_at": "2025-10-18T..."
}
```

### Error Handling
- Continue processing on individual file errors
- Log errors with context
- Implement retry logic for network issues
- Track failed files for reprocessing

## Success Criteria
- All documents in dev-docs successfully ingested
- Chunks maintain semantic coherence
- Metadata accurately reflects document structure
- Process completes with clear logging
- Error recovery works as expected
- Performance acceptable for document volume

## Next Steps
1. Begin Phase 1: Analyze dev-docs structure and set up Upstash client
2. Implement document reading and basic chunking
3. Add metadata extraction and upsert functionality
4. Enhance with logging and error handling
5. Test and validate complete workflow