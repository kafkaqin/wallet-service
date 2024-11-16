package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"time"
	"wallet-service/models"
	wallet_logger "wallet-service/pkg/logger"
)

type WalletService interface {
	Deposit(ctx context.Context, senderID, receiverID int, amount decimal.Decimal, transactionType models.TransactionType) error
	Withdraw(ctx context.Context, senderID, receiverID int, amount decimal.Decimal, transactionType models.TransactionType) error
	Transfer(ctx context.Context, senderID, receiverID int, amount decimal.Decimal) error
	GetBalance(ctx context.Context, userID int) (decimal.Decimal, error)
	GetTransactionHistory(userID int) ([]models.Transaction, error)
}

type walletService struct {
	db     *sqlx.DB
	redis  *redis.Client
	logger *wallet_logger.Logger
}

var _ WalletService = &walletService{}

// NewWalletService service
func NewWalletService(logger *wallet_logger.Logger, db *sqlx.DB, redis *redis.Client) WalletService {
	return &walletService{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

// Deposit 存款
func (s *walletService) Deposit(ctx context.Context, senderID, receiverID int, amount decimal.Decimal, transactionType models.TransactionType) error {
	if amount.LessThan(decimal.Zero) {
		return errors.New("amount must be greater than zero")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		s.logger.Error(ctx, "Deposit Failed to begin transaction", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				s.logger.Error(ctx, "Deposit Failed Rollback transaction", zap.Int("senderID", senderID),
					zap.Int("receiverID", receiverID), zap.Error(err))
			}
		}
	}()

	res, err := tx.Exec("UPDATE wallets SET balance = balance + $1 WHERE user_id = $2", amount, senderID)
	if err != nil {
		s.logger.Error(ctx, "Deposit Failed to Exec transaction", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		query := `
			INSERT INTO wallets (user_id, balance)
			VALUES ($1, $2)
			ON CONFLICT (user_id)
			DO UPDATE SET balance = wallets.balance + EXCLUDED.balance
	   `
		_, err = tx.Exec(query, senderID, amount)
		if err != nil {
			s.logger.Error(ctx, "Deposit Failed to Exec transaction", zap.Int("senderID", senderID),
				zap.Int("receiverID", receiverID), zap.Error(err))
			return err
		}
	}

	if len(transactionType) == 0 {
		transactionType = models.DepositTransactionType
	}

	_, err = tx.Exec("INSERT INTO transactions (sender_user_id, receiver_user_id, transaction_type, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
		senderID, receiverID, transactionType, amount, time.Now())
	if err != nil {
		s.logger.Error(ctx, "Deposit Failed insert into  transactions ", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err
	}

	// 更新缓存中的余额

	err = s.redis.IncrByFloat(ctx, fmt.Sprintf("wallet:balance:%d", senderID), amount.InexactFloat64()).Err()
	if err != nil {
		// 记录日志，不影响事务
		s.logger.Warn(ctx, "Deposit Failed to update Redis cache:", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
	}

	return tx.Commit()
}

func (s *walletService) DepositWithTx(ctx context.Context, tx *sqlx.Tx, senderID, receiverID int, amount decimal.Decimal, transactionType models.TransactionType) error {

	if amount.LessThan(decimal.Zero) {
		return errors.New("amount must be greater than zero")
	}

	res, err := tx.Exec("UPDATE wallets SET balance = balance + $1 WHERE user_id = $2", amount, senderID)
	if err != nil {
		s.logger.Warn(ctx, "depositWithTx Failed to update Redis cache:", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		query := `
			INSERT INTO wallets (user_id, balance)
			VALUES ($1, $2)
			ON CONFLICT (user_id) 
			DO UPDATE SET balance = wallets.balance + EXCLUDED.balance
        `
		_, err = tx.Exec(query, senderID, amount)
		if err != nil {
			s.logger.Warn(ctx, "depositWithTx Failed to update Redis cache:", zap.Int("senderID", senderID),
				zap.Int("receiverID", receiverID), zap.Error(err))
			return err
		}
	}

	if len(transactionType) == 0 {
		transactionType = models.DepositTransactionType
	}

	if tx == nil {
		tx, err = s.db.Beginx()
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec("INSERT INTO transactions (receiver_user_id,sender_user_id, transaction_type, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
		senderID, receiverID, transactionType, amount, time.Now())
	if err != nil {
		s.logger.Warn(ctx, "depositWithTx Failed to update Redis cache:", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err
	}

	// 更新缓存中的余额
	err = s.redis.IncrByFloat(ctx, fmt.Sprintf("wallet:balance:%d", senderID), amount.InexactFloat64()).Err()
	if err != nil {
		// 记录日志，不影响事务
		s.logger.Warn(ctx, "depositWithTx Failed to update Redis cache:", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
	}

	return nil

}

// Withdraw 取款
func (s *walletService) Withdraw(ctx context.Context, senderID, receiverID int, amount decimal.Decimal, transactionType models.TransactionType) error {

	if amount.LessThan(decimal.Zero) {
		return errors.New("amount must be greater than zero")
	}

	var balance decimal.Decimal
	// 从 Redis 缓存中查询余额
	cacheBalance, err := s.redis.Get(ctx, fmt.Sprintf("wallet:balance:%d", senderID)).Result()

	switch {
	case err == redis.Nil:
		// 缓存不存在，从数据库查询余额
		err = s.db.Get(&balance, "SELECT balance FROM wallets WHERE user_id = $1", senderID)
		if err != nil {
			s.logger.Error(ctx, "Withdraw Failed to select from wallets pg", zap.Int("senderID", senderID),
				zap.Int("receiverID", receiverID), zap.Error(err))
			return errors.New("insufficient funds")
		}

	case err != nil:
		// Redis 获取余额失败
		s.logger.Error(ctx, "Withdraw Failed to get Balance from redis", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err

	default:
		// 从缓存中读取余额
		balance, err = decimal.NewFromString(cacheBalance)
		if err != nil {
			s.logger.Error(ctx, "Withdraw Failed decimal.NewFromString", zap.Int("senderID", senderID),
				zap.Int("receiverID", receiverID), zap.Error(err))
			return err
		}
		// 如果缓存的余额为零，重新从数据库查询
		if balance.IsZero() {
			err = s.db.Get(&balance, "SELECT balance FROM wallets WHERE user_id = $1", senderID)
			if err != nil {
				s.logger.Error(ctx, "Withdraw Failed to select from wallets pg", zap.Int("senderID", senderID),
					zap.Int("receiverID", receiverID), zap.Error(err))
				return errors.New("insufficient funds")
			}
		}
	}

	// 检查余额是否足够
	if balance.LessThan(amount) {
		return errors.New("insufficient funds balance")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				s.logger.Error(ctx, "Withdraw Failed Rollback transaction", zap.Int("senderID", senderID),
					zap.Int("receiverID", receiverID), zap.Error(err))
			}
		}
	}()

	_, err = tx.Exec("UPDATE wallets SET balance = balance - $1 WHERE user_id = $2", amount, senderID)
	if err != nil {
		s.logger.Error(ctx, "Withdraw Failed to Exec transaction update wallets", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err
	}

	if len(transactionType) == 0 {
		transactionType = models.WithdrawTransactionType
	}

	_, err = tx.Exec("INSERT INTO transactions (sender_user_id, receiver_user_id, transaction_type, amount, created_at) VALUES ($1, $2, $3, $4, $5)",
		senderID, receiverID, transactionType, amount, time.Now())
	if err != nil {
		s.logger.Error(ctx, "Withdraw Failed to Exec transaction: insert into pg transactions", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err
	}

	// 更新 Redis 缓存
	err = s.redis.Set(ctx, fmt.Sprintf("wallet:balance:%d", senderID), balance.Sub(amount).String(), 0).Err()
	if err != nil {
		// 记录日志，不影响事务
		s.logger.Warn(ctx, "Withdraw Failed to update Redis cache:", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
	}

	return tx.Commit()
}

func (s *walletService) WithdrawWithTx(ctx context.Context, tx *sqlx.Tx, senderID int, amount decimal.Decimal) error {

	if amount.LessThan(decimal.Zero) {
		return errors.New("amount must be greater than zero")
	}

	var balance decimal.Decimal
	// 从 Redis 缓存中查询余额
	cacheBalance, err := s.redis.Get(ctx, fmt.Sprintf("wallet:balance:%d", senderID)).Result()
	switch {
	case err == redis.Nil:
		// 缓存不存在，从数据库查询余额
		err = s.db.Get(&balance, "SELECT balance FROM wallets WHERE user_id = $1", senderID)
		if err != nil {
			s.logger.Error(ctx, "Withdraw Failed get balance from pg", zap.Int("senderID", senderID),
				zap.Error(err))
			return errors.New("insufficient funds")
		}

	case err != nil:
		// 如果获取 Redis 时出现错误，返回错误
		return err

	default:
		// 从缓存中读取余额
		balance, err = decimal.NewFromString(cacheBalance)
		if err != nil {
			return err
		}
		if balance.IsZero() {
			err = s.db.Get(&balance, "SELECT balance FROM wallets WHERE user_id = $1", senderID)
			if err != nil {
				s.logger.Error(ctx, "Withdraw Failed get balance from pg ", zap.Int("senderID", senderID),
					zap.Error(err))
				return errors.New("insufficient funds")
			}
		}
	}

	if tx == nil {
		tx, err = s.db.Beginx()
		if err != nil {
			return err
		}
	}

	// 检查余额是否足够
	if balance.LessThan(amount) {
		return errors.New("insufficient funds balance")
	}

	_, err = tx.Exec("UPDATE wallets SET balance = balance - $1 WHERE user_id = $2", amount, senderID)
	if err != nil {
		s.logger.Error(ctx, "Withdraw Failed to Exec transaction: UPDATE wallets ", zap.Int("senderID", senderID),
			zap.Error(err))
		return err
	}

	// 更新 Redis 缓存
	err = s.redis.Set(ctx, fmt.Sprintf("wallet:balance:%d", senderID), balance.Sub(amount).String(), 0).Err()
	if err != nil {
		// 记录日志，不影响事务
		s.logger.Warn(ctx, "Withdraw Failed to update Redis cache:", zap.Int("senderID", senderID),
			zap.Error(err))
	}

	return nil
}

// Transfer 转账
func (s *walletService) Transfer(ctx context.Context, senderID, receiverID int, amount decimal.Decimal) error {

	if amount.LessThan(decimal.Zero) {
		return errors.New("amount must be greater than zero")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		s.logger.Error(ctx, "Transfer Failed to begin transaction:", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err
	}
	// 确保事务回滚
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				s.logger.Error(ctx, "Transfer Failed Rollback transaction", zap.Int("senderID", senderID),
					zap.Int("receiverID", receiverID), zap.Error(err))
			}
		}
	}()

	err = s.WithdrawWithTx(ctx, tx, senderID, amount)
	if err != nil {
		s.logger.Error(ctx, "Transfer Failed withdrawWithTx:", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err
	}

	err = s.DepositWithTx(ctx, tx, receiverID, senderID, amount, models.TransferTransactionType)
	if err != nil {
		s.logger.Error(ctx, "Transfer Failed depositWithTx:", zap.Int("senderID", senderID),
			zap.Int("receiverID", receiverID), zap.Error(err))
		return err
	}

	return tx.Commit()
}

// GetBalance 查询余额
func (s *walletService) GetBalance(ctx context.Context, userID int) (decimal.Decimal, error) {
	var balance decimal.Decimal
	// 尝试从 Redis 获取缓存中的余额
	cacheBalance, err := s.redis.Get(ctx, fmt.Sprintf("wallet:balance:%d", userID)).Result()

	switch {
	case err == redis.Nil:
		// 缓存不存在，从数据库查询余额
		err = s.db.Get(&balance, "SELECT balance FROM wallets WHERE user_id = $1", userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return balance, errors.New("wallet not found")
			}
			s.logger.Error(ctx, "GetBalance Failed get balance from pg:", zap.Int("userID", userID),
				zap.Error(err))
			return balance, err
		}
		// 查询成功后，将数据缓存到 Redis
		err = s.redis.Set(ctx, fmt.Sprintf("wallet:balance:%d", userID), balance.String(), 0).Err()
		if err != nil {
			// 记录日志，不影响主流程
			s.logger.Warn(ctx, "GetBalance Failed to cache balance:", zap.Error(err))
		}

	case err != nil:
		// Redis 查询失败，记录日志并返回错误
		s.logger.Error(ctx, "GetBalance Failed get balance from cache:", zap.Int("userID", userID),
			zap.Error(err))
		return balance, err

	default:
		// 从缓存中读取余额
		balance, err = decimal.NewFromString(cacheBalance)
		if err != nil {
			s.logger.Error(ctx, "GetBalance Failed to decimal.NewFromString from redis ", zap.Int("userID", userID),
				zap.Error(err))
			return balance, err
		}
	}

	return balance, nil
}

// GetTransactionHistory 获取交易历史
func (s *walletService) GetTransactionHistory(userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := s.db.Select(&transactions, "SELECT * FROM transactions WHERE sender_user_id = $1", userID)
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
