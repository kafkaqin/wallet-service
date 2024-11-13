package rdsLimit

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type RdsLimit struct {
	countPerSeconds int
	key             string
	limiter         *redis_rate.Limiter
	logger          *zap.Logger
}

func NewRdsLimit(rdb *redis.Client, key string, countPerSeconds int) *RdsLimit {
	var r = &RdsLimit{
		countPerSeconds: countPerSeconds,
		key:             key,
		limiter:         redis_rate.NewLimiter(rdb),
		logger: zap.L().With(
			zap.String("module", "rds_limit"),
			zap.Int("countPerSeconds", countPerSeconds),
			zap.String("key", key),
		),
	}
	return r
}

func (r *RdsLimit) Allow(ctx context.Context) bool {
	result, err := r.limiter.Allow(ctx, r.key, redis_rate.PerSecond(r.countPerSeconds))
	if err != nil {
		r.logger.Info("call allow failed", zap.Error(err))
		return false
	}
	return result.Allowed == 1
}

func (r *RdsLimit) AllowN(ctx context.Context, n int) bool {
	limit := redis_rate.Limit{
		Burst:  100,                              // 突发量：允许每秒最多 100 次请求
		Rate:   100,                              // 速率：每秒 100 次请求
		Period: time.Duration(r.countPerSeconds), // 周期：1 秒
	}
	result, err := r.limiter.AllowN(ctx, r.key, limit, n)
	if err != nil {
		r.logger.Info("call allow failed", zap.Error(err))
		return false
	}
	return result.Allowed == n
}
