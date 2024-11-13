package postgresx

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/cast"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"wallet-service/pkg/config"
)

var oncePostgres = sync.Once{}
var PostgresClient *sqlx.DB

// InitDB init PostgresClient
func InitDB() *sqlx.DB {
	oncePostgres.Do(func() {
		if PostgresClient == nil {
			cfg := config.GetConfig().Postgres

			dbHost := os.Getenv("DB_HOST")
			if dbHost != "" {
				cfg.Host = dbHost
			}

			dbPort := os.Getenv("DB_PORT")
			if dbPort != "" {
				cfg.Port = cast.ToInt(dbPort)
			}

			dbUser := os.Getenv("DB_USER")
			if dbUser != "" {
				cfg.User = dbUser
			}

			dbPassword := os.Getenv("DB_PASSWORD")
			if dbPassword != "" {
				cfg.Password = dbPassword
			}

			dbName := os.Getenv("DB_NAME")
			if dbName != "" {
				cfg.Database = dbName
			}

			SSLMode := os.Getenv("SSLMODE")
			if SSLMode != "" {
				cfg.SSLMode = SSLMode
			}
			connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
			PostgresClient = sqlx.MustConnect("postgres", connStr)

			dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password,
				cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode)

			migrateDB(dataSourceName)
		}
	})

	return PostgresClient
}

// GetProjectRoot 获取项目根目录的绝对路径
func GetProjectRoot() (string, error) {
	// 获取当前文件的路径
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get caller's path")
	}
	// 通过上级目录找到项目根目录
	projectRoot := filepath.Join(filepath.Dir(currentFilePath), "..", "..")

	// 转换为绝对路径
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("unable to get absolute path: %v", err)
	}

	return absProjectRoot, nil
}

// GetDB is get postgres instance
func GetDB() *sqlx.DB {
	if PostgresClient == nil {
		return InitDB()
	}
	return PostgresClient
}
