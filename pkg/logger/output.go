package logger

import (
	"errors"
	"os"
)

const (
	keyOutputPath = "LOG_OUTPUT_PATH"
)

func getOutputPaths() (string, error) {
	var p = os.Getenv(keyOutputPath)
	if len(p) == 0 {
		return "", errors.New("emtpy path")
	}
	return p, nil
}
