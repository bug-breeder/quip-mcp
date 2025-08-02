package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigManager_LoadSave(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create config manager with custom path
	cm := &ConfigManager{configPath: configPath}

	// Test saving and loading
	originalConfig := &Config{
		QuipAPIToken: "test-token-12345",
	}

	// Save config
	err := cm.Save(originalConfig)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedConfig, err := cm.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify
	if loadedConfig.QuipAPIToken != originalConfig.QuipAPIToken {
		t.Errorf("Expected token %s, got %s", originalConfig.QuipAPIToken, loadedConfig.QuipAPIToken)
	}
}

func TestConfigManager_EnvironmentOverride(t *testing.T) {
	// Create temporary directory for test
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Create config manager with custom path
	cm := &ConfigManager{configPath: configPath}

	// Save a config with one token
	fileConfig := &Config{
		QuipAPIToken: "file-token",
	}
	err := cm.Save(fileConfig)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Set environment variable
	envToken := "env-token"
	t.Setenv("QUIP_API_TOKEN", envToken)

	// Load config - should get env token
	loadedConfig, err := cm.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify environment variable takes precedence
	if loadedConfig.QuipAPIToken != envToken {
		t.Errorf("Expected env token %s, got %s", envToken, loadedConfig.QuipAPIToken)
	}
}

func TestConfigManager_NoConfigFile(t *testing.T) {
	// Create config manager with non-existent path
	cm := &ConfigManager{configPath: "/tmp/non-existent-config.yaml"}

	// Load config - should not fail even if file doesn't exist
	config, err := cm.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Should have empty token
	if config.QuipAPIToken != "" {
		t.Errorf("Expected empty token, got %s", config.QuipAPIToken)
	}
}

func TestConfigManager_HasValidToken(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected bool
	}{
		{
			name:     "valid token",
			token:    "valid-token-123",
			expected: true,
		},
		{
			name:     "empty token",
			token:    "",
			expected: false,
		},
		{
			name:     "short token",
			token:    "short",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")

			// Create config manager with custom path
			cm := &ConfigManager{configPath: configPath}

			// Save config if token is provided
			if tt.token != "" {
				config := &Config{QuipAPIToken: tt.token}
				err := cm.Save(config)
				if err != nil {
					t.Fatalf("Failed to save config: %v", err)
				}
			}

			// Test HasValidToken
			result := cm.HasValidToken()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetConfigPath(t *testing.T) {
	// Save original environment
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	originalHome := os.Getenv("HOME")

	// Clean up after test
	defer func() {
		if originalXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", originalXDG)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
		}
	}()

	tests := []struct {
		name        string
		xdgConfig   string
		homeDir     string
		expectedEnd string
	}{
		{
			name:        "XDG_CONFIG_HOME set",
			xdgConfig:   "/custom/config",
			homeDir:     "/home/user",
			expectedEnd: "/custom/config/quip-mcp/config.yaml",
		},
		{
			name:        "HOME directory fallback",
			xdgConfig:   "",
			homeDir:     "/home/user",
			expectedEnd: "/home/user/.config/quip-mcp/config.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.xdgConfig != "" {
				os.Setenv("XDG_CONFIG_HOME", tt.xdgConfig)
			} else {
				os.Unsetenv("XDG_CONFIG_HOME")
			}
			os.Setenv("HOME", tt.homeDir)

			// Get config path
			path := getConfigPath()

			// Check if path ends with expected suffix
			if path != tt.expectedEnd {
				t.Errorf("Expected path to be %s, got %s", tt.expectedEnd, path)
			}
		})
	}
}
