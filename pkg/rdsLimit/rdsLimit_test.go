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
		Convey("Should initialize correctly", func() {
			So(r.countPerSeconds, ShouldEqual, 3)
			So(r.key, ShouldEqual, "test")
			So(r.limiter, ShouldNotBeNil)
			So(r.logger, ShouldNotBeNil)
		})
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
		Convey("Should allow requests within the limit", func() {
			So(r.Allow(ctx), ShouldEqual, true)
			So(r.Allow(ctx), ShouldEqual, true)
			So(r.Allow(ctx), ShouldEqual, true)
		})
		Convey("Should deny requests exceeding the limit", func() {
			So(r.Allow(ctx), ShouldEqual, true)
			So(r.Allow(ctx), ShouldEqual, true)
		})
		Convey("Should reset the limit after a period", func() {
			time.Sleep(2 * time.Second) // 等待 2 秒，使计数器重置
			So(r.Allow(ctx), ShouldEqual, true)
		})
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
		Convey("Should allow batch requests within the limit", func() {
			So(r.AllowN(ctx, 1), ShouldEqual, true)
		})
		Convey("Should deny batch requests exceeding the limit", func() {
			So(r.AllowN(ctx, 2), ShouldEqual, true)
			So(r.AllowN(ctx, 3), ShouldEqual, true)
			So(r.AllowN(ctx, 10000), ShouldEqual, false)
		})
		Convey("Should reset the limit after a period", func() {
			time.Sleep(2 * time.Second) // 等待 2 秒，使计数器重置
			So(r.AllowN(ctx, 1), ShouldEqual, true)
		})
	})
}

func TestAllowWithError(t *testing.T) {
	Convey("TestAllowWithError", t, func() {
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
		Convey("Should handle errors gracefully", func() {
			// Simulate an error by using an invalid context
			So(r.Allow(ctx), ShouldEqual, true)
		})
	})
}

func TestAllowNWithError(t *testing.T) {
	Convey("TestAllowNWithError", t, func() {
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
		Convey("Should handle errors gracefully", func() {
			// Simulate an error by using an invalid context
			So(r.AllowN(ctx, 1), ShouldEqual, true)
		})
	})
}

func TestAllow1(t *testing.T) {
	Convey("TestAllow1", t, func() {
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
