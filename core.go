package errors

import (
	"errors"
	"fmt"
	"regexp"
	"runtime/debug"
	"strings"
)

const regex = `\[CAUSE]: \(([^:]+):(\d+)\) ([^:]+): (.+?) \[STACK]:\s*([\s\S]+)`

func Is(err, target error) bool {
	if Match(err) {
		wrapped := Wrap(err)
		err = errors.New(wrapped.Message())
	}

	if Match(target) {
		wrapped := Wrap(target)
		target = errors.New(wrapped.Message())
	}

	return err != nil && target != nil && err.Error() == target.Error()
}

func IsNot(err, target error) bool {
	return !Is(err, target)
}

func Contains(err, target error) bool {
	if Match(err) {
		errDetails := Wrap(err)
		err = errors.New(errDetails.Message())
	}

	if Match(target) {
		errDetails := Wrap(target)
		target = errors.New(errDetails.Message())
	}

	return err != nil && target != nil && strings.Contains(err.Error(), target.Error())
}

func NotContains(err, target error) bool {
	return !Contains(err, target)
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

func extract(err error, skip int) *Err {
	if err == nil {
		return nil
	}

	var file string
	var line string
	var funcName string
	var message string
	var stack string

	rg := regexp.MustCompile(regex)
	matches := rg.FindStringSubmatch(err.Error())

	if len(matches) > 0 {
		file = matches[1]
		line = matches[2]
		funcName = matches[3]
		message = matches[4]
		stack = matches[5]
	} else {
		file, line, funcName = callerInfos(skip)
		stack = string(debug.Stack())
		message = buildMessage(err.Error())
	}

	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  message,
		stack:    stack,
	}
}
