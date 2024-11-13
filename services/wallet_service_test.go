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
	"time"
	"wallet-service/models"
	"wallet-service/pkg/logger"
)

func TestWalletService_Deposit(t *testing.T) {
	// 创建 mock Redis 客户端
	client, mockRedis := redismock.NewClientMock()

	// 创建 mock DB 和 mock Logger
	db, mockDB, err := sqlmock.New() // 创建一个 mock DB 和 sqlmock 实例
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres") // 将 sqlmock 的 DB 传入 sqlx.DB

	mockLogger := logger.NewLogger()

	// 创建 WalletService
	service := NewWalletService(mockLogger, sqlxDB, client)

	// 准备测试数据
	senderID := 1
	receiverID := 2
	amount := decimal.NewFromFloat(100.0)
	transactionType := models.DepositTransactionType

	// 设置 mock DB 的期望行为
	mockDB.ExpectBegin() // 开始事务
	mockDB.ExpectExec(`UPDATE wallets SET balance = balance \+ \$1 WHERE user_id = \$2`).
		WithArgs(amount, senderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectExec("INSERT INTO transactions").WithArgs(senderID, receiverID, transactionType, amount, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1)) // 执行插入操作
	mockDB.ExpectCommit()                                                                                                                                             // 提交事务
	//ExpectRollback()
	// 设置 mock Redis 的期望行为
	mockRedis.ExpectSet(fmt.Sprintf("wallet:balance:%d", senderID), amount.String(), time.Duration(0)).SetVal("OK")

	// 执行 Deposit 方法
	err = service.Deposit(context.Background(), senderID, receiverID, amount, transactionType)

	// 断言没有错误
	assert.NoError(t, err)
}

func TestWalletService_Withdraw(t *testing.T) {
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

func TestWalletService_Transfer(t *testing.T) {
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
	amount := decimal.NewFromFloat(100.0)

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

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", senderID)).SetVal(amount.String())
	mockRedis.ExpectSet(fmt.Sprintf("wallet:balance:%d", senderID), amount.String(), time.Duration(0)).SetVal("OK")
	mockRedis.ExpectSet(fmt.Sprintf("wallet:balance:%d", receiverID), amount.String(), time.Duration(0)).SetVal("OK")

	// 执行 Transfer 方法
	err = service.Transfer(context.Background(), senderID, receiverID, amount)

	// 断言没有错误
	assert.NoError(t, err)
}

func TestWalletService_GetBalance(t *testing.T) {
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

	// 设置 mock DB 的期望行为
	mockDB.ExpectQuery(`SELECT balance FROM wallets WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(expectedBalance))

	// 设置 mock Redis 的期望行为
	mockRedis.ExpectGet(fmt.Sprintf("wallet:balance:%d", userID)).SetVal(expectedBalance.String())

	// 执行 GetBalance 方法
	balance, err := service.GetBalance(context.Background(), userID)

	// 断言没有错误，并且余额正确
	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
}

func TestWalletService_GetTransactionHistory(t *testing.T) {
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
	userID := 1
	expectedTransactions := []models.Transaction{
		{SenderUserID: 1, ReceiverUserID: 2, Amount: decimal.NewFromFloat(100.00).Round(2), TransactionType: models.DepositTransactionType},
	}

	// 设置 mock DB 的期望行为
	mockDB.ExpectQuery(`SELECT \* FROM transactions WHERE sender_user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"sender_user_id", "receiver_user_id", "amount", "transaction_type"}).
			AddRow(1, 2, "100.00", models.DepositTransactionType))

	// 执行 GetTransactionHistory 方法
	transactions, err := service.GetTransactionHistory(userID)

	// 断言没有错误，并且交易记录正确
	assert.NoError(t, err)
	assert.Equal(t, expectedTransactions, transactions)
}
