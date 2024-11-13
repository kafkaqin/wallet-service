package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type TransactionType string

const (
	DepositTransactionType  TransactionType = "deposit"
	WithdrawTransactionType TransactionType = "withdraw"
	TransferTransactionType TransactionType = "transfer"
)

type Transaction struct {
	ID              int             `db:"id"`
	SenderUserID    int             `db:"sender_user_id"`
	ReceiverUserID  int             `db:"receiver_user_id"`
	TransactionType TransactionType `db:"transaction_type"` // "deposit", "withdraw", "transfer"
	Amount          decimal.Decimal `db:"amount"`           // 使用 decimal.Decimal 处理金额
	CreatedAt       time.Time       `db:"created_at"`
}
