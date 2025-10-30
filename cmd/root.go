package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/typicalfo/prj-start/logger"
)

var (
	cfgFile string
	verbose bool
	version = "v0.1.5"
	commit  = "unknown"
	date    = "unknown"
	builtBy = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prj-start",
	Short: "Vector document processing and MCP server",
	Long: `prj-start is a Go-based tool for document processing and MCP server functionality.

Commands:
  ingest    - Process and ingest documents into Upstash Vector database
  init      - Initialize configuration
  mcp       - Start MCP server for querying (coming soon)

Use 'prj-start help <command>' for more information about a specific command.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s by: %s)", version, commit, date, builtBy),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			logger.SetLogLevel("debug")
		}
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
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
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
