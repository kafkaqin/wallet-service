package services

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"wallet-service/models"
	"wallet-service/pkg/logger"
)

func TestWalletService_Transfer_Success(t *testing.T) {
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
	mockRedis.ExpectSet(fmt.Sprintf("wallet:balance:%d", senderID), "50", 0).SetVal("OK")
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", receiverID)).SetVal("200")
	mockRedis.ExpectSet(fmt.Sprintf("wallet:balance:%d", receiverID), "250", 0).SetVal("OK")

	// 设置 mock DB 的期望行为
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE wallets SET balance = balance - \$1 WHERE user_id = \$2`).
		WithArgs(amount, senderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectExec(`UPDATE wallets SET balance = balance \+ \$1 WHERE user_id = \$2`).
		WithArgs(amount, receiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectExec("INSERT INTO transactions").WithArgs(receiverID, senderID, models.TransferTransactionType, amount, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()

	// 执行 Transfer 方法
	err = service.Transfer(context.Background(), senderID, receiverID, amount)

	// 断言没有错误
	assert.NoError(t, err)
}

func TestWalletService_Transfer_NegativeAmount(t *testing.T) {
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

	// 执行 Transfer 方法
	err = service.Transfer(context.Background(), senderID, receiverID, amount)

	// 断言返回错误
	assert.EqualError(t, err, "amount must be greater than zero")
}

func TestWalletService_Transfer_BeginTransactionError(t *testing.T) {
	// 创建 mock Redis 客户端
	client, _ := redismock.NewClientMock()

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

	// 设置 mock DB 的期望行为
	mockDB.ExpectBegin().WillReturnError(fmt.Errorf("begin transaction error"))

	// 执行 Transfer 方法
	err = service.Transfer(context.Background(), senderID, receiverID, amount)

	// 断言返回错误
	assert.Error(t, err)
}

func TestWalletService_Transfer_WithdrawError(t *testing.T) {
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
		WillReturnError(fmt.Errorf("withdraw error"))
	mockDB.ExpectRollback()

	// 执行 Transfer 方法
	err = service.Transfer(context.Background(), senderID, receiverID, amount)

	// 断言返回错误
	assert.Error(t, err)
}

func TestWalletService_Transfer_DepositError(t *testing.T) {
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
	mockRedis.ExpectSet(fmt.Sprintf("wallet:balance:%d", senderID), "50", 0).SetVal("OK")
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", receiverID)).SetVal("200")

	// 设置 mock DB 的期望行为
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE wallets SET balance = balance - \$1 WHERE user_id = \$2`).
		WithArgs(amount, senderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectExec(`UPDATE wallets SET balance = balance + \$1 WHERE user_id = \$2`).
		WithArgs(amount, receiverID).
		WillReturnError(fmt.Errorf("deposit error"))
	mockDB.ExpectRollback()

	// 执行 Transfer 方法
	err = service.Transfer(context.Background(), senderID, receiverID, amount)

	// 断言返回错误
	assert.Error(t, err)
}

func TestWalletService_Transfer_CommitError(t *testing.T) {
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
	mockRedis.ExpectSet(fmt.Sprintf("wallet:balance:%d", senderID), "50", 0).SetVal("OK")
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", receiverID)).SetVal("200")
	mockRedis.ExpectSet(fmt.Sprintf("wallet:balance:%d", receiverID), "250", 0).SetVal("OK")

	// 设置 mock DB 的期望行为
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE wallets SET balance = balance - \$1 WHERE user_id = \$2`).
		WithArgs(amount, senderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectExec(`UPDATE wallets SET balance = balance + \$1 WHERE user_id = \$2`).
		WithArgs(amount, receiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit().WillReturnError(fmt.Errorf("commit error"))
	mockDB.ExpectRollback()

	// 执行 Transfer 方法
	err = service.Transfer(context.Background(), senderID, receiverID, amount)

	// 断言返回错误
	assert.Error(t, err)
}
