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
	"runtime/debug"
	"strconv"
	"strings"
)

func renderStackByPolicy(raw []byte) string {
	switch getPolicy() {
	case PolicyDetailed:
		return normalizeStack(raw)
	case PolicyNormal:
		return normalizeStack(filterBoilerplateFrames(raw, 0))
	case PolicyNative:
		return normalizeStack(filterBoilerplateFrames(raw, 5))
	default:
		return normalizeStack(raw)
	}
}

func normalizeStack(raw []byte) string {
	s := string(raw)

	s = strings.TrimSpace(s)
	if s == "" {
		return "<empty>"
	}

	return s
}

func filterBoilerplateFrames(raw []byte, keepLast int) []byte {
	lines := strings.Split(string(raw), "\n")

	type frame struct {
		fn   string
		file string
	}

	frames := make([]frame, 0, len(lines)/2)

	// debug.Stack geralmente vem assim:
	// goroutine X [running]:
	// runtime/debug.Stack(...)
	// ...
	// \t/path/file.go:line +0x...
	//
	// Vamos ignorar o header e tentar agrupar por pares (fn line + file line)
	for i := 0; i < len(lines)-1; i++ {
		fnLine := lines[i]
		fileLine := lines[i+1]

		// heurística: fileLine costuma começar com tab e conter ".go:"
		if !looksLikeFileLine(fileLine) {
			continue
		}

		if isBoilerplate(fnLine, fileLine) {
			i++ // consumiu o par
			continue
		}

		frames = append(frames, frame{
			fn:   strings.TrimRight(fnLine, "\r"),
			file: strings.TrimRight(fileLine, "\r"),
		})

		i++ // consumiu o par
	}

	// Se não conseguimos parsear frames, fallback: devolve raw
	if len(frames) == 0 {
		return raw
	}

	// keepLast
	if keepLast > 0 && len(frames) > keepLast {
		frames = frames[len(frames)-keepLast:]
	}

	// Remonta mantendo pares
	var out []string
	for _, fr := range frames {
		out = append(out, fr.fn, fr.file)
	}
	return []byte(strings.Join(out, "\n"))
}

func looksLikeFileLine(s string) bool {
	// debug.Stack usa "\t/path/file.go:123 +0x..."
	return strings.Contains(s, ".go:") && strings.HasPrefix(s, "\t")
}

func isBoilerplate(fnLine, fileLine string) bool {
	fn := strings.TrimSpace(fnLine)
	file := strings.TrimSpace(fileLine)

	// Ajuste livre: esses são defaults “bons” pra microserviço/gateway
	dropContains := []string{
		"/runtime/", "runtime.",
		"/testing/", "testing.",
		"/reflect/", "reflect.",
		"/net/http/", "net/http.",
		"/runtime/debug", "runtime/debug.",
		"/pkg/mod/",
	}

	for _, d := range dropContains {
		if strings.Contains(fn, d) || strings.Contains(file, d) {
			return true
		}
	}

	return false
}

func escapeSpecialChars(s string) string {
	replacer := strings.NewReplacer(
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
	)
	return replacer.Replace(s)
}

func callerInfos(skipCaller int) (fileName string, line string, funcName string) {
	pc, file, lineNo, ok := runtime.Caller(skipCaller + 1)
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
	ss := filterMsg(v)
	for _, i := range v {
		ss = append(ss, toString(i))
	}
	message := cleanMessage(strings.TrimRight(fmt.Sprintln(ss...), "\n"))
	if len(message) == 0 {
		message = "<empty>"
	}
	return escapeSpecialChars(message)
}

func buildDebugStack() []byte {
	return debug.Stack()
}

func buildMessageByFormat(format string, v ...any) string {
	return escapeSpecialChars(cleanMessage(fmt.Sprintf(format, filterMsg(v...)...)))
}

func cleanMessage(msg string) string {
	msg = strings.ReplaceAll(msg, "[CODE]", "")
	msg = strings.ReplaceAll(msg, "[METADATA]", "")
	msg = strings.ReplaceAll(msg, "[STACK]", "")
	msg = strings.ReplaceAll(msg, "[CAUSE]", "")

	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(msg, " ")
}

func filterMsg(v ...any) []any {
	for i, iv := range v {
		ivError, ok := iv.(error)
		if ok {
			errDetail := Wrap(ivError)
			if errDetail != nil {
				v[i] = errDetail.message
			}
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
	if a == nil {
		return "", errors.New("error convert to string, it is null")
	}

	reflectValue := reflect.ValueOf(a)
	reflectType := reflectValue.Type()

	if reflectValue.Kind() == reflect.Ptr || reflectValue.Kind() == reflect.Interface {
		if reflectValue.IsNil() {
			return "", errors.New("error convert to string, it is null")
		} else if implementsStringer(reflectType) {
			return reflectValue.Interface().(fmt.Stringer).String(), nil
		} else if implementsError(reflectType) {
			return reflectValue.Interface().(error).Error(), nil
		}
		return toStringWithErr(reflectValue.Elem().Interface())
	}

	if implementsStringer(reflectType) {
		return reflectValue.Interface().(fmt.Stringer).String(), nil
	} else if implementsError(reflectType) {
		return reflectValue.Interface().(error).Error(), nil
	}

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
	case reflect.Slice:
		if reflectValue.Type().Elem().Kind() == reflect.Uint8 {
			return string(reflectValue.Bytes()), nil
		}
		marshal, _ := json.Marshal(reflectValue.Interface())
		return string(marshal), nil
	case reflect.Array:
		if reflectValue.Type().Elem().Kind() == reflect.Uint8 {
			bytes := make([]byte, reflectValue.Len())
			for i := 0; i < reflectValue.Len(); i++ {
				bytes[i] = byte(reflectValue.Index(i).Uint())
			}
			return string(bytes), nil
		}
		marshal, _ := json.Marshal(reflectValue.Interface())
		return string(marshal), nil
	case reflect.Map, reflect.Struct:
		marshal, _ := json.Marshal(reflectValue.Interface())
		return string(marshal), nil
	default:
		return "", fmt.Errorf("error convert to string, unsupported type %s", reflectValue.Kind().String())
	}
}

func implementsStringer(reflectType reflect.Type) bool {
	if reflectType == nil {
		return false
	}
	return reflectType.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem())
}

func implementsError(reflectType reflect.Type) bool {
	if reflectType == nil {
		return false
	}
	return reflectType.Implements(reflect.TypeOf((*error)(nil)).Elem())
}
