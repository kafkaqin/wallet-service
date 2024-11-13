package logger

import "sync"

var once = sync.Once{}

var logger *Logger

func instance() *Logger {
	once.Do(func() {
		logger = NewLogger()
	})
	return logger
}
