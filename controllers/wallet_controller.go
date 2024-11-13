package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"strconv"
	"wallet-service/models"
	wallet_logger "wallet-service/pkg/logger"
	"wallet-service/pkg/rdsLimit"
	"wallet-service/services"
)

const limitCount = 100

type WalletController struct {
	walletService services.WalletService
	redis         *redis.Client
	logger        *wallet_logger.Logger
}

// NewWalletController new wallet controller
func NewWalletController(logger *wallet_logger.Logger, redis *redis.Client, service services.WalletService) *WalletController {
	return &WalletController{
		walletService: service,
		redis:         redis,
		logger:        logger,
	}
}

func generateTraceID() string {
	// 生成一个新的 UUID 作为 Trace ID
	traceID := uuid.New()
	return traceID.String()
}

// Deposit 存款
func (wc *WalletController) Deposit(c *gin.Context) {
	wc.logger.WithField("traceID", generateTraceID())
	ctx := c.Request.Context()
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		wc.logger.Error(ctx, "WalletController Deposit strconv.Atoi",
			zap.Int("userID", userID), zap.Error(err))
		handleError(c, CODE_INVALID_PARAMS, err)
		return
	}
	if !rdsLimit.NewRdsLimit(wc.redis, fmt.Sprintf("Deposit:%d", userID), 1).AllowN(ctx, limitCount) { //限频
		handleError(c, CODE_REQUEST_TOO_QUICKLY, errors.New("trigger limit exceeded"))
		return
	}
	var request struct{ Amount decimal.Decimal }
	if err := c.BindJSON(&request); err != nil {
		wc.logger.Error(ctx, "WalletController Deposit BindJSON",
			zap.Int("userID", userID), zap.Error(err))
		handleError(c, CODE_INVALID_PARAMS, err)
		return
	}
	if err := wc.walletService.Deposit(ctx, userID, userID, request.Amount, models.DepositTransactionType); err != nil {
		wc.logger.Error(ctx, "WalletController Deposit walletService",
			zap.Int("userID", userID), zap.Error(err))
		handleError(c, CODE_INTERNALSERVER, err)
		return
	}
	handleSuccess(c, gin.H{"status": "Deposit successful"})
}

// Withdraw 取款
func (wc *WalletController) Withdraw(c *gin.Context) {
	wc.logger.WithField("traceID", generateTraceID())
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		handleError(c, CODE_INVALID_PARAMS, err)
		return
	}
	ctx := c.Request.Context()
	if !rdsLimit.NewRdsLimit(wc.redis, fmt.Sprintf("Withdraw:%d", userID), 1).AllowN(ctx, limitCount) { //限频
		handleError(c, CODE_REQUEST_TOO_QUICKLY, errors.New("trigger limit exceeded"))
		return
	}

	var request struct{ Amount decimal.Decimal }
	if err := c.BindJSON(&request); err != nil {
		wc.logger.Error(ctx, "WalletController Deposit BindJSON",
			zap.Int("userID", userID), zap.Error(err))
		handleError(c, CODE_INVALID_PARAMS, err)
		return
	}
	if err := wc.walletService.Withdraw(ctx, userID, userID, request.Amount, models.WithdrawTransactionType); err != nil {
		wc.logger.Error(ctx, "WalletController Deposit BindJSON",
			zap.Int("userID", userID), zap.Error(err))
		handleError(c, CODE_INTERNALSERVER, err)
		return
	}
	handleSuccess(c, gin.H{"status": "Withdraw successful"})
}

// Transfer  转账
func (wc *WalletController) Transfer(c *gin.Context) {
	wc.logger.WithField("traceID", generateTraceID())
	senderID, err := strconv.Atoi(c.Param("sender_id"))
	if err != nil {
		handleError(c, CODE_INVALID_PARAMS, err)
		return
	}
	ctx := c.Request.Context()
	if !rdsLimit.NewRdsLimit(wc.redis, fmt.Sprintf("Transfer:%d", senderID), 1).AllowN(ctx, limitCount) { //限频
		handleError(c, CODE_REQUEST_TOO_QUICKLY, errors.New("trigger limit exceeded"))
		return
	}
	receiverID, err := strconv.Atoi(c.Param("receiver_id"))
	if err != nil {
		wc.logger.Error(ctx, "WalletController Transfer strconv.Atoi",
			zap.Int("senderID", senderID), zap.Int("receiverID", receiverID), zap.Error(err))
		handleError(c, CODE_INVALID_PARAMS, err)
		return
	}

	var request struct{ Amount decimal.Decimal }
	if err := c.BindJSON(&request); err != nil {
		wc.logger.Error(ctx, "WalletController Transfer BindJSON",
			zap.Int("senderID", senderID), zap.Int("receiverID", receiverID), zap.Error(err))
		handleError(c, CODE_INVALID_PARAMS, err)
		return
	}
	if err := wc.walletService.Transfer(ctx, senderID, receiverID, request.Amount); err != nil {
		wc.logger.Error(ctx, "WalletController Transfer walletService ",
			zap.Int("senderID", senderID), zap.Int("receiverID", receiverID), zap.Error(err))
		handleError(c, CODE_INTERNALSERVER, err)
		return
	}
	handleSuccess(c, gin.H{"status": "Transfer successful"})
}

// GetBalance 查询余额
func (wc *WalletController) GetBalance(c *gin.Context) {
	wc.logger.WithField("traceID", generateTraceID())
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		handleError(c, CODE_INVALID_PARAMS, err)
		return
	}

	ctx := c.Request.Context()
	if !rdsLimit.NewRdsLimit(wc.redis, fmt.Sprintf("GetBalance:%d", userID), 1).AllowN(ctx, limitCount) { //限频
		handleError(c, CODE_REQUEST_TOO_QUICKLY, errors.New("trigger limit exceeded"))
		return
	}
	balance, err := wc.walletService.GetBalance(ctx, userID)
	if err != nil {
		wc.logger.Error(ctx, "WalletController GetBalance",
			zap.Int("userID", userID), zap.Error(err))
		handleError(c, CODE_INTERNALSERVER, err)
		return
	}
	handleSuccess(c, gin.H{"balance": balance})
}

// GetTransactionHistory 获取交易历史
func (wc *WalletController) GetTransactionHistory(c *gin.Context) {
	wc.logger.WithField("traceID", generateTraceID())
	ctx := c.Request.Context()
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		handleError(c, CODE_INVALID_PARAMS, err)
		return
	}

	if !rdsLimit.NewRdsLimit(wc.redis, fmt.Sprintf("GetTransactionHistory:%d", userID), 1).AllowN(ctx, limitCount) { //限频
		handleError(c, CODE_REQUEST_TOO_QUICKLY, errors.New("trigger limit exceeded"))
		return
	}

	transactions, err := wc.walletService.GetTransactionHistory(userID)
	if err != nil {
		wc.logger.Error(ctx, "WalletController GetTransactionHistory",
			zap.Int("userID", userID), zap.Error(err))
		handleError(c, CODE_INTERNALSERVER, err)
		return
	}
	if len(transactions) == 0 {
		handleError(c, CODE_NOT_FOUND, errors.New("No transaction history"))
		return
	}
	handleSuccess(c, transactions)
}
