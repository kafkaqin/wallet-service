package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"testing"
	"wallet-service/models"
	wallet_logger "wallet-service/pkg/logger"
)

func TestWalletController_GetTransactionHistory_Success(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100,
		MinIdleConns: 25,
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
	assert.Contains(t, w.Body.String(), `"data":`)
	assert.Contains(t, w.Body.String(), `"ID":1`)
	assert.Contains(t, w.Body.String(), `"ID":2`)

	// 验证方法调用
	mockService.AssertExpectations(t)
}

func TestWalletController_GetTransactionHistory_NoTransactions(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100,
		MinIdleConns: 25,
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
	transactions := []models.Transaction{}

	// 设置 mock WalletService 的期望行为，返回空的交易历史
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
	assert.Contains(t, w.Body.String(), `"error_code":100001`) // 未找到交易记录

	// 验证方法调用
	mockService.AssertExpectations(t)
}

func TestWalletController_GetTransactionHistory_InvalidParams(t *testing.T) {
	// 模拟 Redis 客户端
	ctx := context.Background()
	mockService := new(MockWalletService)
	var redisCli = redis.NewClient(&redis.Options{
		PoolSize:     100,
		MinIdleConns: 25,
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

	// 创建 HTTP 请求
	router := gin.Default()
	router.GET("/wallet/:user_id/transactions", controller.GetTransactionHistory)

	// 模拟请求，传递无效的 user_id
	req := httptest.NewRequest("GET", "/wallet/fs/transactions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":100004`) // 无效参数

	// 验证方法调用
	mockService.AssertExpectations(t)
}
