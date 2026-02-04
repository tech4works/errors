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

func Is(err, target error) bool {
	if err == nil || target == nil {
		return false
	}

	errString := err.Error()
	if Match(err) {
		wrapped := Wrap(err)
		errString = wrapped.Snapshot()
	}

	targetString := target.Error()
	if Match(target) {
		wrapped := Wrap(target)
		targetString = wrapped.Snapshot()
	}

	return errString == targetString || contains(err, target)
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

func Match(err error) bool {
	regex := regexp.MustCompile(regex)
	return err != nil && regex.MatchString(err.Error())
}

func Inherit(err error, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCaller(2, msg...)

	extracted, parent := extract(err, 2)
	if extracted {
		if e.Message() == "<empty>" {
			e.message = parent.message
		}
		e.code = parent.code
		e.metadata = parent.metadata
		e.stack = inheritsStackWithError(parent, e.stack)
	} else {
		e.message = inheritsSimpleMessageByErrs(parent, e)
	}

	return e
}

func Inheritf(err error, format string, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCallerf(2, format, msg...)

	extracted, parent := extract(err, 2)
	if extracted {
		if e.Message() == "<empty>" {
			e.message = parent.message
		}
		e.code = parent.code
		e.metadata = parent.metadata
		e.stack = inheritsStackWithError(parent, e.stack)
	} else {
		e.message = inheritsSimpleMessageByErrs(parent, e)
	}
	return e
}

func InheritWithSkipCaller(err error, skipCaller int, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCaller(skipCaller, msg...)

	extracted, parent := extract(err, 2)
	if extracted {
		if e.Message() == "<empty>" {
			e.message = parent.message
		}
		e.code = parent.code
		e.metadata = parent.metadata
		e.stack = inheritsStackWithError(parent, e.stack)
	} else {
		e.message = inheritsSimpleMessageByErrs(parent, e)
	}
	return e
}

func InheritWithCode(err error, code string, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithCode(code, msg...)

	extracted, parent := extract(err, 2)
	if extracted {
		if e.Message() == "<empty>" {
			e.message = parent.message
		}
		e.metadata = parent.metadata
		e.stack = inheritsStackWithError(parent, e.stack)
	} else {
		e.message = inheritsSimpleMessageByErrs(parent, e)
	}
	return e
}

func InheritWithSkipCallerf(err error, skipCaller int, format string, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCallerf(skipCaller, format, msg...)

	extracted, parent := extract(err, 2)
	if extracted {
		if e.Message() == "<empty>" {
			e.message = parent.message
		}
		e.code = parent.code
		e.metadata = parent.metadata
		e.stack = inheritsStackWithError(parent, e.stack)
	} else {
		e.message = inheritsSimpleMessageByErrs(parent, e)
	}
	return e
}

func InheritWithAll(err error, skipCaller int, code string, metadata map[string]any, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithAll(skipCaller, code, metadata, msg...)

	extracted, parent := extract(err, 2)
	if extracted {
		if e.Message() == "<empty>" {
			e.message = parent.message
		}
		for k, v := range parent.metadata {
			if _, exists := e.metadata[k]; !exists {
				e.metadata[k] = v
			}
		}
		e.stack = inheritsStackWithError(parent, e.stack)
	} else {
		e.message = inheritsSimpleMessageByErrs(parent, e)
	}
	return e
}

func InheritWithSkipCallerAndCode(err error, skipCaller int, code string, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCallerAndCode(skipCaller, code, msg...)

	extracted, parent := extract(err, 2)
	if extracted {
		if e.Message() == "<empty>" {
			e.message = parent.message
		}
		e.code = parent.code
		e.metadata = parent.metadata
		e.stack = inheritsStackWithError(parent, e.stack)
	} else {
		e.message = inheritsSimpleMessageByErrs(parent, e)
	}
	return e
}

func Wrap(err error, msg ...any) *Err {
	extracted, e := extract(err, 2)
	if e == nil {
		return nil
	} else if len(msg) == 0 || !extracted {
		return e
	}

	e.message = inheritsMessage(e, buildMessage(msg...))
	e.stack = inheritsStack(e, buildDebugStack())

	return e
}

func WrapWithSkipCaller(err error, skip int, msg ...any) *Err {
	extracted, e := extract(err, skip)
	if e == nil {
		return nil
	}

	file, line, funcName := callerInfos(skip)
	e.file = file
	e.line = line
	e.funcName = funcName

	if len(msg) == 0 {
		return e
	}

	e.message = inheritsMessage(e, msg...)

	if !extracted {
		return e
	}

	e.stack = inheritsStack(e, buildDebugStack())

	return e
}

func WrapWithSkipCallerf(err error, skip int, format string, msg ...any) *Err {
	extracted, e := extract(err, skip)
	if e == nil {
		return nil
	}

	file, line, funcName := callerInfos(skip)
	e.file = file
	e.line = line
	e.funcName = funcName

	if len(msg) == 0 {
		return e
	}

	e.message = inheritsMessage(e, buildMessageByFormat(format, msg...))

	if !extracted {
		return e
	}

	e.stack = inheritsStack(e, buildDebugStack())

	return e
}

func Wrapf(err error, format string, msg ...any) *Err {
	extracted, e := extract(err, 2)
	if e == nil {
		return nil
	}

	e.message = inheritsMessage(e, buildMessageByFormat(format, msg...))

	if !extracted {
		return e
	}

	e.stack = inheritsStack(e, buildDebugStack())

	return e
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

func contains(err, target error) bool {
	errString := err.Error()
	if Match(err) {
		errDetails := Wrap(err)
		errString = errDetails.Snapshot()
	}

	targetString := target.Error()
	if Match(target) {
		errDetails := Wrap(target)
		if len(errDetails.Code()) > 0 {
			targetString = errDetails.Code()
		} else {
			targetString = errDetails.Message()
		}
	}

	return err != nil && target != nil && strings.Contains(errString, targetString)
}

func extract(err error, skip int) (bool, *Err) {
	if err == nil {
		return false, nil
	}

	rg := regexp.MustCompile(regex)
	matches := rg.FindStringSubmatch(err.Error())

	if len(matches) == 0 {
		file, line, funcName := callerInfos(skip + 1)
		return false, &Err{
			file:     file,
			line:     line,
			funcName: funcName,
			message:  buildMessage(err.Error()),
			stack:    buildDebugStack(),
		}
	}

	e := &Err{
		file:     matches[2],
		line:     matches[3],
		funcName: matches[4],
		message:  matches[5],
		stack:    matches[7],
	}

	// code opcional
	if matches[1] != "" {
		e.code = matches[1]
	}

	// metadata opcional
	if matches[6] != "" {
		metaStr := matches[6]

		var meta map[string]any
		if err := json.Unmarshal([]byte(metaStr), &meta); err == nil {
			e.metadata = meta
		} else {
			e.metadata = map[string]any{
				"raw": metaStr,
			}
		}
	}

	return true, e
}

func inheritsSimpleMessageByErrs(parent, e *Err) string {
	return escapeSpecialChars(fmt.Sprintf("%s | inherited by: %s", e.Message(), parent.Message()))
}

func inheritsMessage(parent *Err, msg ...any) string {
	return escapeSpecialChars(fmt.Sprintf("%s | inherited by: %s", buildMessage(msg...), parent.Snapshot()))
}

func inheritsStack(parent *Err, currentStack string) string {
	return escapeSpecialChars(fmt.Sprintf("%s----------------\n\t\t|\tinherited by: %s", currentStack, parent.stack))
}

func inheritsStackWithError(parent *Err, currentStack string) string {
	return escapeSpecialChars(fmt.Sprintf("%s----------------\n\t\t|\tinherited by: %s", currentStack, parent.Error()))
}
