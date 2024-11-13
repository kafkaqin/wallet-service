package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Wallet struct {
	UserID    int             `db:"user_id" json:"user_id"`
	Balance   decimal.Decimal `db:"balance" json:"balance"` // 使用 decimal.Decimal 处理金额
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}
