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
	"testing"
	wallet_logger "wallet-service/pkg/logger"
)

func TestWalletController_GetBalance_Success(t *testing.T) {
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
	expectedBalance := decimal.NewFromFloat(100.0)

	// 设置 mock WalletService 的期望行为
	mockService.On("GetBalance", ctx, userID).Return(expectedBalance, nil)

	// 创建 HTTP 请求
	router := gin.Default()
	router.GET("/wallet/:user_id/balance", controller.GetBalance)

	// 模拟请求，传递有效的 user_id
	req := httptest.NewRequest("GET", fmt.Sprintf("/wallet/%d/balance", userID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), `"balance":"100"`)

	// 验证方法调用
	mockService.AssertExpectations(t)
}

func TestWalletController_GetBalance_Failed(t *testing.T) {
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

	// 设置 mock WalletService 的期望行为，模拟查询余额失败
	mockService.On("GetBalance", ctx, userID).Return(decimal.Zero, errors.New("internal server error"))

	// 创建 HTTP 请求
	router := gin.Default()
	router.GET("/wallet/:user_id/balance", controller.GetBalance)

	// 模拟请求，传递有效的 user_id
	req := httptest.NewRequest("GET", fmt.Sprintf("/wallet/%d/balance", userID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":100003`) // 服务器错误

	// 验证方法调用
	mockService.AssertExpectations(t)
}

func TestWalletController_GetBalance_InvalidParams(t *testing.T) {
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
	router.GET("/wallet/:user_id/balance", controller.GetBalance)

	// 模拟请求，传递无效的 user_id
	req := httptest.NewRequest("GET", "/wallet/fs/balance", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言返回结果
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), `"error_code":100004`) // 无效参数

	// 验证方法调用
	mockService.AssertExpectations(t)
}
