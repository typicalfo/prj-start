# Project Refactoring Plan: Command Structure Reorganization

## Current State Analysis
The project currently has a basic command structure with:
- `main.go` - Entry point that calls `cmd.Execute()`
- `cmd/root.go` - Root command with ingest logic in the RunE function
- `cmd/init.go` - Configuration setup command

## Problem Statement
The current structure has the document ingestion logic embedded directly in the root command's RunE function. This makes it difficult to:
1. Add new commands (like the upcoming MCP server)
2. Maintain separation of concerns
3. Test individual functionality
4. Extend the application with additional features

## Refactoring Goals
1. Move ingestion logic to a dedicated command
2. Clean up the root command to be a proper coordinator
3. Prepare structure for adding MCP server command
4. Maintain backward compatibility
5. Follow Go command-line application best practices

## Implementation Plan

### Phase 1: Create Ingest Command
1. Create `cmd/ingest.go` with dedicated ingest command
2. Move ingestion logic from `root.go:RunE` to the new ingest command
3. Update root command to show help when no subcommand is provided
4. Add ingest command to root command's subcommands

### Phase 2: Update Root Command
1. Remove the ingestion logic from root command's RunE
2. Set root command to display help/usage when called directly
3. Ensure all flags are properly migrated to ingest command
4. Update command descriptions and help text

### Phase 3: Prepare for MCP Server
1. Create placeholder `cmd/mcp.go` for future MCP server command
2. Update root command help to mention upcoming MCP functionality
3. Ensure configuration supports both ingest and MCP use cases

### Phase 4: Testing and Validation
1. Test all existing functionality works with new structure
2. Verify backward compatibility
3. Update documentation and help text
4. Run full test suite

## File Changes Required

### New Files
- `cmd/ingest.go` - New dedicated ingest command
- `cmd/mcp.go` - Placeholder for MCP server command

### Modified Files
- `cmd/root.go` - Remove ingestion logic, update to coordinator role
- `main.go` - No changes needed
- Documentation files to reflect new command structure

## Backward Compatibility
- Default behavior (`prj-start`) will show help and suggest `prj-start ingest`
- `prj-start --folder /path` will become `prj-start ingest --folder /path`
- All existing flags will be available on the ingest command
- Configuration and init command remain unchanged

## Success Criteria
1. All existing functionality preserved
2. Clean separation between commands
3. Easy to add new commands
4. Clear help and usage information
5. No breaking changes to configuration

## Next Steps After Refactoring
1. Implement MCP server command in `cmd/mcp.go`
2. Add vector database query functionality
3. Implement MCP protocol handlers
4. Add comprehensive testing for MCP features