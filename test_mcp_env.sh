#!/bin/bash

# Test MCP server without config to isolate the issue
echo "Testing MCP server without config..."

# Test with environment variables instead
export UPSTASH_VECTOR_REST_URL="https://test-vector.upstash.io"
export UPSTASH_VECTOR_REST_TOKEN="test-token"

# Create a simple test that sends initialization and then lists tools
(
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-06-18", "capabilities": {"tools": {}}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}'
echo '{"jsonrpc": "2.0", "method": "notifications/initialized", "params": {}}'
sleep 1
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}'
) | ./main mcp --debug

echo "Test completed."