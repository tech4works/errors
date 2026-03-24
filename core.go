package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const regex = `^(?:\[CODE]: (.+?) )?` + // 1: code (optional)
	`\[CAUSE]: \(([^:]+):(\d+)\) ([^:]+): (.+?)` + // 2: file, 3: line, 4: func, 5: msg
	`(?: \[METADATA]: (.+?))?` + // 6: metadata (optional)
	` \[STACK]:\s*([\s\S]+)$` // 7: stack

const inheritSep = "----------------\n\t\t|\tinherited by: "

var compiledRegex = regexp.MustCompile(regex)

func Is(err, target error) bool {
	if err == nil || target == nil {
		return false
	}

	var mErr, mTarget *Err

	errString := err.Error()
	if As(err, &mErr) {
		if mErr.IsTarget() {
			return false
		}
		errString = mErr.Snapshot()
	}

	targetString := target.Error()
	targetContains := targetString
	if As(target, &mTarget) {
		targetString = mTarget.Snapshot()
		if mTarget.HasCode() {
			targetContains = mTarget.Code()
		} else {
			targetContains = mTarget.Message()
		}
	}

	return errString == targetString || strings.Contains(errString, targetContains)
}

func IsNot(err, target error) bool {
	return !Is(err, target)
}

func ContainsCode(err error, code string) bool {
	if err == nil || code == "" {
		return false
	}

	pattern := fmt.Sprintf(`\[CODE\]:\s*%s(\s|\])`, regexp.QuoteMeta(code))
	re := regexp.MustCompile(pattern)

	return re.MatchString(err.Error())
}

func NotContainsCode(err error, code string) bool {
	return !ContainsCode(err, code)
}

func Contains(errs []error, target error) bool {
	if len(errs) == 0 || target == nil {
		return false
	}

	for _, err := range errs {
		if Is(err, target) {
			return true
		}
	}
	return false
}

func NotContains(errs []error, target error) bool {
	return !Contains(errs, target)
}

func Only(errs []error, target error) bool {
	if len(errs) == 0 || target == nil {
		return false
	}

	for _, err := range errs {
		if IsNot(err, target) {
			return false
		}
	}
	return true
}

func NotOnly(errs []error, target error) bool {
	return !Only(errs, target)
}

func As(err error, target any) bool {
	if err == nil || target == nil {
		return false
	}

	if t, ok := target.(**Err); ok {
		extracted, parsed := extract(err, 3)
		if parsed == nil {
			return false
		}
		if extracted {
			*t = parsed
		}
		return extracted
	}

	if t, ok := target.(*Err); ok {
		extracted, parsed := extract(err, 3)
		if extracted && parsed != nil {
			*t = *parsed
		}
		return extracted
	}

	if isValidAsTarget(target) {
		return errors.As(err, target)
	}

	return false
}

func Wrap(err error) *Err {
	if err == nil {
		return nil
	}
	_, e := extract(err, 3)
	return e
}

func Inherit(err error, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCaller(3, msg...)

	extracted, parent := extract(err, 3)

	e.parent = parent
	if e.Message() == "<empty>" {
		e.message = parent.message
	}
	if extracted {
		e.code = parent.code
		e.metadata = parent.metadata
	}

	return e
}

func Inheritf(err error, format string, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCallerf(3, format, msg...)

	extracted, parent := extract(err, 3)

	e.parent = parent
	if e.Message() == "<empty>" {
		e.message = parent.message
	}
	if extracted {
		e.code = parent.code
		e.metadata = parent.metadata
	}

	return e
}

func InheritWithSkipCaller(err error, skipCaller int, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCaller(skipCaller+1, msg...)

	extracted, parent := extract(err, skipCaller+1)

	e.parent = parent
	if e.Message() == "<empty>" {
		e.message = parent.message
	}
	if extracted {
		e.code = parent.code
		e.metadata = parent.metadata
	}

	return e
}

func InheritWithCode(err error, code string, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithCode(code, msg...)

	extracted, parent := extract(err, 3)

	e.parent = parent

	if e.Message() == "<empty>" {
		e.message = parent.message
	}

	if extracted {
		e.metadata = parent.metadata
	}

	return e
}

func InheritWithSkipCallerf(err error, skipCaller int, format string, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCallerf(skipCaller+1, format, msg...)

	extracted, parent := extract(err, 3)

	e.parent = parent
	if e.Message() == "<empty>" {
		e.message = parent.message
	}

	if extracted {
		e.code = parent.code
		e.metadata = parent.metadata
	}

	return e
}

func InheritWithAll(err error, skipCaller int, code string, metadata map[string]any, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithAll(skipCaller, code, metadata, msg...)

	extracted, parent := extract(err, 3)

	e.parent = parent
	if e.Message() == "<empty>" {
		e.message = parent.message
	}

	if extracted {
		for k, v := range parent.metadata {
			if _, exists := e.metadata[k]; !exists {
				e.metadata[k] = v
			}
		}
	}

	return e
}

func InheritWithSkipCallerAndCode(err error, skipCaller int, code string, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCallerAndCode(skipCaller, code, msg...)

	extracted, parent := extract(err, 3)

	e.parent = parent
	if e.Message() == "<empty>" {
		e.message = parent.message
	}

	if extracted {
		e.code = parent.code
		e.metadata = parent.metadata
	}

	return e
}

func InheritAsSlice(err error, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{Inherit(err, msg...)}
}

func InheritAsSlicef(err error, format string, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{Inheritf(err, format, msg...)}
}

func InheritWithSkipCallerAsSlice(err error, skipCaller int, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{InheritWithSkipCaller(err, skipCaller, msg...)}
}

func InheritWithCodeAsSlice(err error, code string, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{InheritWithCode(err, code, msg...)}
}

func InheritWithSkipCallerAsSlicef(err error, skipCaller int, format string, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{InheritWithSkipCallerf(err, skipCaller, format, msg...)}
}

func InheritWithAllAsSlice(err error, skipCaller int, code string, metadata map[string]any, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{InheritWithAll(err, skipCaller, code, metadata, msg...)}
}

func InheritWithSkipCallerAndCodeAsSlice(err error, skipCaller int, code string, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{InheritWithSkipCallerAndCode(err, skipCaller, code, msg...)}
}

func Join(errs []error, sep string) error {
	return errors.New(JoinToString(errs, sep))
}

func JoinInherit(errs []error, sep string, msg ...any) error {
	if errs == nil {
		return nil
	}
	return Inherit(Join(errs, sep), msg...)
}

func JoinInheritf(errs []error, sep, format string, msg ...any) *Err {
	if errs == nil {
		return nil
	}
	return Inheritf(Join(errs, sep), format, msg...)
}

func JoinInheritAsSlice(errs []error, sep string, msg ...any) []error {
	if errs == nil {
		return nil
	}
	return InheritAsSlice(Join(errs, sep), msg...)
}

func JoinInheritAsSlicef(errs []error, sep, format string, msg ...any) []error {
	if errs == nil {
		return nil
	}
	return InheritAsSlicef(Join(errs, sep), format, msg...)
}

func JoinInheritWithSkipCaller(errs []error, sep string, skipCaller int, msg ...any) *Err {
	if errs == nil {
		return nil
	}

	return InheritWithSkipCaller(Join(errs, sep), skipCaller, msg...)
}

func JoinInheritWithCode(errs []error, sep, code string, msg ...any) *Err {
	if errs == nil {
		return nil
	}

	return InheritWithCode(Join(errs, sep), code, msg...)
}

func JoinInheritWithSkipCallerf(errs []error, sep string, skipCaller int, format string, msg ...any) *Err {
	if errs == nil {
		return nil
	}

	return InheritWithSkipCallerf(Join(errs, sep), skipCaller, format, msg...)
}

func JoinInheritWithAll(errs []error, sep string, skipCaller int, code string, metadata map[string]any, msg ...any) *Err {
	if errs == nil {
		return nil
	}

	return InheritWithAll(Join(errs, sep), skipCaller, code, metadata, msg...)
}

func JoinInheritWithSkipCallerAndCode(errs []error, sep string, skipCaller int, code string, msg ...any) *Err {
	if errs == nil {
		return nil
	}

	return InheritWithSkipCallerAndCode(Join(errs, sep), skipCaller, code, msg...)
}

func JoinToString(errs []error, sep string) (result string) {
	if len(errs) == 0 {
		return ""
	}
	for i, err := range errs {
		dt := Wrap(err)
		if dt != nil {
			result += dt.message
		}
		if i < len(errs)-1 {
			result += sep
		}
	}
	return result
}

func extract(err error, skipCaller int) (bool, *Err) {
	if err == nil {
		return false, nil
	}

	var e *Err
	if errors.As(err, &e) {
		return true, e
	}

	rg := compiledRegex
	matches := rg.FindStringSubmatch(err.Error())

	if len(matches) == 0 {
		file, line, funcName := callerInfos(skipCaller + 1)
		return false, &Err{
			file:     file,
			line:     line,
			funcName: funcName,
			message:  buildMessage(err.Error()),
			stack:    buildDebugStack(),
		}
	}

	e = &Err{
		file:     matches[2],
		line:     matches[3],
		funcName: matches[4],
		message:  matches[5],
		stack:    []byte(matches[7]),
	}

	if matches[1] != "" {
		e.code = matches[1]
	}

	if matches[6] != "" {
		metaStr := matches[6]

		var meta map[string]any
		if err = json.Unmarshal([]byte(metaStr), &meta); err == nil {
			e.metadata = meta
		} else {
			e.metadata = map[string]any{
				"raw": metaStr,
			}
		}
	}

	return true, e
}

func isValidAsTarget(target any) bool {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false
	}
	elem := rv.Type().Elem()
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	return elem.Implements(errorType) || elem.Kind() == reflect.Interface
}
