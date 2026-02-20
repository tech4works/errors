package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const regex = `^(?:\[CODE]: (.+?) )?` + // 1: code (optional)
	`\[CAUSE]: \(([^:]+):(\d+)\) ([^:]+): (.+?)` + // 2: file, 3: line, 4: func, 5: msg
	`(?: \[METADATA]: (.+?))?` + // 6: metadata (optional)
	` \[STACK]:\s*([\s\S]+)$` // 7: stack

const inheritSep = "----------------\n\t\t|\tinherited by: "

func Is(err, target error) bool {
	if err == nil || target == nil {
		return false
	}

	var mErr, mTarget Err

	errString := err.Error()
	if As(err, &mErr) {
		if mErr.IsTarget() {
			return false
		}
		errString = mErr.Snapshot()
	}

	targetString := target.Error()
	targetContains := targetString
	if As(err, &mTarget) {
		targetString = mTarget.Snapshot()
		if mTarget.HasCode() {
			targetContains = mTarget.Code()
		} else {
			targetContains = mTarget.Snapshot()
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
	} else if errors.As(err, target) {
		return true
	} else if t, ok := target.(**Err); ok {
		if !regexp.MustCompile(regex).MatchString(err.Error()) {
			return false
		}

		extracted, parsed := extract(err, 3)

		*t = parsed
		return extracted
	}

	return false
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

	e := NewWithSkipCallerf(skipCaller, format, msg...)

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

func Wrap(err error) *Err {
	_, extract := extract(err, 3)
	return extract
}

func Join(errs []error, sep string) error {
	return errors.New(JoinToString(errs, sep))
}

func JoinToString(errs []error, sep string) (result string) {
	for i, err := range errs {
		dt := Wrap(err)
		result += dt.message
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

	rg := regexp.MustCompile(regex)
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
