package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/typicalfo/prj-start/config"
	"github.com/typicalfo/prj-start/processor"
)

var (
	ingestFolder string
)

// ingestCmd represents the ingest command
var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest documents into Upstash Vector database",
	Long: `Ingest documents from a specified folder into Upstash Vector database.
The tool processes documents, intelligently chunks them based on content type,
and upserts them with rich metadata for enhanced search and retrieval.

Examples:
  prj-start ingest                    # Ingest from current directory
  prj-start ingest --folder ./docs    # Ingest from specific folder
  prj-start ingest -f ./docs -v       # Ingest with verbose output`,
	RunE: runIngest,
}

func init() {
	rootCmd.AddCommand(ingestCmd)
	ingestCmd.Flags().StringVarP(&ingestFolder, "folder", "f", "", "folder to scan for documents (default is current directory)")
}

func runIngest(cmd *cobra.Command, args []string) error {
	// If no folder specified, use current directory
	if ingestFolder == "" {
		ingestFolder = "."
	}

	// Check if folder exists
	if _, err := os.Stat(ingestFolder); os.IsNotExist(err) {
		return fmt.Errorf("folder '%s' does not exist", ingestFolder)
	}

	// Load configuration
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w\n\nUse 'prj-start init' to set up your configuration", err)
	}

	// Check if Upstash configuration is complete
	if !cfg.HasUpstashConfig() {
		return fmt.Errorf("Upstash configuration is incomplete\n\nUse 'prj-start init' to set up your configuration")
	}

	// Validate folder
	if err := processor.ValidateFolder(ingestFolder); err != nil {
		return err
	}

	// Process documents
	ctx := context.Background()
	if err := processor.ProcessFolder(ctx, cfg, ingestFolder); err != nil {
		return fmt.Errorf("failed to process folder: %w", err)
	}

	return nil
}
