package utility

import (
	"fmt"
	"runtime"
)

// return a error with format string
// auto create prefix - [file:line]
func Errorf(format string, a ...interface{}) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf(format, a...)
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	s := fmt.Sprintf(format, a...)
	return fmt.Errorf("[%s:%d]: %s", file, line, s)
}
