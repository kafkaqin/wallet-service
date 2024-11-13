package logger

import "os"

const (
	keyLogLevel = "LOG_LOGLEVEL"
)

func setLogLevelFromEnviron() {
	var levelStr = os.Getenv(keyLogLevel)
	err := ParseAndSetLogLevel(levelStr)
	if err != nil {
		SetLogLevel(DebugLevel) // 默认设置为DEBUG等级
	} else {
		//do nothing
	}
}
