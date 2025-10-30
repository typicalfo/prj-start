package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
	"github.com/typicalfo/prj-start/config"
	"github.com/upstash/vector-go"
)

// mcpCmd represents the mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server for querying vector database",
	Long: `Start the Model Context Protocol (MCP) server to enable agents
to query ingested documents in Upstash Vector database.

The MCP server provides:
- Natural language querying capabilities
- Vector similarity search
- Metadata-based filtering
- Real-time document retrieval

This command uses your existing Upstash Vector configuration.

Examples:
  prj-start mcp                    # Start MCP server
  prj-start mcp --debug           # Start with debug logging`,
	RunE: runMCP,
}

var (
	debug bool
)

func init() {
	rootCmd.AddCommand(mcpCmd)
	mcpCmd.Flags().BoolVar(&debug, "debug", false, "enable debug logging")
}

func runMCP(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check Upstash configuration
	if !cfg.HasUpstashConfig() {
		return fmt.Errorf("Upstash configuration incomplete\n\nUse 'prj-start init' to set up your configuration")
	}

	if debug {
		log.Printf("Loaded config: URL=%s, Token=%s", cfg.Upstash.URL, cfg.Upstash.Token)
	}

	// Create vector client
	vectorClient := createVectorClient(cfg)

	// Create MCP server
	server := createMCPServer()

	// Add tools
	addVectorQueryTool(server, vectorClient)
	addMetadataQueryTool(server, vectorClient)
	addListNamespacesTool(server, vectorClient)
	addGetDocumentTool(server, vectorClient)

	// Add resources
	addHelpResource(server)

	// Start server with stdio transport
	log.Println("Starting MCP server...")
	if debug {
		log.Println("Debug mode enabled")
	}

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("MCP server failed: %w", err)
	}

	return nil
}

func createMCPServer() *mcp.Server {
	return mcp.NewServer(&mcp.Implementation{
		Name:    "prj-start-vector-db",
		Version: "1.0.0",
	}, nil)
}

func createVectorClient(cfg *config.Config) *vector.Index {
	if cfg.Upstash.URL != "" && cfg.Upstash.Token != "" {
		if debug {
			log.Printf("Creating vector client with URL: %s", cfg.Upstash.URL)
		}
		return vector.NewIndex(cfg.Upstash.URL, cfg.Upstash.Token)
	}
	if debug {
		log.Println("Creating vector client from environment")
	}
	return vector.NewIndexFromEnv()
}
