package logger

import (
	"sync"
	"time"
)

var onceLocation = sync.Once{}

var loc *time.Location

func Get() *time.Location {
	onceLocation.Do(func() {
		l, err := time.LoadLocation("Asia/Shanghai") //设置时区
		if err != nil {
			l = time.FixedZone("CST", 8*3600)
		}
		loc = l
	})
	return loc
}
