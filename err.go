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
	code     string
	message  string
	metadata map[string]any
	stack    string
}

func New(msg ...any) *Err {
	file, line, funcName := callerInfos(2)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		stack:    string(debug.Stack()),
	}
}

func NewWithCode(code string, msg ...any) *Err {
	file, line, funcName := callerInfos(2)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessage(msg...),
		stack:    string(debug.Stack()),
	}
}

func NewWithMetadata(metadata map[string]any, msg ...any) *Err {
	file, line, funcName := callerInfos(2)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		metadata: metadata,
		stack:    string(debug.Stack()),
	}
}

func NewWithCodeAndMetadata(code string, metadata map[string]any, msg ...any) *Err {
	file, line, funcName := callerInfos(2)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessage(msg...),
		metadata: metadata,
		stack:    string(debug.Stack()),
	}
}

func Newf(format string, msg ...any) *Err {
	file, line, funcName := callerInfos(2)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		stack:    string(debug.Stack()),
	}
}

func NewWithCodef(format, code string, msg ...any) *Err {
	file, line, funcName := callerInfos(2)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessageByFormat(format, msg...),
		stack:    string(debug.Stack()),
	}
}

func NewWithMetadataf(format string, metadata map[string]any, msg ...any) *Err {
	file, line, funcName := callerInfos(2)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		metadata: metadata,
		stack:    string(debug.Stack()),
	}
}

func NewWithCodeAndMetadataf(format, code string, metadata map[string]any, msg ...any) *Err {
	file, line, funcName := callerInfos(2)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessageByFormat(format, msg...),
		metadata: metadata,
		stack:    string(debug.Stack()),
	}
}

func NewWithSkipCaller(skipCaller int, msg ...any) *Err {
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		stack:    string(debug.Stack()),
	}
}

func NewWithSkipCallerAndCode(skipCaller int, code string, msg ...any) *Err {
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessage(msg...),
		stack:    string(debug.Stack()),
	}
}

func NewWithAll(skipCaller int, code string, metadata map[string]any, msg ...any) *Err {
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessage(msg...),
		metadata: metadata,
		stack:    string(debug.Stack()),
	}
}

func NewWithSkipCallerWithCodef(skipCaller int, format, code string, msg ...any) *Err {
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessageByFormat(format, msg...),
		stack:    string(debug.Stack()),
	}
}

func NewWithSkipCallerf(skipCaller int, format string, msg ...any) *Err {
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		stack:    string(debug.Stack()),
	}
}

func NewWithAllf(skipCaller int, format, code string, metadata map[string]any, msg ...any) *Err {
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessageByFormat(format, msg...),
		metadata: metadata,
		stack:    string(debug.Stack()),
	}
}

func (e *Err) Error() string {
	var code string
	if len(e.Code()) > 0 {
		code = fmt.Sprintf("[CODE]: %s ", e.Code())
	}
	var metadata string
	if len(e.Metadata()) > 0 {
		metadata = fmt.Sprintf("[METADATA]: %s", toString(e.Metadata()))
	}
	return fmt.Sprint(code, "[CAUSE]: ", e.Cause().Error(), metadata, " [STACK]: ", e.Stack())
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

func (e *Err) Code() string {
	return e.code
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

func (e *Err) Metadata() map[string]any {
	return e.metadata
}

func (e *Err) Stack() string {
	return e.stack
}

func (e *Err) String() string {
	return e.Error()
}
