package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigWatch(t *testing.T) {

	// Temporary file path for testing
	tempFilePath := "test_config.yml"

	// Prepare initial config content
	initialContent := `
wallet_service:
  host: "localhost"
  port: 8080
`
	// Create a temporary config file
	err := os.WriteFile(tempFilePath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tempFilePath) // Clean up

	// Initialize the config watcher
	cfg := MustUseFileConfig[ServerConfig](tempFilePath)

	// Simulate file content change
	updatedContent := `
wallet_service:
  host: "127.0.0.1"
  port: 9090
`
	// Write updated content to the file after a delay (to simulate file change)
	go func() {
		time.Sleep(2 * time.Second) // Simulate delay for file change
		err := os.WriteFile(tempFilePath, []byte(updatedContent), 0644)
		if err != nil {
			t.Errorf("Failed to write to temp config file: %v", err)
		}
	}()

	// Wait for the watcher to pick up the file change
	time.Sleep(3 * time.Second)

	// Get the current config value after update
	configValue := cfg.Get()

	// Verify if the config was updated properly
	assert.Equal(t, "127.0.0.1", configValue.WalletService.Host)
	assert.Equal(t, 9090, configValue.WalletService.Port)

}

func TestConfigWatchError(t *testing.T) {

	// Temporary file path for testing
	tempFilePath := "test_config_invalid.yml"

	// Create an invalid config file
	err := os.WriteFile(tempFilePath, []byte("invalid: yaml"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tempFilePath) // Clean up

	// Replace logger in config package

	// Initialize the config watcher (this should trigger an error)
	cfg := MustUseFileConfig[ServerConfig](tempFilePath)
	// Simulate file content change
	updatedContent := `
server:
  host: "127.0.0.1"
  port: 9090
`
	// Write updated content to the file after a delay (to simulate file change)
	go func() {
		time.Sleep(2 * time.Second) // Simulate delay for file change
		err := os.WriteFile(tempFilePath, []byte(updatedContent), 0644)
		if err != nil {
			t.Errorf("Failed to write to temp config file: %v", err)
		}
	}()
	// Verify if the config was updated properly
	configValue := cfg.Get()
	assert.NotEqual(t, "127.0.0.1", configValue.WalletService.Host)
	assert.NotEqual(t, 9090, configValue.WalletService.Port)
	// Wait for the watcher to pick up the file change
	time.Sleep(3 * time.Second)

}

func TestStartWatch(t *testing.T) {
	// Temporary file path for testing
	tempFilePath := "test_config.yml"

	// Prepare initial config content
	initialContent := `
wallet_service:
  host: "localhost"
  port: 8080
`
	// Create a temporary config file
	err := os.WriteFile(tempFilePath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tempFilePath) // Clean up

	// Initialize the config watcher
	cfg := NewCfg[ServerConfig]()
	err = cfg.StartWatch(tempFilePath)
	if err != nil {
		t.Fatalf("Failed to start watch: %v", err)
	}

	// Verify initial config values
	configValue := cfg.Get()
	assert.Equal(t, "localhost", configValue.WalletService.Host)
	assert.Equal(t, 8080, configValue.WalletService.Port)

	// Simulate file content change
	updatedContent := `
wallet_service:
  host: "127.0.0.1"
  port: 9090
`
	// Write updated content to the file after a delay (to simulate file change)
	go func() {
		time.Sleep(2 * time.Second) // Simulate delay for file change
		err := os.WriteFile(tempFilePath, []byte(updatedContent), 0644)
		if err != nil {
			t.Errorf("Failed to write to temp config file: %v", err)
		}
	}()

	// Wait for the watcher to pick up the file change
	time.Sleep(3 * time.Second)

	// Get the current config value after update
	configValue = cfg.Get()

	// Verify if the config was updated properly
	assert.Equal(t, "127.0.0.1", configValue.WalletService.Host)
	assert.Equal(t, 9090, configValue.WalletService.Port)
}

func TestStartWatchWithError(t *testing.T) {
	// Temporary file path for testing
	tempFilePath := "test_config_invalid.yml"

	// Create an invalid config file
	err := os.WriteFile(tempFilePath, []byte("invalid: yaml"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tempFilePath) // Clean up

	// Initialize the config watcher (this should trigger an error)
	cfg := NewCfg[ServerConfig]()
	err = cfg.StartWatch(tempFilePath)
	if err != nil {
		t.Errorf("Expected an error but got none")
	}

	// Verify that the config is not updated
	configValue := cfg.Get()
	assert.Equal(t, "", configValue.WalletService.Host)
	assert.Equal(t, 0, configValue.WalletService.Port)

	// Simulate file content change
	updatedContent := `
wallet_service:
  host: "127.0.0.1"
  port: 9090
`
	// Write updated content to the file after a delay (to simulate file change)
	go func() {
		time.Sleep(2 * time.Second) // Simulate delay for file change
		err := os.WriteFile(tempFilePath, []byte(updatedContent), 0644)
		if err != nil {
			t.Errorf("Failed to write to temp config file: %v", err)
		}
	}()

	// Wait for the watcher to pick up the file change
	time.Sleep(3 * time.Second)

	// Get the current config value after update
	configValue = cfg.Get()

	// Verify if the config was updated properly
	assert.Equal(t, "127.0.0.1", configValue.WalletService.Host)
	assert.Equal(t, 9090, configValue.WalletService.Port)
}

func TestInit(t *testing.T) {
	// Temporary file path for testing
	tempFilePath := "test_config_init.yml"

	// Prepare initial config content
	initialContent := `
wallet_service:
  host: "localhost"
  port: 8080
`
	// Create a temporary config file
	err := os.WriteFile(tempFilePath, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tempFilePath) // Clean up

	// Initialize the config with custom path
	Init(tempFilePath)

	// Verify initial config values
	configValue := GetConfig()
	assert.Equal(t, "localhost", configValue.WalletService.Host)
	assert.Equal(t, 8080, configValue.WalletService.Port)

	// Simulate file content change
	updatedContent := `
wallet_service:
  host: "127.0.0.1"
  port: 9090
`
	// Write updated content to the file after a delay (to simulate file change)
	go func() {
		time.Sleep(2 * time.Second) // Simulate delay for file change
		err := os.WriteFile(tempFilePath, []byte(updatedContent), 0644)
		if err != nil {
			t.Errorf("Failed to write to temp config file: %v", err)
		}
	}()

	// Wait for the watcher to pick up the file change
	time.Sleep(3 * time.Second)

	// Get the current config value after update
	configValue = GetConfig()

	// Verify if the config was updated properly
	assert.Equal(t, "127.0.0.1", configValue.WalletService.Host)
	assert.Equal(t, 9090, configValue.WalletService.Port)
}
