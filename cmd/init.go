package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/typicalfo/prj-start/config"
)

var (
	force      bool
	configPath string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration for prj-start",
	Long: `Initialize configuration for prj-start by setting up Upstash Vector credentials
and other settings. This command will guide you through the setup process interactively.

The configuration will be saved to:
- ~/.config/prj-start/config.yaml (Linux/macOS)
- %LOCALAPPDATA%/prj-start/config.yaml (Windows)
- Or ./prj-start.yaml if no home directory is available

You can also specify a custom config file with --config flag.`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&force, "force", "F", false, "overwrite existing configuration")
	initCmd.Flags().StringVar(&configPath, "config", "", "config file path (default is standard location)")
}

func runInit(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🚀 Welcome to prj-start Configuration Setup")
	fmt.Println("==========================================")
	fmt.Println()

	// Determine config file path
	targetPath := configPath
	if targetPath == "" {
		targetPath = config.GetDefaultConfigPath()
	}

	// Check if config already exists
	if _, err := os.Stat(targetPath); err == nil && !force {
		fmt.Printf("Configuration file already exists at: %s\n", targetPath)
		fmt.Print("Do you want to overwrite it? (y/N): ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Configuration setup cancelled.")
			return nil
		}
	}

	fmt.Printf("Configuration will be saved to: %s\n", targetPath)
	fmt.Println()

	// Create new config
	cfg := &config.Config{
		DefaultNamespace: "default",
		BatchSize:        10,
		LogLevel:         "info",
	}

	// Get Upstash configuration
	fmt.Println("📋 Upstash Vector Configuration")
	fmt.Println("-------------------------------")
	fmt.Println("You can get these values from: https://console.upstash.com")
	fmt.Println()

	cfg.Upstash.URL = promptForInput(reader, "Upstash Vector REST URL", "https://your-vector-url.upstash.io", true)
	cfg.Upstash.Token = promptForInput(reader, "Upstash Vector REST Token", "", true)
	cfg.Upstash.IndexURL = promptForInput(reader, "Upstash Vector Index URL (optional)", "", false)

	fmt.Println()
	fmt.Println("🔧 Optional Configuration")
	fmt.Println("------------------------")

	cfg.Upstash.Email = promptForInput(reader, "Upstash Email (for MCP server)", "", false)
	cfg.Upstash.APIKey = promptForInput(reader, "Upstash API Key (for MCP server)", "", false)

	fmt.Println()
	cfg.DefaultNamespace = promptForInput(reader, "Default namespace", "default", false)

	logLevel := promptForInput(reader, "Log level (debug, info, warn, error)", "info", false)
	if logLevel != "" {
		cfg.LogLevel = logLevel
	}

	// Save configuration
	if err := config.SaveConfig(cfg, targetPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ Configuration saved successfully!")
	fmt.Printf("📁 Config file: %s\n", targetPath)
	fmt.Println()
	fmt.Println("🎯 Next steps:")
	fmt.Printf("1. Process documents: prj-start --folder /path/to/docs\n")
	fmt.Printf("2. Or use current directory: prj-start\n")
	fmt.Printf("3. For help: prj-start --help\n")
	fmt.Println()
	fmt.Println("💡 Tip: You can override any setting with environment variables:")
	fmt.Println("   export UPSTASH_VECTOR_REST_URL=your-url")
	fmt.Println("   export UPSTASH_VECTOR_REST_TOKEN=your-token")

	return nil
}

func promptForInput(reader *bufio.Reader, prompt, defaultValue string, required bool) string {
	for {
		if defaultValue != "" {
			fmt.Printf("%s [%s]: ", prompt, defaultValue)
		} else {
			fmt.Printf("%s: ", prompt)
		}

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			input = defaultValue
		}

		if required && input == "" {
			fmt.Println("⚠️  This field is required. Please provide a value.")
			continue
		}

		return input
	}
}
