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

func NewByParentWithSkipCaller(parent error, skipCaller int, msg ...any) *Err {
	if parent == nil {
		return NewWithSkipCaller(skipCaller+1, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		parent:   Wrap(parent),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		stack:    buildDebugStack(),
	}
}

func NewByParentWithSkipCallerf(parent error, skipCaller int, format string, msg ...any) *Err {
	if parent == nil {
		return NewWithSkipCallerf(skipCaller+1, format, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		parent:   Wrap(parent),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		stack:    buildDebugStack(),
	}
}

func NewByParentWithSkipCallerAndCode(parent error, skipCaller int, code string, msg ...any) *Err {
	if parent == nil {
		return NewWithSkipCallerAndCode(skipCaller+1, code, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
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

func NewByParentWithSkipCallerAndCodef(parent error, skipCaller int, code, format string, msg ...any) *Err {
	if parent == nil {
		return NewWithSkipCallerAndCodef(skipCaller+1, code, format, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
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

func NewByParentWithAll(parent error, skipCaller int, code string, metadata map[string]any, msg ...any) *Err {
	if parent == nil {
		return NewWithAll(skipCaller+1, code, metadata, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
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

func NewByParentWithAllf(parent error, skipCaller int, code string, metadata map[string]any, format string, msg ...any) *Err {
	if parent == nil {
		return NewWithAllf(skipCaller+1, code, metadata, format, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
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

func NewByParentAsSlice(parent error, msg ...any) []error {
	return []error{NewByParent(parent, msg...)}
}

func NewByParentAsSlicef(parent error, format string, msg ...any) []error {
	return []error{NewByParentf(parent, format, msg...)}
}

func NewByParentWithCodeAsSlice(parent error, code string, msg ...any) []error {
	return []error{NewByParentWithCode(parent, code, msg...)}
}

func NewByParentWithCodeAsSlicef(parent error, code, format string, msg ...any) []error {
	return []error{NewByParentWithCodef(parent, code, format, msg...)}
}

func NewByParentWithMetadataAsSlice(parent error, metadata map[string]any, msg ...any) []error {
	return []error{NewByParentWithMetadata(parent, metadata, msg...)}
}

func NewByParentWithMetadataAsSlicef(parent error, metadata map[string]any, format string, msg ...any) []error {
	return []error{NewByParentWithMetadataf(parent, metadata, format, msg...)}
}

func NewByParentWithCodeAndMetadataAsSlice(parent error, code string, metadata map[string]any, msg ...any) []error {
	return []error{NewByParentWithCodeAndMetadata(parent, code, metadata, msg...)}
}

func NewByParentWithCodeAndMetadataAsSlicef(parent error, code string, metadata map[string]any, format string, msg ...any) []error {
	return []error{NewByParentWithCodeAndMetadataf(parent, code, metadata, format, msg...)}
}

func NewByParentWithSkipCallerAsSlice(parent error, skipCaller int, msg ...any) []error {
	return []error{NewByParentWithSkipCaller(parent, skipCaller, msg...)}
}

func NewByParentWithSkipCallerAsSlicef(parent error, skipCaller int, format string, msg ...any) []error {
	return []error{NewByParentWithSkipCallerf(parent, skipCaller, format, msg...)}
}

func NewByParentWithSkipCallerAndCodeAsSlice(parent error, skipCaller int, code string, msg ...any) []error {
	return []error{NewByParentWithSkipCallerAndCode(parent, skipCaller, code, msg...)}
}

func NewByParentWithSkipCallerAndCodeAsSlicef(parent error, skipCaller int, code, format string, msg ...any) []error {
	return []error{NewByParentWithSkipCallerAndCodef(parent, skipCaller, code, format, msg...)}
}

func NewByParentWithAllAsSlice(parent error, skipCaller int, code string, metadata map[string]any, msg ...any) []error {
	return []error{NewByParentWithAll(parent, skipCaller, code, metadata, msg...)}
}

func NewByParentWithAllAsSlicef(parent error, skipCaller int, code string, metadata map[string]any, format string, msg ...any) []error {
	return []error{NewByParentWithAllf(parent, skipCaller, code, metadata, format, msg...)}
}

func NewByChain(errs []error, msg ...any) *Err {
	if len(errs) == 0 {
		return New(msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		stack:    buildDebugStack(),
	}
}

func NewByChainf(errs []error, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return Newf(format, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		stack:    buildDebugStack(),
	}
}

func NewByChainWithCode(errs []error, code string, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithCode(code, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessage(msg...),
		stack:    buildDebugStack(),
	}
}

func NewByChainWithCodef(errs []error, code, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithCodef(code, format, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessageByFormat(format, msg...),
		stack:    buildDebugStack(),
	}
}

func NewByChainWithMetadata(errs []error, metadata map[string]any, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithMetadata(metadata, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		metadata: metadata,
		stack:    buildDebugStack(),
	}
}

func NewByChainWithMetadataf(errs []error, metadata map[string]any, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithMetadataf(metadata, format, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		metadata: metadata,
		stack:    buildDebugStack(),
	}
}

func NewByChainWithCodeAndMetadata(errs []error, code string, metadata map[string]any, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithCodeAndMetadata(code, metadata, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessage(msg...),
		metadata: metadata,
		stack:    buildDebugStack(),
	}
}

func NewByChainWithCodeAndMetadataf(errs []error, code string, metadata map[string]any, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithCodeAndMetadataf(code, metadata, format, msg...)
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessageByFormat(format, msg...),
		metadata: metadata,
		stack:    buildDebugStack(),
	}
}

func NewByChainWithSkipCaller(errs []error, skipCaller int, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithSkipCaller(skipCaller+1, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessage(msg...),
		stack:    buildDebugStack(),
	}
}

func NewByChainWithSkipCallerf(errs []error, skipCaller int, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithSkipCallerf(skipCaller+1, format, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		message:  buildMessageByFormat(format, msg...),
		stack:    buildDebugStack(),
	}
}

func NewByChainWithSkipCallerAndCode(errs []error, skipCaller int, code string, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithSkipCallerAndCode(skipCaller+1, code, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessage(msg...),
		stack:    buildDebugStack(),
	}
}

func NewByChainWithSkipCallerAndCodef(errs []error, skipCaller int, code, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithSkipCallerAndCodef(skipCaller+1, code, format, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessageByFormat(format, msg...),
		stack:    buildDebugStack(),
	}
}

func NewByChainWithAll(errs []error, skipCaller int, code string, metadata map[string]any, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithAll(skipCaller+1, code, metadata, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessage(msg...),
		metadata: metadata,
		stack:    buildDebugStack(),
	}
}

func NewByChainWithAllf(errs []error, skipCaller int, code string, metadata map[string]any, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return NewWithAllf(skipCaller+1, code, metadata, format, msg...)
	}
	file, line, funcName := callerInfos(skipCaller + 1)
	return &Err{
		parent:   inheritChainFromSlice(errs),
		file:     file,
		line:     line,
		funcName: funcName,
		code:     code,
		message:  buildMessageByFormat(format, msg...),
		metadata: metadata,
		stack:    buildDebugStack(),
	}
}

func NewByChainAsSlice(errs []error, msg ...any) []error {
	return []error{NewByChain(errs, msg...)}
}

func NewByChainAsSlicef(errs []error, format string, msg ...any) []error {
	return []error{NewByChainf(errs, format, msg...)}
}

func NewByChainWithCodeAsSlice(errs []error, code string, msg ...any) []error {
	return []error{NewByChainWithCode(errs, code, msg...)}
}

func NewByChainWithCodeAsSlicef(errs []error, code, format string, msg ...any) []error {
	return []error{NewByChainWithCodef(errs, code, format, msg...)}
}

func NewByChainWithMetadataAsSlice(errs []error, metadata map[string]any, msg ...any) []error {
	return []error{NewByChainWithMetadata(errs, metadata, msg...)}
}

func NewByChainWithMetadataAsSlicef(errs []error, metadata map[string]any, format string, msg ...any) []error {
	return []error{NewByChainWithMetadataf(errs, metadata, format, msg...)}
}

func NewByChainWithCodeAndMetadataAsSlice(errs []error, code string, metadata map[string]any, msg ...any) []error {
	return []error{NewByChainWithCodeAndMetadata(errs, code, metadata, msg...)}
}

func NewByChainWithCodeAndMetadataAsSlicef(errs []error, code string, metadata map[string]any, format string, msg ...any) []error {
	return []error{NewByChainWithCodeAndMetadataf(errs, code, metadata, format, msg...)}
}

func NewByChainWithSkipCallerAsSlice(errs []error, skipCaller int, msg ...any) []error {
	return []error{NewByChainWithSkipCaller(errs, skipCaller, msg...)}
}

func NewByChainWithSkipCallerAsSlicef(errs []error, skipCaller int, format string, msg ...any) []error {
	return []error{NewByChainWithSkipCallerf(errs, skipCaller, format, msg...)}
}

func NewByChainWithSkipCallerAndCodeAsSlice(errs []error, skipCaller int, code string, msg ...any) []error {
	return []error{NewByChainWithSkipCallerAndCode(errs, skipCaller, code, msg...)}
}

func NewByChainWithSkipCallerAndCodeAsSlicef(errs []error, skipCaller int, code, format string, msg ...any) []error {
	return []error{NewByChainWithSkipCallerAndCodef(errs, skipCaller, code, format, msg...)}
}

func NewByChainWithAllAsSlice(errs []error, skipCaller int, code string, metadata map[string]any, msg ...any) []error {
	return []error{NewByChainWithAll(errs, skipCaller, code, metadata, msg...)}
}

func NewByChainWithAllAsSlicef(errs []error, skipCaller int, code string, metadata map[string]any, format string, msg ...any) []error {
	return []error{NewByChainWithAllf(errs, skipCaller, code, metadata, format, msg...)}
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
