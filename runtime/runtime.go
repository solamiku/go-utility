package runtime

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
)

// return heap objs num in memory
func LookupHeapObjs() int64 {
	p := pprof.Lookup("heap")
	buff := bytes.NewBuffer(make([]byte, 0))
	p.WriteTo(buff, 2)
	rx := regexp.MustCompile(`#\s*(HeapObjects)\s*=\s*(\d+)`)
	rd := bufio.NewReader(buff)
	for line, err := rd.ReadString('\n'); err == nil; line, err = rd.ReadString('\n') {
		l := len(line)
		if l > 0 {
			line = line[:l-1]
		}
		match := rx.FindStringSubmatch(line)
		if len(match) > 0 {
			i, _ := strconv.Atoi(match[2])
			return int64(i)
		}
	}
	return 0
}

// get short file name at now caller
func callerShortfile(file string, lastsep_ ...rune) string {
	lastsep := '/'
	if len(lastsep_) > 0 {
		lastsep = lastsep_[0]
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == byte(lastsep) {
			short = file[i+1:]
			break
		}
	}
	return short
}

// add prefix to error
func Error(err error) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return err
	}
	file = callerShortfile(file)
	return fmt.Errorf("[%s:%d]: %v", file, line, err)
}

// create error auto add prefix
func Errof(str string, args ...interface{}) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf(str, args...)
	}
	file = callerShortfile(file)
	s := fmt.Sprintf(str, args...)
	return fmt.Errorf("[%s:%d]:%s", file, line, s)
}

// return panic call stack info
func GetPanicRoutineCallstack() string {
	s := ""
	if e := recover(); e != nil {
		s += fmt.Sprintf("%s\n", e)
		lv := 3
		for {
			_, file, line, ok := runtime.Caller(lv)
			if !ok {
				break
			}
			file = callerShortfile(file)
			s += fmt.Sprintf("%d) %s: %d\n", lv, file, line)
			lv++
		}
	}
	return s
}

// write call stack to writer
func WriteRoutineCallstack(lv int, wr io.Writer) {
	for {
		_, file, line, ok := runtime.Caller(lv)
		if !ok {
			break
		}
		file = callerShortfile(file)
		s := fmt.Sprintf("%d) %s: %d\n", lv, file, line)
		wr.Write([]byte(s))
		lv++
	}
}
func WriteRoutineCallstackFull(lv int, wr io.Writer) {
	for {
		s := CallInfo(lv)
		if len(s) == 0 {
			break
		}
		s1 := fmt.Sprintf("%d) %s\n", lv, s)
		wr.Write([]byte(s1))
		lv++
	}
}

// return call info, format - filename:line(funcname)
// lv=0 is CallInfo  lv=1 is which call CallInfo
func CallInfo(lv int) string {
	pc, file, line, ok := runtime.Caller(lv)
	if !ok {
		return ""
	}
	file = callerShortfile(file)
	funcName := runtime.FuncForPC(pc).Name()
	funcName = callerShortfile(funcName)
	fn := callerShortfile(funcName, ')')
	if len(fn) < len(funcName) {
		if len(fn) > 1 && fn[0] == '.' {
			fn = fn[1:]
		}
		funcName = fn
	} else {
		funcName = callerShortfile(funcName, '.')
	}
	s := fmt.Sprintf("%s:%d(%s)", file, line, funcName)
	return s
}

// return call info, format- filename(funcname)
// lv=0 is CallInfo  lv=1 is which call CallInfo
func SimpleCallInfo(lv int) string {
	pc, file, _, ok := runtime.Caller(lv)
	if !ok {
		return ""
	}
	file = callerShortfile(file)
	funcName := runtime.FuncForPC(pc).Name()
	funcName = callerShortfile(funcName)
	fn := callerShortfile(funcName, ')')
	if len(fn) < len(funcName) {
		if len(fn) > 1 && fn[0] == '.' {
			fn = fn[1:]
		}
		funcName = fn
	} else {
		funcName = callerShortfile(funcName, '.')
	}
	s := fmt.Sprintf("%s(%s)", file, funcName)
	return s
}

//create path include current directory
//@param start without "/"
func PathIncludeCurrentDir(path string) string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir + "/" + path
}

func IsPahtExisted(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func MakeDir(path string, mode ...os.FileMode) error {
	defaultMode := os.ModePerm
	if len(mode) > 0 {
		defaultMode = 0
		for _, m := range mode {
			defaultMode = defaultMode | m
		}
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return os.Chmod(path, defaultMode)
}

func ReadFileWithContentType(file string, types []string) (string, error) {
	var dat []byte
	var content string
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		content = err.Error()
	} else {
		content = string(dat)
	}
	if len(dat) > 512 {
		dat = dat[:512]
	}
	contentType := http.DetectContentType(dat)
	if len(types) > 0 {
		has := false
		for _, st := range types {
			if strings.Contains(contentType, st) {
				has = true
			}
		}
		if !has {
			content = fmt.Sprintf("file type is:%v, cannot show the content.", contentType)
			return content, Errof("type %s not matched.", contentType)
		}
	}
	return content, nil
}

// 带完整调用参数信息的panic栈信息
func WritePanicStackFull(prefix string, wr io.Writer) {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")

	byteSize := 1 << 11 //2KB
	stack := make([]byte, byteSize)
	stackWriteLen := runtime.Stack(stack, true)
	if stackWriteLen >= byteSize {
		stackWriteLen = byteSize - 1
	}
	// fmt.Printf("ori:\n%s\n", string(stack))
	// 去除开头和换行
	start := bytes.Index(stack, s)
	if start != -1 {
		stack = stack[start:stackWriteLen]
	}
	start = bytes.Index(stack, line) + 1
	if start != -1 {
		stack = stack[start:]
	}
	// 去除换行和结尾
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	// 组装
	stackString := string(stack)
	ss := strings.Split(stackString, "\n")
	filereg, _ := regexp.Compile(`[^/|^\\]+\.go`)
	funcreg, _ := regexp.Compile(`[^/|^\\]+\)`)
	for i := 0; i+1 < len(ss); i += 2 {
		filestr := filereg.FindString(ss[i+1])
		funcstr := funcreg.FindString(ss[i])
		fmt.Fprintf(wr, "%v(%v) %v(%v)\n", prefix, i/2+1, filestr, funcstr)
	}
}
