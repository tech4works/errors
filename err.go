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

func NewWithCodef(code, format string, msg ...any) *Err {
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

func NewWithMetadataf(metadata map[string]any, format string, msg ...any) *Err {
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

func NewWithCodeAndMetadataf(code string, metadata map[string]any, format string, msg ...any) *Err {
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

func NewWithSkipCallerAndCodef(skipCaller int, code, format string, msg ...any) *Err {
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

func NewWithAllf(skipCaller int, code string, metadata map[string]any, format string, msg ...any) *Err {
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

func NewAsSlice(msg ...any) []error {
	return []error{NewWithSkipCaller(2, msg...)}
}

func NewAsSlicef(format string, msg ...any) []error {
	return []error{NewWithSkipCallerf(2, format, msg...)}
}

func NewWithCodeAsSlice(code string, msg ...any) []error {
	return []error{NewWithSkipCallerAndCode(2, code, msg...)}
}

func NewWithMetadataAsSlice(metadata map[string]any, msg ...any) []error {
	return []error{NewWithAll(2, "", metadata, msg...)}
}

func NewWithCodeAndMetadataAsSlice(code string, metadata map[string]any, msg ...any) []error {
	return []error{NewWithAll(2, code, metadata, msg...)}
}

func NewWithCodeAsSlicef(code, format string, msg ...any) []error {
	return []error{NewWithSkipCallerAndCodef(2, code, format, msg...)}
}

func NewWithMetadataAsSlicef(metadata map[string]any, format string, msg ...any) []error {
	return []error{NewWithAllf(2, "", metadata, format, msg...)}
}

func NewWithCodeAndMetadataAsSlicef(code string, metadata map[string]any, format string, msg ...any) []error {
	return []error{NewWithAllf(2, code, metadata, format, msg...)}
}

func NewWithSkipCallerAsSlice(skipCaller int, msg ...any) []error {
	return []error{NewWithSkipCaller(skipCaller+1, msg...)}
}

func NewWithSkipCallerAndCodeAsSlice(skipCaller int, code string, msg ...any) []error {
	return []error{NewWithSkipCallerAndCode(skipCaller+1, code, msg...)}
}

func NewWithAllAsSlice(skipCaller int, code string, metadata map[string]any, msg ...any) []error {
	return []error{NewWithAll(skipCaller+1, code, metadata, msg...)}
}

func NewWithSkipCallerAndCodeAsSlicef(skipCaller int, code, format string, msg ...any) []error {
	return []error{NewWithSkipCallerAndCodef(skipCaller+1, code, format, msg...)}
}

func NewWithSkipCallerAsSlicef(skipCaller int, format string, msg ...any) []error {
	return []error{NewWithSkipCallerf(skipCaller+1, format, msg...)}
}

func NewWithAllAsSlicef(skipCaller int, code string, metadata map[string]any, format string, msg ...any) []error {
	return []error{NewWithAllf(skipCaller+1, code, metadata, format, msg...)}
}

func NewWithRawDataAsSlice(file, line, funcName, code, message string, metadata map[string]any, stack []byte) []error {
	return []error{NewWithRawData(file, line, funcName, code, message, metadata, stack)}
}

func (e *Err) Error() string {
	if e.target {
		s := "[TARGET]: true"
		if e.HasCode() {
			s += " [CODE]: " + e.code
		}
		if len(e.message) > 0 {
			s += " [MESSAGE]: " + e.message
		}
		return s
	}

	code := ""
	if e.HasCode() {
		code = "[CODE]: " + e.code + " "
	}
	metadata := ""
	if e.HasMetadata() {
		metadata = " [METADATA]: " + toString(e.Metadata())
	}
	stack := e.Stack()
	if e.parent != nil {
		stack += "\n" + inheritSep + e.parent.Error()
	}
	return code + "[CAUSE]: " + e.Cause().Error() + metadata + " [STACK]: " + stack
}

func (e *Err) Raw() string {
	if e.target {
		return e.Error()
	}

	code := ""
	if e.HasCode() {
		code = "[CODE]: " + e.code + " "
	}
	metadata := ""
	if e.HasMetadata() {
		metadata = " [METADATA]: " + toString(e.Metadata())
	}
	return code + "[CAUSE]: " + e.Cause().Error() + metadata + " [STACK]: " + normalizeStack(e.stack)
}

func (e *Err) PrintStackTrace() {
	fmt.Println(e.stack)
}

func (e *Err) PrintCause() {
	fmt.Println(e.Cause())
}

func (e *Err) Cause() error {
	if e.target {
		return errors.New(e.Snapshot())
	}
	return errors.New(fmt.Sprint("(", e.file, ":", e.line, ")", " ", e.funcName, ": ", e.Snapshot()))
}

func (e *Err) IsTarget() bool {
	return e.target
}

func (e *Err) HasCode() bool {
	return len(e.code) > 0
}

func (e *Err) Code() string {
	return e.code
}

func (e *Err) HasMessage() bool {
	return len(e.message) > 0
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
	return len(e.metadata) > 0
}

func (e *Err) Metadata() map[string]any {
	return e.metadata
}

func (e *Err) Stack() string {
	stack := renderStackByPolicy(e.stack)
	if e.parent != nil {
		stack += "\n" + inheritSep + e.parent.Error()
	}
	return stack
}

func (e *Err) Snapshot() string {
	s := e.message
	if e.HasCode() {
		s = "[CODE]: " + e.code + " [MESSAGE]: " + s
	}
	return s
}

func (e *Err) String() string {
	return e.Error()
}
