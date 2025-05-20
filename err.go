package errors

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
)

type Err struct {
	file     string
	line     string
	funcName string
	message  string
	stack    string
}

func New(args ...any) *Err {
	msg := buildMessage(args...)
	file, line, funcName := callerInfos(2)
	stack := debug.Stack()
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  msg,
		stack:    string(stack),
	}
}

func Newf(format string, args ...any) *Err {
	msg := buildMessageByFormat(format, args...)
	file, line, funcName := callerInfos(2)
	stack := debug.Stack()
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  msg,
		stack:    string(stack),
	}
}

func NewSkipCaller(skipCaller int, args ...any) *Err {
	msg := buildMessage(args...)
	file, line, funcName := callerInfos(skipCaller + 1)
	stack := debug.Stack()
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  msg,
		stack:    string(stack),
	}
}

func NewSkipCallerf(skipCaller int, format string, args ...any) *Err {
	msg := buildMessageByFormat(format, args...)
	file, line, funcName := callerInfos(skipCaller + 1)
	stack := debug.Stack()
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  msg,
		stack:    string(stack),
	}
}

func (e *Err) Error() string {
	return fmt.Sprint("[CAUSE]: ", e.Cause().Error(), " [STACK]: ", e.stack)
}

func (e *Err) PrintStackTrace() {
	fmt.Print(e.stack)
}

func (e *Err) PrintCause() {
	fmt.Print(e.Cause())
}

func (e *Err) Cause() error {
	return errors.New(fmt.Sprint("(", e.file, ":", e.line, ")", " ", e.funcName, ": ", e.message))
}

func (e *Err) Message() string {
	return e.message
}

func (e *Err) File() string {
	return e.file
}

func (e *Err) Line() int {
	lineNumber, _ := strconv.Atoi(e.line)
	return lineNumber
}

func (e *Err) Func() string {
	return e.funcName
}

func (e *Err) Stack() string {
	return e.stack
}

func (e *Err) String() string {
	return e.Cause().Error()
}
