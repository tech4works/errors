package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
)

const regex = `^(?:\[CODE]: (.+?) )?` + // 1: code (optional)
	`\[CAUSE]: \(([^:]+):(\d+)\) ([^:]+): (.+?)` + // 2: file, 3: line, 4: func, 5: msg
	`(?: \[METADATA]: (.+?))?` + // 6: metadata (optional)
	` \[STACK]:\s*([\s\S]+)$` // 7: stack

func Is(err, target error) bool {
	if Match(err) {
		wrapped := Wrap(err)
		err = errors.New(wrapped.Simple())
	}

	if Match(target) {
		wrapped := Wrap(target)
		target = errors.New(wrapped.Simple())
	}

	return err != nil && target != nil && (err.Error() == target.Error() || contains(err, target))
}

func IsNot(err, target error) bool {
	return !Is(err, target)
}

func Match(err error) bool {
	regex := regexp.MustCompile(regex)
	return err != nil && regex.MatchString(err.Error())
}

func Wrap(err error, msg ...any) *Err {
	e := extract(err, 3)
	if e == nil {
		return nil
	} else if len(msg) == 0 {
		return e
	}
	e.message = fmt.Sprintf("%s err: %s", buildMessage(msg...), e.message)
	return e
}

func WrapSkipCaller(err error, skip int, msg ...any) *Err {
	e := extract(err, skip)
	if e == nil {
		return nil
	}

	file, line, funcName := callerInfos(skip)
	e.file = file
	e.line = line
	e.funcName = funcName
	e.stack = string(debug.Stack())

	if len(msg) == 0 {
		return e
	}
	e.message = fmt.Sprintf("%s err: %s", buildMessage(msg...), e.message)
	return e
}

func WrapSkipCallerf(err error, skip int, format string, msg ...any) *Err {
	e := extract(err, skip)
	if e == nil {
		return nil
	}

	file, line, funcName := callerInfos(skip)
	e.file = file
	e.line = line
	e.funcName = funcName
	e.stack = string(debug.Stack())
	if len(msg) == 0 {
		return e
	}
	e.message = fmt.Sprintf("%s err: %s", buildMessageByFormat(format, msg...), e.message)
	return e
}

func Wrapf(err error, format string, msg ...any) *Err {
	e := extract(err, 3)
	if e == nil {
		return nil
	}

	e.message = fmt.Sprintf("%s err: %s", buildMessageByFormat(format, msg...), e.message)
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
	if Match(err) {
		errDetails := Wrap(err)
		err = errors.New(errDetails.Simple())
	}

	if Match(target) {
		errDetails := Wrap(target)
		if len(errDetails.Code()) > 0 {
			target = errors.New(errDetails.Code())
		} else {
			target = errors.New(errDetails.Message())
		}
	}

	return err != nil && target != nil && strings.Contains(err.Error(), target.Error())
}

func extract(err error, skip int) *Err {
	if err == nil {
		return nil
	}

	rg := regexp.MustCompile(regex)
	matches := rg.FindStringSubmatch(err.Error())

	if len(matches) == 0 {
		file, line, funcName := callerInfos(skip)
		return &Err{
			file:     file,
			line:     line,
			funcName: funcName,
			message:  buildMessage(err.Error()),
			stack:    string(debug.Stack()),
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

	return e
}
