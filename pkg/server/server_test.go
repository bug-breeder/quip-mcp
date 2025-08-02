package server

import (
	"testing"
)

func TestNew(t *testing.T) {
	token := "test-token"
	server := New(token)

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.mcpServer == nil {
		t.Fatal("Expected MCP server to be initialized, got nil")
	}

	if server.quipClient == nil {
		t.Fatal("Expected Quip client to be initialized, got nil")
	}
}

func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		expected  string
	}{
		{
			name:      "zero timestamp",
			timestamp: 0,
			expected:  "Unknown",
		},
		{
			name:      "valid timestamp",
			timestamp: 1640995200000000, // 2022-01-01 00:00:00 UTC in microseconds
			expected:  "1640995200",     // Same in seconds
		},
		{
			name:      "another valid timestamp",
			timestamp: 1609459200000000, // 2021-01-01 00:00:00 UTC in microseconds
			expected:  "1609459200",     // Same in seconds
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimestamp(tt.timestamp)
			if result != tt.expected {
				t.Errorf("formatTimestamp(%d) = %s, expected %s", tt.timestamp, result, tt.expected)
			}
		})
	}
}

func TestTruncateText(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		maxLength int
		expected  string
	}{
		{
			name:      "text shorter than limit",
			text:      "Hello",
			maxLength: 10,
			expected:  "Hello",
		},
		{
			name:      "text equal to limit",
			text:      "Hello World",
			maxLength: 11,
			expected:  "Hello World",
		},
		{
			name:      "text longer than limit",
			text:      "This is a very long text that should be truncated",
			maxLength: 20,
			expected:  "This is a very long...",
		},
		{
			name:      "text with spaces at end",
			text:      "Hello World   ",
			maxLength: 10,
			expected:  "Hello Worl...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateText(tt.text, tt.maxLength)
			if result != tt.expected {
				t.Errorf("truncateText(%q, %d) = %q, expected %q", tt.text, tt.maxLength, result, tt.expected)
			}
		})
	}
}

// TestQuipClientIntegration tests that the server correctly integrates with the Quip client
func TestQuipClientIntegration(t *testing.T) {
	// This is an integration test that verifies the server properly wraps the Quip client
	token := "test-token"
	server := New(token)

	// Verify that the Quip client is properly configured
	if server.quipClient == nil {
		t.Fatal("Expected Quip client to be initialized")
	}

	// We can't test the actual API calls without a real token and internet connection,
	// but we can verify the structure is correct
}
