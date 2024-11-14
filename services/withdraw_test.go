package services

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wallet-service/models"
	"wallet-service/pkg/logger"
)

func TestWalletService_Withdraw_Success(t *testing.T) {
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
	senderID := 1
	receiverID := 2
	amount := decimal.NewFromFloat(50.0)
	transactionType := models.WithdrawTransactionType

	// 设置 mock DB 的期望行为
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE wallets SET balance = balance - \$1 WHERE user_id = \$2`).
		WithArgs(amount, senderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectExec("INSERT INTO transactions").WithArgs(senderID, receiverID, transactionType, amount, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()

	// 设置 mock Redis 的期望行为
	// 首先 GET 操作（查询余额）
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", senderID)).SetVal("100")

	// 然后 SET 操作（更新余额）
	mockRedis.ExpectSet(fmt.Sprintf("wallet:balance:%d", senderID), "50", time.Duration(0)).SetVal("OK")

	// 执行 Withdraw 方法
	err = service.Withdraw(context.Background(), senderID, receiverID, amount, transactionType)

	// 断言没有错误
	assert.NoError(t, err)
}

func TestWalletService_Withdraw_InsufficientFunds(t *testing.T) {
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
	senderID := 1
	receiverID := 2
	amount := decimal.NewFromFloat(150.0)

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", senderID)).SetVal("100")

	// 执行 Withdraw 方法
	err = service.Withdraw(context.Background(), senderID, receiverID, amount, "")

	// 断言返回错误
	assert.EqualError(t, err, "insufficient funds balance")
}

func TestWalletService_Withdraw_RedisError(t *testing.T) {
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
	senderID := 1
	receiverID := 2
	amount := decimal.NewFromFloat(50.0)

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", senderID)).SetErr(fmt.Errorf("redis error"))

	// 执行 Withdraw 方法
	err = service.Withdraw(context.Background(), senderID, receiverID, amount, "")

	// 断言返回错误
	assert.Error(t, err)
}

func TestWalletService_Withdraw_DBError(t *testing.T) {
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
	senderID := 1
	receiverID := 2
	amount := decimal.NewFromFloat(50.0)

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", senderID)).SetVal("100")

	// 设置 mock DB 的期望行为
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE wallets SET balance = balance - \$1 WHERE user_id = \$2`).
		WithArgs(amount, senderID).
		WillReturnError(fmt.Errorf("database error"))
	mockDB.ExpectRollback()

	// 执行 Withdraw 方法
	err = service.Withdraw(context.Background(), senderID, receiverID, amount, "")

	// 断言返回错误
	assert.Error(t, err)
}

func TestWalletService_Withdraw_RedisUpdateError(t *testing.T) {
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
	senderID := 1
	receiverID := 2
	amount := decimal.NewFromFloat(50.0)

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", senderID)).SetVal("100")
	mockRedis.ExpectSet(fmt.Sprintf("wallet:balance:%d", senderID), "50", time.Duration(0)).SetErr(fmt.Errorf("redis update error"))

	// 设置 mock DB 的期望行为
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE wallets SET balance = balance - \$1 WHERE user_id = \$2`).
		WithArgs(amount, senderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectExec("INSERT INTO transactions").WithArgs(senderID, receiverID, models.WithdrawTransactionType, amount, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()

	// 执行 Withdraw 方法
	err = service.Withdraw(context.Background(), senderID, receiverID, amount, "")

	// 断言没有错误
	assert.NoError(t, err)
}

func TestWalletService_Withdraw_InvalidParams(t *testing.T) {
	// 创建 mock Redis 客户端
	client, _ := redismock.NewClientMock()

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
	senderID := 1
	receiverID := 2
	amount := decimal.NewFromFloat(-50.0)

	// 执行 Withdraw 方法
	err = service.Withdraw(context.Background(), senderID, receiverID, amount, models.WithdrawTransactionType)

	// 断言返回错误
	assert.EqualError(t, err, "amount must be greater than zero")
}

func TestWalletService_Withdraw_DBSelectError(t *testing.T) {
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
	senderID := 1
	receiverID := 2
	amount := decimal.NewFromFloat(50.0)

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", senderID)).SetErr(redis.Nil)

	// 设置 mock DB 的期望行为
	mockDB.ExpectQuery(`SELECT balance FROM wallets WHERE user_id = \$1`).
		WithArgs(senderID).
		WillReturnError(fmt.Errorf("database select error"))

	// 执行 Withdraw 方法
	err = service.Withdraw(context.Background(), senderID, receiverID, amount, "")

	// 断言返回错误
	assert.Error(t, err)
}
