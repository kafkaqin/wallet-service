package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"wallet-service/controllers"
	"wallet-service/pkg/config"
	"wallet-service/pkg/logger"
	"wallet-service/pkg/postgresx"
	"wallet-service/pkg/redisx"
	"wallet-service/services"
)

func main() {
	config.Init()
	postgresx.InitDB()
	redisx.InitRedis()
	l := logger.NewLogger()
	walletService := services.NewWalletService(l, postgresx.GetDB(), redisx.GetRedisClient())
	walletController := controllers.NewWalletController(l, redisx.GetRedisClient(), walletService)

	router := gin.Default()
	router.POST("/wallet/:user_id/deposit", walletController.Deposit)
	router.POST("/wallet/:user_id/withdraw", walletController.Withdraw)
	router.POST("/wallet/transfer/:sender_id/to/:receiver_id", walletController.Transfer)
	router.GET("/wallet/:user_id/balance", walletController.GetBalance)
	router.GET("/wallet/:user_id/transactions", walletController.GetTransactionHistory)

	err := router.Run(":8080") // 启动服务在8080端口(暂时不用配置文件里的端口)
	if err != nil {
		panic(err)
	}
	fmt.Println("server running on port 8080")
}
