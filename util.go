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
	"sync"
)

var (
	prefixesToRemove []string
	prefixMutex      sync.RWMutex
	framesToSanitize []string
	frameMutex       sync.RWMutex
)

func SetPrefixesToSanitize(prefixes ...string) {
	prefixMutex.Lock()
	defer prefixMutex.Unlock()
	prefixesToRemove = prefixes
}

func SetFrameToSanitize(frames ...string) {
	frameMutex.Lock()
	defer frameMutex.Unlock()
	framesToSanitize = frames
}

func renderStackByPolicy(raw []byte) string {
	switch getPolicy() {
	case PolicyNormal:
		return normalizeStack(sanitizeBoilerplateFrames(raw, 10))
	case PolicyNative:
		return normalizeStack(sanitizeBoilerplateFrames(raw, 5))
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

	// Sanitize frames that should be removed entirely
	s = sanitizeFramesInStack(s)

	// Remove parâmetros das funções, deixando só o nome
	// Exemplo: "func({{0x...}, ...}, ...)" -> "func()"
	re := regexp.MustCompile(`\([^)]*\)`)
	s = re.ReplaceAllString(s, "()")

	// Remove offsets hexadecimais (+0x5e, +0xe3, etc)
	// Exemplo: "/path/file.go:42 +0x5e" -> "/path/file.go:42"
	hexOffsetRe := regexp.MustCompile(`\s+\+0x[0-9a-fA-F]+`)
	s = hexOffsetRe.ReplaceAllString(s, "")

	// Sanitiza prefixos configurados
	s = sanitizePrefixes(s)

	return s
}

func sanitizeFramesInStack(s string) string {
	frameMutex.RLock()
	defer frameMutex.RUnlock()

	if len(framesToSanitize) == 0 {
		return s
	}

	lines := strings.Split(s, "\n")
	var result []string

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		shouldSkip := false

		// Check if current line matches any frame to sanitize
		for _, frameMatch := range framesToSanitize {
			if strings.Contains(line, frameMatch) {
				shouldSkip = true
				break
			}
		}

		// If this line matches, skip it and the next line (file line)
		if shouldSkip {
			// Skip the next line if it looks like a file line
			if i+1 < len(lines) && looksLikeFileLine(lines[i+1]) {
				i++ // Skip the file line
			}
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func sanitizePrefixes(s string) string {
	prefixMutex.RLock()
	defer prefixMutex.RUnlock()

	for _, prefix := range prefixesToRemove {
		s = strings.ReplaceAll(s, prefix, "")
	}
	return s
}

func sanitizeBoilerplateFrames(raw []byte, keepFirst int) []byte {
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

	// keepFirst - mantém os primeiros N frames (do topo)
	if keepFirst > 0 && len(frames) > keepFirst {
		frames = frames[:keepFirst]
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
		"github.com/tech4works/errors", // Própria lib de errors
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
	ss := filterMsg(v...)
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
	result := make([]any, len(v))
	for i, iv := range v {
		if ivError, ok := iv.(error); ok {
			errDetail := Wrap(ivError)
			if errDetail != nil {
				result[i] = errDetail.message
				continue
			}
		}
		result[i] = iv
	}
	return result
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

func optionalField(label, value string) string {
	if len(value) == 0 {
		return ""
	}
	return label + value
}
