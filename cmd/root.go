package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/typicalfo/prj-start/config"
	"github.com/typicalfo/prj-start/logger"
	"github.com/typicalfo/prj-start/processor"
)

var (
	cfgFile string
	folder  string
	verbose bool
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prj-start",
	Short: "Vector document ingestion system for development documentation",
	Long: `prj-start is a Go-based tool for intelligently chunking and ingesting 
development documentation into Upstash Vector database for enhanced search and retrieval.

The tool processes documents from any specified folder, automatically chunks them 
based on content type, and upserts them to Upstash Vector with rich metadata.

Perfect for:
- Development teams with code documentation
- API documentation and examples  
- Recipe collections and tutorials
- Knowledge base management

Use 'prj-start init' to set up your configuration first.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s by: %s)", version, commit, date, builtBy),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			logger.SetLogLevel("debug")
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no folder specified, use current directory
		if folder == "" {
			folder = "."
		}

		// Check if folder exists
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			return fmt.Errorf("folder '%s' does not exist", folder)
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
		if err := processor.ValidateFolder(folder); err != nil {
			return err
		}

		// Process documents
		ctx := context.Background()
		if err := processor.ProcessFolder(ctx, cfg, folder); err != nil {
			return fmt.Errorf("failed to process folder: %w", err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/prj-start/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&folder, "folder", "f", "", "folder to scan for documents (default is current directory)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Add custom help template
	rootCmd.SetHelpTemplate(`{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)

	// Add completion command
	rootCmd.AddCommand(completionCmd)
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `To load completions:

Bash:
  $ source <(prj-start completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ prj-start completion bash > /etc/bash_completion.d/prj-start
  # macOS:
  $ prj-start completion bash > /usr/local/etc/bash_completion.d/prj-start

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ prj-start completion zsh > "${fpath[1]}/_prj-start"

  # You will need to start a new shell for this setup to take effect.

fish:
  $ prj-start completion fish | source

  # To load completions for each session, execute once:
  $ prj-start completion fish > ~/.config/fish/completions/prj-start.fish

PowerShell:
  PS> prj-start completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> prj-start completion powershell > prj-start.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}
