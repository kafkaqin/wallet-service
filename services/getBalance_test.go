package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"wallet-service/pkg/logger"
)

func TestWalletService_GetBalance_SuccessFromRedis(t *testing.T) {
	// 创建 mock Redis 客户端
	client, mockRedis := redismock.NewClientMock()

	// 创建 mock DB 和 mock Logger
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")

	mockLogger := logger.NewLogger()

	// 创建 WalletService
	service := NewWalletService(mockLogger, sqlxDB, client)

	// 准备测试数据
	userID := 1
	expectedBalance := decimal.NewFromFloat(100.0).Round(0)

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", userID)).SetVal(expectedBalance.String())

	// 执行 GetBalance 方法
	balance, err := service.GetBalance(context.Background(), userID)

	// 断言没有错误，并且余额正确
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
}

func TestWalletService_GetBalance_SuccessFromDB(t *testing.T) {
	// 创建 mock Redis 客户端
	client, mockRedis := redismock.NewClientMock()

	// 创建 mock DB 和 mock Logger
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")

	mockLogger := logger.NewLogger()

	// 创建 WalletService
	service := NewWalletService(mockLogger, sqlxDB, client)

	// 准备测试数据
	userID := 1
	expectedBalance := decimal.NewFromFloat(100.0).Round(0)

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", userID)).SetErr(redis.Nil)

	// 设置 mock DB 的期望行为
	mockDB.ExpectQuery(`SELECT balance FROM wallets WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(expectedBalance))

	// 执行 GetBalance 方法
	balance, err := service.GetBalance(context.Background(), userID)

	// 断言没有错误，并且余额正确
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
}

func TestWalletService_GetBalance_RedisError(t *testing.T) {
	// 创建 mock Redis 客户端
	client, mockRedis := redismock.NewClientMock()

	// 创建 mock DB 和 mock Logger
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")

	mockLogger := logger.NewLogger()

	// 创建 WalletService
	service := NewWalletService(mockLogger, sqlxDB, client)

	// 准备测试数据
	userID := 1

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", userID)).SetErr(fmt.Errorf("redis error"))

	// 执行 GetBalance 方法
	_, err = service.GetBalance(context.Background(), userID)

	// 断言返回错误
	assert.Error(t, err)
}

func TestWalletService_GetBalance_DBError(t *testing.T) {
	// 创建 mock Redis 客户端
	client, mockRedis := redismock.NewClientMock()

	// 创建 mock DB 和 mock Logger
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")

	mockLogger := logger.NewLogger()

	// 创建 WalletService
	service := NewWalletService(mockLogger, sqlxDB, client)

	// 准备测试数据
	userID := 1

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", userID)).SetErr(redis.Nil)

	// 设置 mock DB 的期望行为
	mockDB.ExpectQuery(`SELECT balance FROM wallets WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(fmt.Errorf("database error"))

	// 执行 GetBalance 方法
	_, err = service.GetBalance(context.Background(), userID)

	// 断言返回错误
	assert.Error(t, err)
}

func TestWalletService_GetBalance_WalletNotFound(t *testing.T) {
	// 创建 mock Redis 客户端
	client, mockRedis := redismock.NewClientMock()

	// 创建 mock DB 和 mock Logger
	db, mockDB, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")

	mockLogger := logger.NewLogger()

	// 创建 WalletService
	service := NewWalletService(mockLogger, sqlxDB, client)

	// 准备测试数据
	userID := 1

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", userID)).SetErr(redis.Nil)

	// 设置 mock DB 的期望行为
	mockDB.ExpectQuery(`SELECT balance FROM wallets WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(errors.New("wallet not found"))

	// 执行 GetBalance 方法
	_, err = service.GetBalance(context.Background(), userID)

	// 断言返回错误
	assert.EqualError(t, err, "wallet not found")
}
