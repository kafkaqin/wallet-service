package rdsLimit

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestNewRLim(t *testing.T) {
	Convey("TestNewRLim", t, func() {
		ctx := context.Background()
		var err error
		var redisCli = redis.NewClient(&redis.Options{
			PoolSize:     100, //最大链接数
			MinIdleConns: 25,  //空闲链接数
			Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
			Password:     "",
			DB:           0,
		})
		_, err = redisCli.Ping(ctx).Result()
		if err != nil {
			log.Fatal(err)
		}
		redisCli.FlushDB(ctx) //清空redis记录，方便测试

		var r = NewRdsLimit(redisCli, "test", 3)
		So(r.Allow(ctx), ShouldEqual, true)
		So(r.Allow(ctx), ShouldEqual, true)
		So(r.Allow(ctx), ShouldEqual, true)
		So(r.Allow(ctx), ShouldEqual, false)
		So(r.Allow(ctx), ShouldEqual, false)
		So(r.AllowN(ctx, 1), ShouldEqual, false)
	})
}
