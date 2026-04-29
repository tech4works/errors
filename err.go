package errors

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	file, line, funcName := callerInfos(1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		stack:    buildDebugStack(),
	}
}

func NewWithCode(code string, msg ...any) *Err {
	file, line, funcName := callerInfos(1)
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
	file, line, funcName := callerInfos(1)
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
	file, line, funcName := callerInfos(1)
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
	file, line, funcName := callerInfos(1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		stack:    buildDebugStack(),
	}
}

func NewWithCodef(code, format string, msg ...any) *Err {
	file, line, funcName := callerInfos(1)
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
	file, line, funcName := callerInfos(1)
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
	file, line, funcName := callerInfos(1)
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

func NewByParent(parent error, msg ...any) *Err {
	if parent == nil {
		return New(msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   Wrap(parent),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		stack:    buildDebugStack(),
	}
}

func NewByParentf(parent error, format string, msg ...any) *Err {
	if parent == nil {
		return Newf(format, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   Wrap(parent),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		stack:    buildDebugStack(),
	}
}

func NewByParentWithCode(parent error, code string, msg ...any) *Err {
	if parent == nil {
		return NewWithCode(code, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   Wrap(parent),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessage(msg...),
		stack:    buildDebugStack(),
	}
}

func NewByParentWithCodef(parent error, code, format string, msg ...any) *Err {
	if parent == nil {
		return NewWithCodef(code, format, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   Wrap(parent),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessageByFormat(format, msg...),
		stack:    buildDebugStack(),
	}
}

func NewByParentWithMetadata(parent error, metadata map[string]any, msg ...any) *Err {
	if parent == nil {
		return NewWithMetadata(metadata, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   Wrap(parent),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		metadata: metadata,
		stack:    buildDebugStack(),
	}
}

func NewByParentWithMetadataf(parent error, metadata map[string]any, format string, msg ...any) *Err {
	if parent == nil {
		return NewWithMetadataf(metadata, format, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   Wrap(parent),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		metadata: metadata,
		stack:    buildDebugStack(),
	}
}

func NewByParentWithCodeAndMetadata(parent error, code string, metadata map[string]any, msg ...any) *Err {
	if parent == nil {
		return NewWithCodeAndMetadata(code, metadata, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   Wrap(parent),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessage(msg...),
		metadata: metadata,
		stack:    buildDebugStack(),
	}
}

func NewByParentWithCodeAndMetadataf(parent error, code string, metadata map[string]any, format string, msg ...any) *Err {
	if parent == nil {
		return NewWithCodeAndMetadataf(code, metadata, format, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   Wrap(parent),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessageByFormat(format, msg...),
		metadata: metadata,
		stack:    buildDebugStack(),
	}
}

func NewAsSlice(msg ...any) []error {
	return []error{NewWithSkipCaller(1, msg...)}
}

func NewAsSlicef(format string, msg ...any) []error {
	return []error{NewWithSkipCallerf(1, format, msg...)}
}

func NewWithCodeAsSlice(code string, msg ...any) []error {
	return []error{NewWithSkipCallerAndCode(1, code, msg...)}
}

func NewWithMetadataAsSlice(metadata map[string]any, msg ...any) []error {
	return []error{NewWithAll(1, "", metadata, msg...)}
}

func NewWithCodeAndMetadataAsSlice(code string, metadata map[string]any, msg ...any) []error {
	return []error{NewWithAll(1, code, metadata, msg...)}
}

func NewWithCodeAsSlicef(code, format string, msg ...any) []error {
	return []error{NewWithSkipCallerAndCodef(1, code, format, msg...)}
}

func NewWithMetadataAsSlicef(metadata map[string]any, format string, msg ...any) []error {
	return []error{NewWithAllf(1, "", metadata, format, msg...)}
}

func NewWithCodeAndMetadataAsSlicef(code string, metadata map[string]any, format string, msg ...any) []error {
	return []error{NewWithAllf(1, code, metadata, format, msg...)}
}

func NewWithSkipCallerAsSlice(skipCaller int, msg ...any) []error {
	return []error{NewWithSkipCaller(skipCaller, msg...)}
}

func NewWithSkipCallerAndCodeAsSlice(skipCaller int, code string, msg ...any) []error {
	return []error{NewWithSkipCallerAndCode(skipCaller, code, msg...)}
}

func NewWithAllAsSlice(skipCaller int, code string, metadata map[string]any, msg ...any) []error {
	return []error{NewWithAll(skipCaller, code, metadata, msg...)}
}

func NewWithSkipCallerAndCodeAsSlicef(skipCaller int, code, format string, msg ...any) []error {
	return []error{NewWithSkipCallerAndCodef(skipCaller, code, format, msg...)}
}

func NewWithSkipCallerAsSlicef(skipCaller int, format string, msg ...any) []error {
	return []error{NewWithSkipCallerf(skipCaller, format, msg...)}
}

func NewWithAllAsSlicef(skipCaller int, code string, metadata map[string]any, format string, msg ...any) []error {
	return []error{NewWithAllf(skipCaller, code, metadata, format, msg...)}
}

func NewWithRawDataAsSlice(file, line, funcName, code, message string, metadata map[string]any, stack []byte) []error {
	return []error{NewWithRawData(file, line, funcName, code, message, metadata, stack)}
}

func Recovery(msg ...any) error {
	if r := recover(); r != nil {
		msg = append(msg, fmt.Sprintf("recovery: %v", r))
		return NewWithSkipCaller(1, msg...)
	}
	return nil
}

func RecoveryWithCode(code string, msg ...any) error {
	if r := recover(); r != nil {
		msg = append(msg, fmt.Sprintf("recovery: %v", r))
		return NewWithSkipCallerAndCode(1, code, msg...)
	}
	return nil
}

func RecoveryWithMetadata(metadata map[string]any, msg ...any) error {
	if r := recover(); r != nil {
		msg = append(msg, fmt.Sprintf("recovery: %v", r))
		return NewWithAll(1, "", metadata, msg...)
	}
	return nil
}

func RecoveryWithCodeAndMetadata(code string, metadata map[string]any, msg ...any) error {
	if r := recover(); r != nil {
		msg = append(msg, fmt.Sprintf("recovery: %v", r))
		return NewWithAll(1, code, metadata, msg...)
	}
	return nil
}

func Recoveryf(format string, msg ...any) error {
	if r := recover(); r != nil {
		format = format + " recovery: %v"
		msg = append(msg, r)
		return NewWithSkipCallerf(1, format, msg...)
	}
	return nil
}

func RecoveryWithCodef(code, format string, msg ...any) error {
	if r := recover(); r != nil {
		format = format + " recovery: %v"
		msg = append(msg, r)
		return NewWithSkipCallerAndCodef(1, code, format, msg...)
	}
	return nil
}

func RecoveryWithMetadataf(metadata map[string]any, format string, msg ...any) error {
	if r := recover(); r != nil {
		format = format + " recovery: %v"
		msg = append(msg, r)
		return NewWithAllf(1, "", metadata, format, msg...)
	}
	return nil
}

func RecoveryWithCodeAndMetadataf(code string, metadata map[string]any, format string, msg ...any) error {
	if r := recover(); r != nil {
		format = format + " recovery: %v"
		msg = append(msg, r)
		return NewWithAllf(1, code, metadata, format, msg...)
	}
	return nil
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

	stack := e.Stack()
	if e.parent != nil {
		stack += "\n" + inheritSep + e.parent.Error()
	}
	return e.Description() + " [STACK]: " + stack
}

func (e *Err) Raw() string {
	if e.target {
		return e.Error()
	}

	return e.Description() + " [STACK]: " + normalizeStack(e.stack)
}

func (e *Err) Description() string {
	code := ""
	if e.HasCode() {
		code = "[CODE]: " + e.code + " "
	}
	metadata := ""
	if e.HasMetadata() {
		metadata = " [METADATA]: " + toString(e.Metadata())
	}
	return code + "[CAUSE]: " + e.Cause().Error() + metadata
}

func (e *Err) PrintStack() {
	fmt.Println(e.Stack())
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

func (e *Err) HasParent() bool {
	return e.parent != nil
}

func (e *Err) Parent() *Err {
	return e.parent
}

// Chain returns the full parent chain as a slice of *Err.
// The first element is the direct parent, the last is the deepest ancestor.
// Returns nil if there is no parent.
func (e *Err) Chain() []*Err {
	var chain []*Err
	for cur := e.parent; cur != nil; cur = cur.parent {
		chain = append(chain, cur)
	}
	return chain
}

// RootCause returns the deepest ancestor in the parent chain.
// If there is no parent, returns the receiver itself.
func (e *Err) RootCause() *Err {
	cur := e
	for cur.parent != nil {
		cur = cur.parent
	}
	return cur
}

// ChainAt returns the parent at the given index in the chain.
// Index 0 is the direct parent, 1 is the grandparent, and so on.
// Panics if the index is out of range.
func (e *Err) ChainAt(index int) *Err {
	return e.Chain()[index]
}

// ChainLen returns the number of parents in the chain.
// Returns 0 if there is no parent.
func (e *Err) ChainLen() int {
	n := 0
	for cur := e.parent; cur != nil; cur = cur.parent {
		n++
	}
	return n
}

func (e *Err) Stack() string {
	stack := renderStackByPolicy(e.stack)
	if e.parent != nil {
		stack += "\n" + inheritSep + e.parent.Error()
	}
	return stack
}

func (e *Err) StackAsSlice() []string {
	stack := renderStackByPolicy(e.stack)
	lines := strings.Split(stack, "\n")

	var result []string
	for _, line := range lines {
		// Remove \t e \n, trim spaces
		cleaned := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(line, "\t", ""), "\n", ""))
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}

	// Se tem parent, adiciona as linhas do parent também
	if e.parent != nil {
		parentLines := e.parent.StackAsSlice()
		result = append(result, "---------------- [INHERITED BY]:")
		result = append(result, e.parent.Description()+" [STACK]:")
		result = append(result, parentLines...)
	}

	return result
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
