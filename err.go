package errors

import (
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"
)

type Policy uint32

const (
	PolicyDetailed Policy = iota
	PolicyNormal
	PolicyNative
)

var policy atomic.Uint32

func init() {
	policy.Store(uint32(PolicyNormal))
}

func SetPolicy(m Policy) {
	policy.Store(uint32(m))
}

func getPolicy() Policy {
	return Policy(policy.Load())
}

type Err struct {
	parent *Err

	target   bool
	file     string
	line     string
	funcName string
	code     string
	message  string
	metadata map[string]any
	stack    []byte
}

func TargetWithCode(code string) *Err {
	return &Err{target: true, code: code}
}

func TargetWithMessage(msg ...any) *Err {
	return &Err{target: true, message: buildMessage(msg...)}
}

func New(msg ...any) *Err {
	file, line, funcName := callerInfos(2)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		stack:    buildDebugStack(),
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
		stack:    buildDebugStack(),
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
		stack:    buildDebugStack(),
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
		stack:    buildDebugStack(),
	}
}

func Newf(format string, msg ...any) *Err {
	file, line, funcName := callerInfos(2)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		stack:    buildDebugStack(),
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
		stack:    buildDebugStack(),
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
		stack:    buildDebugStack(),
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
		stack:    buildDebugStack(),
	}
}

func NewWithSkipCaller(skipCaller int, msg ...any) *Err {
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		stack:    buildDebugStack(),
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
		stack:    buildDebugStack(),
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
		stack:    buildDebugStack(),
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
		stack:    buildDebugStack(),
	}
}

func NewWithSkipCallerf(skipCaller int, format string, msg ...any) *Err {
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		stack:    buildDebugStack(),
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
		stack:    buildDebugStack(),
	}
}

func NewWithRawData(file, line, funcName, code, message string, metadata map[string]any, stack []byte) *Err {
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  message,
		metadata: metadata,
		stack:    stack,
	}
}

func (e *Err) Error() string {
	var code string
	if e.HasCode() {
		code = fmt.Sprintf("[CODE]: %s ", e.Code())
	}
	var metadata string
	if e.HasMetadata() {
		metadata = fmt.Sprintf(" [METADATA]: %s", toString(e.Metadata()))
	}

	stackView := renderStackByPolicy(e.stack)

	if e.parent != nil {
		stackView = stackView + "\n" + inheritSep + e.parent.Error()
	}

	return fmt.Sprint(code, "[CAUSE]: ", e.Cause().Error(), metadata, " [STACK]: ", e.Stack())
}

func (e *Err) PrintStackTrace() {
	fmt.Println(e.stack)
}

func (e *Err) PrintCause() {
	fmt.Println(e.Cause())
}

func (e *Err) Cause() error {
	return errors.New(fmt.Sprint("(", e.file, ":", e.line, ")", " ", e.funcName, ": ", e.Snapshot()))
}

func (e *Err) IsTarget() bool {
	return e.target
}

func (e *Err) HasCode() bool {
	return len(e.Code()) > 0
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

func (e *Err) HasMetadata() bool {
	return e.metadata != nil && len(e.metadata) > 0
}

func (e *Err) Metadata() map[string]any {
	return e.metadata
}

func (e *Err) Stack() string {
	return renderStackByPolicy(e.stack)
}

func (e *Err) Snapshot() string {
	s := e.Message()
	if e.HasCode() {
		s = fmt.Sprintf("[CODE]: %s [MESSAGE]: %s", e.Code(), s)
	}
	return s
}

func (e *Err) String() string {
	return e.Error()
}

func (e *Err) Raw() string {
	var code string
	if e.HasCode() {
		code = fmt.Sprintf("[CODE]: %s ", e.Code())
	}
	var metadata string
	if e.HasMetadata() {
		metadata = fmt.Sprintf(" [METADATA]: %s", toString(e.Metadata()))
	}

	return fmt.Sprint(code, "[CAUSE]: ", e.Cause().Error(), metadata, " [STACK]: ", normalizeStack(e.stack))
}
