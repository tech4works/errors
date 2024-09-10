package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func callerInfos(skip int) (fileName string, line string, funcName string) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		pc, file, lineNo, _ = runtime.Caller(1)
	}
	funcInfo := runtime.FuncForPC(pc).Name()
	dir, fileBase := filepath.Split(file)
	dirBase := filepath.Base(dir)
	name := formatFuncName(funcInfo)

	if lineNo < 1 {
		lineNo = 1
	}

	return dirBase + "/" + fileBase, strconv.Itoa(lineNo), name
}

func buildMessage(v ...any) string {
	var ss []any
	for _, i := range v {
		ss = append(ss, toString(i))
	}
	return cleanMessage(strings.TrimRight(fmt.Sprintln(ss...), "\n"))
}

func buildMessageByFormat(format string, v ...any) string {
	return cleanMessage(fmt.Sprintf(format, filterMsg(v...)...))
}

func cleanMessage(msg string) string {
	msg = strings.ReplaceAll(msg, "[STACK]", "")
	msg = strings.ReplaceAll(msg, "[CAUSE]", "")

	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(msg, " ")
}

func filterMsg(v ...any) []any {
	for i, iv := range v {
		ivError, ok := iv.(error)
		if ok {
			errDetail := Details(ivError)
			v[i] = errDetail.message
		}
	}
	return v
}

func formatFuncName(name string) string {
	name = path.Base(name)
	split := strings.Split(name, ".")
	return split[len(split)-1]
}

func toString(a any) string {
	s, err := toStringWithErr(a)
	if err != nil {
		return fmt.Sprint(a)
	}
	return s
}

func toStringWithErr(a any) (string, error) {
	reflectValue := reflect.ValueOf(a)

	switch reflectValue.Kind() {
	case reflect.String:
		return reflectValue.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(reflectValue.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(reflectValue.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(reflectValue.Float(), 'g', -1, 64), nil
	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(reflectValue.Complex(), 'g', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(reflectValue.Bool()), nil
	case reflect.Array, reflect.Slice:
		if reflectValue.Type().Elem().Kind() == reflect.Uint8 {
			return string(reflectValue.Bytes()), nil
		}
		marshal, _ := json.Marshal(reflectValue.Interface())
		return string(marshal), nil
	case reflect.Map, reflect.Struct:
		marshal, _ := json.Marshal(reflectValue.Interface())
		return string(marshal), nil
	case reflect.Ptr, reflect.Interface:
		if reflectValue.IsNil() {
			return "", errors.New("error convert to string, it is null")
		} else if err, ok := a.(error); ok {
			if IsDetailed(err) {
				details := Details(err)
				return details.Message(), nil
			}
			return err.Error(), nil
		}
		return toStringWithErr(reflectValue.Elem().Interface())
	default:
		return "", fmt.Errorf("error convert to string, unsupported type %s", reflectValue.Kind().String())
	}
}
