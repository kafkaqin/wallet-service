package redisx

import (
	"errors"
	"github.com/go-redis/redismock/v9"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
	"wallet-service/pkg/config"
)

func TestNewClientError(t *testing.T) {
	// Mock config to simulate Redis config in the environment
	mockConfig := &config.ServerConfig{
		Redis: config.Redis{
			Host:            "localhost",
			Port:            6379,
			PoolSize:        10,
			PoolTimeout:     5,
			MinIdleConns:    1,
			MaxIdleConns:    5,
			ConnMaxIdleTime: 10,
			Password:        "",
			DB:              0,
		},
	}
	// Temporary file path for testing
	tempFilePath := "test_config.yml"
	initialContent, err := yaml.Marshal(mockConfig)
	if err != nil {
		t.Error(err)
	}
	// Create a temporary config file
	err = os.WriteFile(tempFilePath, initialContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tempFilePath) // Clean up

	config.MustUseFileConfig[config.ServerConfig](tempFilePath)

	// Create a mock Redis client that will return an error on Ping
	_, mockRedis := redismock.NewClientMock()

	// Expecting Ping to return an error
	mockRedis.ExpectPing().SetErr(errors.New("PING failed"))

	// Call InitRedis, this should cause a panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic but did not occur")
		}
	}()

	// Initialize the Redis client which should panic
	InitRedis()
}
