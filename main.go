package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bug-breeder/quip-mcp/pkg/config"
	"github.com/bug-breeder/quip-mcp/pkg/server"
)

// Version information (set by build flags)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Command line flags
	var (
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show help information")
		setupConfig = flag.Bool("setup", false, "Run interactive configuration setup")
		showConfig  = flag.Bool("config", false, "Show current configuration")
		configPath  = flag.String("config-path", "", "Path to configuration file")
	)
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("quip-mcp %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	// Handle help flag
	if *showHelp {
		showUsage()
		os.Exit(0)
	}

	// Initialize config manager
	var configManager *config.ConfigManager
	if *configPath != "" {
		// TODO: Support custom config path
		configManager = config.New()
	} else {
		configManager = config.New()
	}

	// Handle setup flag
	if *setupConfig {
		if err := configManager.SetupInteractive(); err != nil {
			log.Fatalf("Configuration setup failed: %v", err)
		}
		os.Exit(0)
	}

	// Handle config display flag
	if *showConfig {
		showCurrentConfig(configManager)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := configManager.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Check if we have a valid token
	if cfg.QuipAPIToken == "" {
		fmt.Println("❌ No Quip API token found!")
		fmt.Println()
		fmt.Println("You can set up your token in one of these ways:")
		fmt.Println()
		fmt.Println("1. 🔧 Interactive setup (recommended):")
		fmt.Println("   quip-mcp --setup")
		fmt.Println()
		fmt.Println("2. 🌍 Environment variable:")
		fmt.Println("   export QUIP_API_TOKEN=\"your-token-here\"")
		fmt.Println("   quip-mcp")
		fmt.Println()
		fmt.Printf("3. 📁 Configuration file (%s):\n", configManager.GetConfigPath())
		fmt.Println("   quip_api_token: your-token-here")
		fmt.Println()
		fmt.Println("Get your token from: https://quip.com/dev/token")
		os.Exit(1)
	}

	// Start the MCP server
	srv := server.New(cfg.QuipAPIToken)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
}

func showUsage() {
	fmt.Println("Quip MCP Server")
	fmt.Println()
	fmt.Println("A Model Context Protocol server for Quip integration.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  quip-mcp [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -version       Show version information")
	fmt.Println("  -help          Show this help message")
	fmt.Println("  -setup         Run interactive configuration setup")
	fmt.Println("  -config        Show current configuration")
	fmt.Println("  -config-path   Path to configuration file")
	fmt.Println()
	fmt.Println("Configuration:")
	fmt.Println("  The server looks for your Quip API token in this order:")
	fmt.Println("  1. QUIP_API_TOKEN environment variable")
	fmt.Println("  2. Configuration file (~/.config/quip-mcp/config.yaml)")
	fmt.Println("  3. Interactive setup if no token found")
	fmt.Println()
	fmt.Println("Setup:")
	fmt.Println("  quip-mcp --setup     # Interactive token setup")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # First-time setup")
	fmt.Println("  quip-mcp --setup")
	fmt.Println()
	fmt.Println("  # Run server (after setup)")
	fmt.Println("  quip-mcp")
	fmt.Println()
	fmt.Println("  # Override with environment variable")
	fmt.Println("  export QUIP_API_TOKEN=\"your-token-here\"")
	fmt.Println("  quip-mcp")
	fmt.Println()
	fmt.Println("  # Show current configuration")
	fmt.Println("  quip-mcp --config")
	fmt.Println()
	fmt.Println("For more information, visit:")
	fmt.Println("  https://github.com/bug-breeder/quip-mcp")
}

func showCurrentConfig(configManager *config.ConfigManager) {
	fmt.Println("📋 Current Configuration")
	fmt.Println("========================")
	fmt.Println()

	cfg, err := configManager.Load()
	if err != nil {
		fmt.Printf("❌ Error loading configuration: %v\n", err)
		return
	}

	fmt.Printf("Config file: %s\n", configManager.GetConfigPath())
	fmt.Println()

	// Check if config file exists
	if _, err := os.Stat(configManager.GetConfigPath()); os.IsNotExist(err) {
		fmt.Println("📂 Config file: Not found")
	} else {
		fmt.Println("📂 Config file: Found")
	}

	// Check environment variable
	if envToken := os.Getenv("QUIP_API_TOKEN"); envToken != "" {
		fmt.Println("🌍 Environment variable: Set (will override config file)")
	} else {
		fmt.Println("🌍 Environment variable: Not set")
	}

	// Show token status (masked for security)
	if cfg.QuipAPIToken != "" {
		maskedToken := maskToken(cfg.QuipAPIToken)
		fmt.Printf("🔑 API Token: %s\n", maskedToken)
		fmt.Println("✅ Status: Ready to run")
	} else {
		fmt.Println("🔑 API Token: Not configured")
		fmt.Println("❌ Status: Setup required")
		fmt.Println()
		fmt.Println("Run 'quip-mcp --setup' to configure your API token.")
	}
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}
