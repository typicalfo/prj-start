#!/bin/bash

# Test main MCP server with 2025 protocol version
echo "Testing main MCP server with 2025 protocol..."

# Create a simple test that sends initialization and then lists tools
(
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-06-18", "capabilities": {"tools": {}}, "clientInfo": {"name": "test-client", "version": "1.0.0"}}'
echo '{"jsonrpc": "2.0", "method": "notifications/initialized", "params": {}}'
sleep 1
echo '{"jsonrpc": "2.0", "id": 2, "method": "tools/list", "params": {}}'
) | ./main mcp --config test-config.yaml --debug

echo "Test completed."