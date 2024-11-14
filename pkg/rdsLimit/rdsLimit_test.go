package rdsLimit

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
	"time"
)

func TestNewRLim(t *testing.T) {
	Convey("TestNewRLim", t, func() {
		ctx := context.Background()
		var err error
		var redisCli = redis.NewClient(&redis.Options{
			PoolSize:     100, // 最大链接数
			MinIdleConns: 25,  // 空闲链接数
			Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			Password:     "",
			DB:           0,
		})
		_, err = redisCli.Ping(ctx).Result()
		if err != nil {
			log.Fatal(err)
		}
		redisCli.FlushDB(ctx) // 清空redis记录，方便测试

		var r = NewRdsLimit(redisCli, "test", 3)
		So(r.Allow(ctx), ShouldEqual, true)
		So(r.Allow(ctx), ShouldEqual, true)
		So(r.Allow(ctx), ShouldEqual, true)
		So(r.Allow(ctx), ShouldEqual, false)
		So(r.Allow(ctx), ShouldEqual, false)
		So(r.AllowN(ctx, 1), ShouldEqual, false)
	})
}

func TestNewRLimX(t *testing.T) {
	Convey("TestNewRLimX", t, func() {
		ctx := context.Background()
		var err error
		var redisCli = redis.NewClient(&redis.Options{
			PoolSize:     100, // 最大链接数
			MinIdleConns: 25,  // 空闲链接数
			Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			Password:     "",
			DB:           0,
		})
		_, err = redisCli.Ping(ctx).Result()
		if err != nil {
			log.Fatal(err)
		}
		redisCli.FlushDB(ctx) // 清空redis记录，方便测试

		var r = NewRdsLimit(redisCli, "test1", 1)
		So(r.AllowN(ctx, 1), ShouldEqual, true)
		So(r.AllowN(ctx, 2), ShouldEqual, true)
		So(r.AllowN(ctx, 3), ShouldEqual, true)
		So(r.AllowN(ctx, 10000), ShouldEqual, false)
	})
}

func TestNewRdsLimit(t *testing.T) {
	Convey("TestNewRdsLimit", t, func() {
		ctx := context.Background()
		var err error
		var redisCli = redis.NewClient(&redis.Options{
			PoolSize:     100, // 最大链接数
			MinIdleConns: 25,  // 空闲链接数
			Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			Password:     "",
			DB:           0,
		})
		_, err = redisCli.Ping(ctx).Result()
		if err != nil {
			log.Fatal(err)
		}
		redisCli.FlushDB(ctx) // 清空redis记录，方便测试

		var r = NewRdsLimit(redisCli, "test", 3)
		So(r.countPerSeconds, ShouldEqual, 3)
		So(r.key, ShouldEqual, "test")
		So(r.limiter, ShouldNotBeNil)
		So(r.logger, ShouldNotBeNil)
	})
}

func TestAllow(t *testing.T) {
	Convey("TestAllow", t, func() {
		ctx := context.Background()
		var err error
		var redisCli = redis.NewClient(&redis.Options{
			PoolSize:     100, // 最大链接数
			MinIdleConns: 25,  // 空闲链接数
			Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			Password:     "",
			DB:           0,
		})
		_, err = redisCli.Ping(ctx).Result()
		if err != nil {
			log.Fatal(err)
		}
		redisCli.FlushDB(ctx) // 清空redis记录，方便测试

		var r = NewRdsLimit(redisCli, "test", 3)
		So(r.Allow(ctx), ShouldEqual, true)
		So(r.Allow(ctx), ShouldEqual, true)
		So(r.Allow(ctx), ShouldEqual, true)
		So(r.Allow(ctx), ShouldEqual, false)
		So(r.Allow(ctx), ShouldEqual, false)
	})
}

func TestAllowN(t *testing.T) {
	Convey("TestAllowN", t, func() {
		ctx := context.Background()
		var err error
		var redisCli = redis.NewClient(&redis.Options{
			PoolSize:     100, // 最大链接数
			MinIdleConns: 25,  // 空闲链接数
			Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			Password:     "",
			DB:           0,
		})
		_, err = redisCli.Ping(ctx).Result()
		if err != nil {
			log.Fatal(err)
		}
		redisCli.FlushDB(ctx) // 清空redis记录，方便测试

		var r = NewRdsLimit(redisCli, "test1", 1)
		So(r.AllowN(ctx, 1), ShouldEqual, true)
		So(r.AllowN(ctx, 2), ShouldEqual, true)
		So(r.AllowN(ctx, 3), ShouldEqual, true)
		So(r.AllowN(ctx, 10000), ShouldEqual, false)

		// 测试 AllowN 在不同时间间隔内的行为
		time.Sleep(2 * time.Second) // 等待 2 秒，使计数器重置
		So(r.AllowN(ctx, 1), ShouldEqual, true)
	})
}
