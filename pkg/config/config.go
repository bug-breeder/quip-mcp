package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	QuipAPIToken string `json:"quip_api_token" yaml:"quip_api_token"`
}

// ConfigManager handles loading and saving configuration
type ConfigManager struct {
	configPath string
}

// New creates a new ConfigManager
func New() *ConfigManager {
	return &ConfigManager{
		configPath: getConfigPath(),
	}
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	// Try XDG_CONFIG_HOME first (Linux/Unix standard)
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "quip-mcp", "config.yaml")
	}

	// Fall back to home directory
	home, err := os.UserHomeDir()
	if err != nil {
		// If we can't get home directory, use current directory
		return ".quip-mcp-config.yaml"
	}

	// Use .config in home directory (follows XDG on Linux, reasonable on other platforms)
	return filepath.Join(home, ".config", "quip-mcp", "config.yaml")
}

// Load loads configuration from file and environment
func (cm *ConfigManager) Load() (*Config, error) {
	config := &Config{}

	// First, try to load from config file
	if err := cm.loadFromFile(config); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	// Override with environment variable if set
	if token := os.Getenv("QUIP_API_TOKEN"); token != "" {
		config.QuipAPIToken = token
	}

	return config, nil
}

// loadFromFile loads configuration from the config file
func (cm *ConfigManager) loadFromFile(config *Config) error {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return err
	}

	// Try to parse as YAML first, then JSON
	if err := yaml.Unmarshal(data, config); err != nil {
		// If YAML fails, try JSON
		if jsonErr := json.Unmarshal(data, config); jsonErr != nil {
			return fmt.Errorf("failed to parse config as YAML or JSON: %w", err)
		}
	}

	return nil
}

// Save saves the configuration to file
func (cm *ConfigManager) Save(config *Config) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with restrictive permissions (only user can read/write)
	if err := os.WriteFile(cm.configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigPath returns the path to the configuration file
func (cm *ConfigManager) GetConfigPath() string {
	return cm.configPath
}

// SetupInteractive prompts the user to configure the API token
func (cm *ConfigManager) SetupInteractive() error {
	fmt.Println("🔧 Quip MCP Server Setup")
	fmt.Println("========================")
	fmt.Println()
	fmt.Println("To use the Quip MCP server, you need a Quip API token.")
	fmt.Println("You can get one from: https://quip.com/dev/token")
	fmt.Println()

	// Prompt for token
	fmt.Print("Enter your Quip API token: ")
	token, err := readPassword()
	if err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Validate token format (basic check)
	if len(token) < 10 {
		return fmt.Errorf("token appears to be too short, please check and try again")
	}

	// Save configuration
	config := &Config{
		QuipAPIToken: token,
	}

	if err := cm.Save(config); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println()
	fmt.Printf("✅ Configuration saved to: %s\n", cm.configPath)
	fmt.Println("🚀 You can now run 'quip-mcp' to start the server!")
	fmt.Println()
	fmt.Println("Note: You can override this token anytime by setting the QUIP_API_TOKEN environment variable.")

	return nil
}

// readPassword reads a password from stdin without echoing
func readPassword() (string, error) {
	// Check if we're in a terminal
	if term.IsTerminal(int(syscall.Stdin)) {
		// Read password without echo
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", err
		}
		fmt.Println() // Print newline after password input
		return string(bytePassword), nil
	}

	// If not in terminal, read normally (for testing/automation)
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(password, "\n"), nil
}

// HasValidToken checks if a valid token is available
func (cm *ConfigManager) HasValidToken() bool {
	config, err := cm.Load()
	if err != nil {
		return false
	}
	return config.QuipAPIToken != "" && len(config.QuipAPIToken) >= 10
}
