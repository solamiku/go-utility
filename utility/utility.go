package utility

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

var _DEFAULT_TIME = time.Unix(0, 0)

//parse interface to string
func Str(v interface{}) string {
	switch a := v.(type) {
	case string:
		return a
	case []byte:
		return string(a)
	case json.Number:
		return string(a)
	case float32:
		return strconv.FormatFloat(float64(a), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(a, 'f', -1, 64)
	case int:
		return strconv.FormatInt(int64(a), 10)
	case int8:
		return strconv.FormatInt(int64(a), 10)
	case int16:
		return strconv.FormatInt(int64(a), 10)
	case int32:
		return strconv.FormatInt(int64(a), 10)
	case int64:
		return strconv.FormatInt(a, 10)
	case uint:
		return strconv.FormatUint(uint64(a), 10)
	case uint8:
		return strconv.FormatUint(uint64(a), 10)
	case uint16:
		return strconv.FormatUint(uint64(a), 10)
	case uint32:
		return strconv.FormatUint(uint64(a), 10)
	case uint64:
		return strconv.FormatUint(a, 10)
	case bool:
		return strconv.FormatBool(a)
	case complex64, complex128:
		return fmt.Sprintf("%v", a)
	case interface{}:
		//which not match above will enter this scope
		refval := reflect.ValueOf(a)
		if refval.Kind() == reflect.Ptr {
			refval = reflect.Indirect(refval)
		}
		// check nil first, else Kind() will crash
		if refval.Interface() == nil {
			return ""
		}
		vv := reflect.ValueOf(refval.Interface())
		switch vv.Type().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Str(vv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Str(vv.Uint())
		case reflect.Float32, reflect.Float64:
			return Str(vv.Float())
		case reflect.String:
			return vv.String()
		case reflect.Array, reflect.Slice:
			switch vv.Type().Elem().Kind() {
			case reflect.Uint8:
				data := refval.Interface().([]byte)
				return string(data)
			}
		case reflect.Bool:
			return Str(vv.Bool())
		case reflect.Complex128, reflect.Complex64:
			return Str(vv.Complex())
		case reflect.Struct:
			//time.Time
			if vv.Type().ConvertibleTo(reflect.TypeOf(_DEFAULT_TIME)) {
				return refval.Convert(reflect.TypeOf(_DEFAULT_TIME)).Interface().(time.Time).Format(TIME_FMT_STR)
			}
			//...
		}
	}
	return ""
}

//
func Int(v interface{}) int {
	return int(Int64(v))
}

func str2int(s string) int64 {
	r, err := strconv.ParseInt(s, 10, 0)
	if err == nil {
		return int64(r)
	} else if strings.Index(s, ".") >= 0 {
		r, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return int64(r)
		}
	}
	return 0
}

func Int64(v interface{}) int64 {
	switch a := v.(type) {
	case string:
		return str2int(a)
	case []byte:
		s := string(a)
		return str2int(s)
	case json.Number:
		s := string(a)
		return str2int(s)
	case float32:
		return int64(a)
	case float64:
		return int64(a)
	case int:
		return int64(a)
	case int8:
		return int64(a)
	case int16:
		return int64(a)
	case int32:
		return int64(a)
	case int64:
		return a
	case uint:
		return int64(a)
	case uint8:
		return int64(a)
	case uint16:
		return int64(a)
	case uint32:
		return int64(a)
	case uint64:
		return int64(a)
	case bool:
		if a {
			return 1
		}
		return 0
	case complex64, complex128:
		return 0
	case interface{}:
		//which not match above will enter this scope
		refval := reflect.ValueOf(a)
		if refval.Kind() == reflect.Ptr {
			refval = reflect.Indirect(refval)
		}
		// check nil first, else Kind() will crash
		if refval.Interface() == nil {
			return 0
		}
		vv := reflect.ValueOf(refval.Interface())
		switch vv.Type().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return vv.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Int64(vv.Uint())
		case reflect.Float32, reflect.Float64:
			return Int64(vv.Float())
		case reflect.String:
			return Int64(vv.String())
		case reflect.Array, reflect.Slice:
			switch vv.Type().Elem().Kind() {
			case reflect.Uint8:
				data := refval.Interface().([]byte)
				return Int64(string(data))
			}
		case reflect.Bool:
			return Int64(vv.Bool())
		case reflect.Complex128, reflect.Complex64:
			return Int64(vv.Complex())
		case reflect.Struct:
			if vv.Type().ConvertibleTo(reflect.TypeOf(_DEFAULT_TIME)) {
				return refval.Convert(reflect.TypeOf(_DEFAULT_TIME)).Interface().(time.Time).Unix()
			}
		}
	}
	return 0
}

func Float32(v interface{}) float32 {
	return float32(Float64(v))
}

func str2float(s string) float64 {
	r, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return r
	}
	return 0.0
}

func Float64(v interface{}) float64 {
	switch a := v.(type) {
	case string:
		return str2float(a)
	case []byte:
		s := string(a)
		return str2float(s)
	case json.Number:
		s := string(a)
		return str2float(s)
	case float32:
		return float64(a)
	case float64:
		return a
	case int:
		return float64(a)
	case int8:
		return float64(a)
	case int16:
		return float64(a)
	case int32:
		return float64(a)
	case int64:
		return float64(a)
	case uint:
		return float64(a)
	case uint8:
		return float64(a)
	case uint16:
		return float64(a)
	case uint32:
		return float64(a)
	case uint64:
		return float64(a)
	case bool:
		if a {
			return 1
		}
		return 0
	case complex64, complex128:
		return 0
	case *simplejson.Json:
		return Float64(a.Interface())
	case interface{}:
		//which not match above will enter this scope
		refval := reflect.ValueOf(a)
		if refval.Kind() == reflect.Ptr {
			refval = reflect.Indirect(refval)
		}
		// check nil first, else Kind() will crash
		if refval.Interface() == nil {
			return 0
		}
		vv := reflect.ValueOf(refval.Interface())
		switch vv.Type().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return Float64(vv.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return Float64(vv.Uint())
		case reflect.Float32, reflect.Float64:
			return (vv.Float())
		case reflect.String:
			return Float64(vv.String())
		case reflect.Array, reflect.Slice:
			switch vv.Type().Elem().Kind() {
			case reflect.Uint8:
				data := refval.Interface().([]byte)
				return Float64(string(data))
			}
		case reflect.Bool:
			return Float64(vv.Bool())
		case reflect.Complex128, reflect.Complex64:
			return Float64(vv.Complex())
		//时间类型
		case reflect.Struct:
			if vv.Type().ConvertibleTo(reflect.TypeOf(_DEFAULT_TIME)) {
				return float64(refval.Convert(reflect.TypeOf(_DEFAULT_TIME)).Interface().(time.Time).Unix())
			}
		}
	}
	return 0
}

// return struct fileds name lower-case
func GetStructNameLower(value interface{}) []string {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	switch f := val; f.Kind() {
	case reflect.Struct:
		m := make([]string, 0, f.NumField())
		for i := 0; i < f.NumField(); i++ {
			field := f.Type().Field(i)
			if field.Name == "" {
				continue
			}
			m = append(m, strings.ToLower(field.Name))
		}
		return m

	default:
		return []string{}
	}
}

//return map's objects array in order by int index
func Map2SliceByIndexInt(src map[int]interface{}, asc ...bool) []interface{} {
	idx := make([]int, 0, len(src))
	for k, _ := range src {
		idx = append(idx, k)
	}
	sort.Ints(idx)
	returns := make([]interface{}, 0, len(src))
	for _, v := range idx {
		returns = append(returns, src[v])
	}
	return returns
}
