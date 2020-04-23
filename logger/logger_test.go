package logger

import "testing"

func Test_logger(t *testing.T) {
	DefaultLogger.Infof("test info output")
	DefaultLogger.Errorf("test err output")
}
