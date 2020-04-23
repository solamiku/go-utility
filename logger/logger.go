package logger

import (
	"fmt"
	"log"

	"github.com/solamiku/go-utility/runtime"
)

type Logger interface {
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

var DefaultLogger Logger

func init() {
	DefaultLogger = &defaultLogger{}
}

type defaultLogger struct{}

func (dl *defaultLogger) Infof(str string, args ...interface{}) {
	info := fmt.Sprintf(str, args...)
	log.Println("[INFO] " + runtime.CallInfo(2) + " " + info)
}

func (dl *defaultLogger) Errorf(str string, args ...interface{}) {
	err := fmt.Sprintf(str, args...)
	log.Println("[ERROR] " + runtime.CallInfo(2) + " " + err)
}
