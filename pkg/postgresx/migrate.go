package postgresx

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func migrateDB(dataSourceName string) {
	// 连接数据库
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	// 初始化迁移
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not start database migration: %v", err)
	}

	absProjectRoot, err := GetProjectRoot()
	if err != nil {
		log.Fatalf("Could not get absProjectRoot : %v", err)
	}

	migrationsPath := filepath.Join(absProjectRoot, "pkg", "postgresx", "migrations")
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 执行迁移
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration up failed: %v", err)
	}
}
