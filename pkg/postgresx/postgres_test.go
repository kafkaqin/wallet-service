package postgresx

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDB(t *testing.T) {
	// Set up environment variables for DB config
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "wallet")
	os.Setenv("SSLMODE", "disable")

	// Initialize DB
	InitDB()

	// Test GetDB, should return the mock DB connection
	db := GetDB()
	assert.NotNil(t, db) // Ensure that DB is not nil

}
