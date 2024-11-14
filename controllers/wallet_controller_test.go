package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet-service/models"
	wallet_logger "wallet-service/pkg/logger"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) Deposit(ctx context.Context, senderID, receiverID int, amount decimal.Decimal, transactionType models.TransactionType) error {
	args := m.Called(ctx, senderID, receiverID, amount, transactionType)
	return args.Error(0)
}

func (m *MockWalletService) Withdraw(ctx context.Context, senderID, receiverID int, amount decimal.Decimal, transactionType models.TransactionType) error {
	args := m.Called(ctx, senderID, receiverID, amount, transactionType)
	return args.Error(0)
}

func (m *MockWalletService) Transfer(ctx context.Context, senderID, receiverID int, amount decimal.Decimal) error {
	args := m.Called(ctx, senderID, receiverID, amount)
	return args.Error(0)
}

func (m *MockWalletService) GetBalance(ctx context.Context, userID int) (decimal.Decimal, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(decimal.Decimal), args.Error(1)
}

func (m *MockWalletService) GetTransactionHistory(userID int) ([]models.Transaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

// 测试 GetBalance 方法
func TestWalletController_GetBalance_Success1(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	// 模拟 WalletService
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100, // 最大链接数
		MinIdleConns: 25,  // 空闲链接数
		Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password:     "",
		DB:           0,
	})
	_, err := redisCli.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	l := wallet_logger.NewLogger()
	// 创建 WalletController
	controller := NewWalletController(l, redisCli, mockService)

	// 准备测试数据
	userID := 1
	expectedBalance := decimal.NewFromFloat(100.0)

	// 设置 mock WalletService 的期望行为
	mockService.On("GetBalance", mock.Anything, userID).Return(expectedBalance, nil)

	// 创建 HTTP 请求
	router := gin.Default()
	router.GET("/wallet/:user_id/balance", controller.GetBalance)

	// 模拟请求，传递有效的 user_id
	req := httptest.NewRequest("GET", fmt.Sprintf("/wallet/%d/balance", userID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"balance":"100"`) // 确保返回了正确的余额

	// 验证方法调用
	mockService.AssertExpectations(t)
}

func TestWalletController_GetBalance_InvalidUserID(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	// 模拟 WalletService
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100, // 最大链接数
		MinIdleConns: 25,  // 空闲链接数
		Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password:     "",
		DB:           0,
	})
	_, err := redisCli.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	l := wallet_logger.NewLogger()
	// 创建 WalletController
	controller := NewWalletController(l, redisCli, mockService)

	// 创建 HTTP 请求
	router := gin.Default()
	router.GET("/wallet/:user_id/balance", controller.GetBalance)

	// 模拟请求，传递有效的 user_id
	req := httptest.NewRequest("GET", fmt.Sprintf("/wallet/%s/balance", "fs"), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "invalid params")

	// 验证方法调用
	mockService.AssertExpectations(t)
}

func TestWalletController_Transfer_Success2(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	// 模拟 WalletService
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100, // 最大链接数
		MinIdleConns: 25,  // 空闲链接数
		Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password:     "",
		DB:           0,
	})
	_, err := redisCli.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	// 创建 WalletController
	l := wallet_logger.NewLogger()
	controller := NewWalletController(l, redisCli, mockService)

	// 模拟 WalletService 的 Transfer 方法
	mockService.On("Transfer", mock.Anything, 1, 2, mock.AnythingOfType("decimal.Decimal")).Return(nil)

	// 创建测试请求数据
	router := gin.Default()
	router.POST("/wallet/transfer/:sender_id/to/:receiver_id", controller.Transfer)

	// 模拟请求
	req := httptest.NewRequest("POST", "/wallet/transfer/1/to/2", strings.NewReader(`{"amount": 30.0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"Transfer successful"`)

	// 验证方法调用
	mockService.AssertExpectations(t)
}

func TestWalletController_Withdraw_Success(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	// 模拟 WalletService
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100, // 最大链接数
		MinIdleConns: 25,  // 空闲链接数
		Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password:     "",
		DB:           0,
	})
	_, err := redisCli.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	l := wallet_logger.NewLogger()
	controller := NewWalletController(l, redisCli, mockService)

	// 模拟 WalletService 的 Withdraw 方法
	mockService.On("Withdraw", mock.Anything, 1, 1, mock.AnythingOfType("decimal.Decimal"), models.WithdrawTransactionType).Return(nil)

	// 创建测试请求数据
	router := gin.Default()
	router.POST("/wallet/:user_id/withdraw", controller.Withdraw)

	// 模拟请求
	req := httptest.NewRequest("POST", "/wallet/1/withdraw", strings.NewReader(`{"amount": 50.0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"Withdraw successful"`)

	// 验证方法调用
	mockService.AssertExpectations(t)
}

func TestWalletController_Deposit_Failed1(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	// 模拟 WalletService
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100, // 最大链接数
		MinIdleConns: 25,  // 空闲链接数
		Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password:     "",
		DB:           0,
	})
	_, err := redisCli.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	l := wallet_logger.NewLogger()
	controller := NewWalletController(l, redisCli, mockService)

	// 准备测试数据
	userID := 1
	amount := decimal.NewFromFloat(100.0).Round(0)

	// 测试：存款操作失败
	t.Run("deposit failed", func(t *testing.T) {

		// 设置 mock WalletService 的期望行为，模拟存款失败
		mockService.On("Deposit", mock.Anything, userID, userID, amount, models.DepositTransactionType).Return(errors.New("internal server error"))

		// 创建 HTTP 请求
		router := gin.Default()
		router.POST("/wallet/:user_id/deposit", controller.Deposit)

		// 模拟请求，传递有效的 user_id 和存款金额
		req := httptest.NewRequest("POST", fmt.Sprintf("/wallet/%d/deposit", userID), nil)
		req.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"amount": %s}`, amount.String())))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 断言返回结果
		assert.Equal(t, 500, w.Code)
		assert.Contains(t, w.Body.String(), `"error_code":100003`) // 服务器错误

		// 验证方法调用
		mockService.AssertExpectations(t)
	})
}

func TestWalletController_Deposit_Success1(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	// 模拟 WalletService
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100, // 最大链接数
		MinIdleConns: 25,  // 空闲链接数
		Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password:     "",
		DB:           0,
	})
	_, err := redisCli.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	l := wallet_logger.NewLogger()
	controller := NewWalletController(l, redisCli, mockService)

	// 准备测试数据
	userID := 1
	amount := decimal.NewFromFloat(100.0).Round(0)

	// 测试：成功存款
	t.Run("success deposit", func(t *testing.T) {

		// 设置 mock WalletService 的期望行为
		mockService.On("Deposit", mock.Anything, userID, userID, amount, models.DepositTransactionType).Return(nil)

		// 创建 HTTP 请求
		router := gin.Default()
		router.POST("/wallet/:user_id/deposit", controller.Deposit)

		// 模拟请求，传递有效的 user_id 和存款金额
		req := httptest.NewRequest("POST", fmt.Sprintf("/wallet/%d/deposit", userID), strings.NewReader(`{"amount": `+amount.String()+`}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 断言返回结果
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), `"status":"Deposit successful"`)

		// 验证方法调用
		mockService.AssertExpectations(t)
	})
}

func TestWalletController_GetTransactionHistory(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	// 模拟 WalletService
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100, // 最大链接数
		MinIdleConns: 25,  // 空闲链接数
		Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password:     "",
		DB:           0,
	})
	_, err := redisCli.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	l := wallet_logger.NewLogger()
	controller := NewWalletController(l, redisCli, mockService)

	// 准备测试数据
	userID := 1
	transactions := []models.Transaction{
		{ID: 1, SenderUserID: userID, ReceiverUserID: 2, Amount: decimal.NewFromFloat(100.0).Round(0), TransactionType: models.DepositTransactionType},
		{ID: 2, SenderUserID: userID, ReceiverUserID: 3, Amount: decimal.NewFromFloat(50.0).Round(0), TransactionType: models.WithdrawTransactionType},
	}

	// 测试：成功获取交易历史
	t.Run("success get transaction history", func(t *testing.T) {
		// 设置 mock WalletService 的期望行为，返回交易历史
		mockService.On("GetTransactionHistory", userID).Return(transactions, nil)

		// 创建 HTTP 请求
		router := gin.Default()
		router.GET("/wallet/:user_id/transactions", controller.GetTransactionHistory)

		// 模拟请求，传递有效的 user_id
		req := httptest.NewRequest("GET", fmt.Sprintf("/wallet/%d/transactions", userID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 断言返回结果
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), `"ID":1`)
		assert.Contains(t, w.Body.String(), `"Amount":"100"`)

		// 验证方法调用
		mockService.AssertExpectations(t)
	})
}

func TestWalletController_GetTransactionHistory_Failed(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	// 模拟 WalletService
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100, // 最大链接数
		MinIdleConns: 25,  // 空闲链接数
		Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password:     "",
		DB:           0,
	})
	_, err := redisCli.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	// 创建 WalletController
	l := wallet_logger.NewLogger()
	controller := NewWalletController(l, redisCli, mockService)

	// 准备测试数据
	userID := 1

	// 测试：没有交易历史
	t.Run("no transaction history", func(t *testing.T) {
		// 设置 mock WalletService 的期望行为，返回空的交易历史
		mockService.On("GetTransactionHistory", userID).Return([]models.Transaction{}, nil)

		// 创建 HTTP 请求
		router := gin.Default()
		router.GET("/wallet/:user_id/transactions", controller.GetTransactionHistory)

		// 模拟请求，传递有效的 user_id
		req := httptest.NewRequest("GET", fmt.Sprintf("/wallet/%d/transactions", userID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 断言返回结果
		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), `"error_code":100001`) // 没有交易历史的错误码

		// 验证方法调用
		mockService.AssertExpectations(t)
	})
}
