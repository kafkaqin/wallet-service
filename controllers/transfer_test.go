package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	wallet_logger "wallet-service/pkg/logger"
)

func TestWalletController_Transfer_Success(t *testing.T) {
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
	senderID := 1
	receiverID := 2
	amount := decimal.NewFromFloat(30.0).Round(0)

	// 设置 mock WalletService 的期望行为
	mockService.On("Transfer", ctx, senderID, receiverID, amount).Return(nil)

	// 创建 HTTP 请求
	router := gin.Default()
	router.POST("/wallet/transfer/:sender_id/to/:receiver_id", controller.Transfer)

	// 模拟请求，传递有效的 sender_id 和 receiver_id
	req := httptest.NewRequest("POST", fmt.Sprintf("/wallet/transfer/%d/to/%d", senderID, receiverID), strings.NewReader(fmt.Sprintf(`{"amount": "%s"}`, amount.String())))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"Transfer successful"`)

	// 验证方法调用
	mockService.AssertExpectations(t)
}

func TestWalletController_Transfer_Failed(t *testing.T) {
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
	senderID := 1
	receiverID := 2
	amount := decimal.NewFromFloat(30.0).Round(0)

	// 设置 mock WalletService 的期望行为，模拟转账失败
	mockService.On("Transfer", ctx, senderID, receiverID, amount).Return(errors.New("internal server error"))

	// 创建 HTTP 请求
	router := gin.Default()
	router.POST("/wallet/transfer/:sender_id/to/:receiver_id", controller.Transfer)

	// 模拟请求，传递有效的 sender_id 和 receiver_id
	req := httptest.NewRequest("POST", fmt.Sprintf("/wallet/transfer/%d/to/%d", senderID, receiverID), strings.NewReader(fmt.Sprintf(`{"amount": "%s"}`, amount.String())))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":100003`) // 服务器错误

	// 验证方法调用
	mockService.AssertExpectations(t)
}

func TestWalletController_Transfer_InvalidParams(t *testing.T) {
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
	router.POST("/wallet/transfer/:sender_id/to/:receiver_id", controller.Transfer)

	// 模拟请求，传递无效的 sender_id 或 receiver_id
	req := httptest.NewRequest("POST", "/wallet/transfer/fs/to/2", strings.NewReader(`{"amount": 30}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":100004`) // 无效参数

	// 验证方法调用
	mockService.AssertExpectations(t)
}
