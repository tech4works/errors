package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const regex = `^(?:\[CODE]: (.+?) )?` + // 1: code (optional)
	`\[CAUSE]: \(([^:]+):(\d+)\) ([^:]+): (.+?)` + // 2: file, 3: line, 4: func, 5: msg
	`(?: \[METADATA]: (.+?))?` + // 6: metadata (optional)
	` \[STACK]:\s*([\s\S]+)` // 7: stack

const snapshotRegex = `^\[CODE]: (.+?) \[MESSAGE]: (.+)$`

const inheritSep = "----------------\n\t\t|\t[INHERITED BY]: "

var (
	compiledRegex         = regexp.MustCompile(regex)
	compiledSnapshotRegex = regexp.MustCompile(snapshotRegex)
)

func Is(err, target error) bool {
	if err == nil || target == nil {
		return false
	}

	var mErr, mTarget *Err

	errString := err.Error()
	if As(err, &mErr) {
		if mErr.IsTarget() {
			return false
		}
		errString = mErr.Error()
	}

	targetString := target.Error()
	targetContains := targetString
	if As(target, &mTarget) {
		targetString = mTarget.Snapshot()
		if mTarget.HasCode() {
			targetContains = mTarget.Code()
		} else {
			targetContains = mTarget.Message()
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
	}

	if t, ok := target.(**Err); ok {
		extracted, parsed := extract(err, 3)
		if parsed == nil {
			return false
		}
		if extracted {
			*t = parsed
		}
		return extracted
	}

	if t, ok := target.(*Err); ok {
		extracted, parsed := extract(err, 3)
		if extracted && parsed != nil {
			*t = *parsed
		}
		return extracted
	}

	if isValidAsTarget(target) {
		return errors.As(err, target)
	}

	return false
}

func Wrap(err error) *Err {
	if err == nil {
		return nil
	}
	_, e := extract(err, 3)
	return e
}

func WrapWithExtracted(err error) (bool, *Err) {
	if err == nil {
		return false, nil
	}
	return extract(err, 3)
}

func Parse(v any) (*Err, error) {
	if v == nil {
		return nil, nil
	}
	s, err := toStringWithErr(v)
	if err != nil {
		return nil, err
	}
	_, e := extract(errors.New(s), 3)
	return e, nil
}

func Inherit(err error, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithSkipCaller(2, msg...)

	extracted, parent := extract(err, 2)

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

	e := NewWithSkipCallerf(2, format, msg...)

	extracted, parent := extract(err, 2)

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

	extracted, parent := extract(err, 2)

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

	e := NewWithSkipCallerf(skipCaller+1, format, msg...)

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

func InheritWithAll(err error, skipCaller int, code string, metadata map[string]any, msg ...any) *Err {
	if err == nil {
		return nil
	}

	e := NewWithAll(skipCaller+1, code, metadata, msg...)

	extracted, parent := extract(err, skipCaller+1)

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

	e := NewWithSkipCallerAndCode(skipCaller+1, code, msg...)

	extracted, parent := extract(err, skipCaller+1)

	e.parent = parent
	if e.Message() == "<empty>" {
		e.message = parent.message
	}

	if extracted {
		e.metadata = parent.metadata
	}

	return e
}

func InheritAsSlice(err error, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{Inherit(err, msg...)}
}

func InheritAsSlicef(err error, format string, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{Inheritf(err, format, msg...)}
}

func InheritWithSkipCallerAsSlice(err error, skipCaller int, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{InheritWithSkipCaller(err, skipCaller, msg...)}
}

func InheritWithCodeAsSlice(err error, code string, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{InheritWithCode(err, code, msg...)}
}

func InheritWithSkipCallerAsSlicef(err error, skipCaller int, format string, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{InheritWithSkipCallerf(err, skipCaller, format, msg...)}
}

func InheritWithAllAsSlice(err error, skipCaller int, code string, metadata map[string]any, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{InheritWithAll(err, skipCaller, code, metadata, msg...)}
}

func InheritWithSkipCallerAndCodeAsSlice(err error, skipCaller int, code string, msg ...any) []error {
	if err == nil {
		return nil
	}
	return []error{InheritWithSkipCallerAndCode(err, skipCaller, code, msg...)}
}

// Chain builds an ordered parent chain from a variadic list of errors.
// The first error is the root; each subsequent error becomes a deeper parent.
// Returns nil if no errors are provided.
func Chain(errs ...error) *Err {
	return inheritChainFromSlice(errs)
}

// ChainFromSlice builds an ordered parent chain from a slice of errors.
// The first error is the root; each subsequent error becomes a deeper parent.
// Returns nil if the slice is empty.
func ChainFromSlice(errs []error) *Err {
	return inheritChainFromSlice(errs)
}

// InheritFromSlice builds an ordered parent chain from a slice of errors.
// When msg is provided, a new root *Err wraps the chain with that message.
// When msg is omitted, the chain itself is returned directly — the first error in the slice is the root.
// Returns nil if the slice is empty.
func InheritFromSlice(errs []error, msg ...any) *Err {
	if len(errs) == 0 {
		return nil
	}
	chain := inheritChainFromSlice(errs)
	if buildMessage(msg...) == "<empty>" {
		return chain
	}
	e := NewWithSkipCaller(2, msg...)
	e.parent = chain
	return e
}

// InheritFromSlicef is the format-string variant of InheritFromSlice.
func InheritFromSlicef(errs []error, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return nil
	}
	e := NewWithSkipCallerf(2, format, msg...)
	e.parent = inheritChainFromSlice(errs)
	return e
}

// InheritFromSliceWithCode is the code variant of InheritFromSlice.
// When msg is omitted, the chain is returned with the code applied to the root node.
func InheritFromSliceWithCode(errs []error, code string, msg ...any) *Err {
	if len(errs) == 0 {
		return nil
	}
	chain := inheritChainFromSlice(errs)
	if buildMessage(msg...) == "<empty>" {
		chain.code = code
		return chain
	}
	e := NewWithSkipCallerAndCode(2, code, msg...)
	e.parent = chain
	return e
}

// InheritFromSliceWithCodef is the code and format-string variant of InheritFromSlice.
func InheritFromSliceWithCodef(errs []error, code, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return nil
	}
	e := NewWithSkipCallerAndCodef(2, code, format, msg...)
	e.parent = inheritChainFromSlice(errs)
	return e
}

// InheritFromSliceWithAll is the skipCaller, code and metadata variant of InheritFromSlice.
// When msg is omitted, the chain is returned with code and metadata applied to the root node.
func InheritFromSliceWithAll(errs []error, skipCaller int, code string, metadata map[string]any, msg ...any) *Err {
	if len(errs) == 0 {
		return nil
	}
	chain := inheritChainFromSlice(errs)
	if buildMessage(msg...) == "<empty>" {
		chain.code = code
		chain.metadata = metadata
		return chain
	}
	e := NewWithAll(skipCaller+1, code, metadata, msg...)
	e.parent = chain
	return e
}

// InheritFromSliceWithAllf is the skipCaller, code, metadata and format-string variant of InheritFromSlice.
func InheritFromSliceWithAllf(errs []error, skipCaller int, code string, metadata map[string]any, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return nil
	}
	e := NewWithAllf(skipCaller+1, code, metadata, format, msg...)
	e.parent = inheritChainFromSlice(errs)
	return e
}

// InheritFromSliceWithSkipCaller is the custom skipCaller variant of InheritFromSlice.
// When msg is omitted, the chain itself is returned directly.
func InheritFromSliceWithSkipCaller(errs []error, skipCaller int, msg ...any) *Err {
	if len(errs) == 0 {
		return nil
	}
	chain := inheritChainFromSlice(errs)
	if buildMessage(msg...) == "<empty>" {
		return chain
	}
	e := NewWithSkipCaller(skipCaller+1, msg...)
	e.parent = chain
	return e
}

// InheritFromSliceWithSkipCallerf is the skipCaller and format-string variant of InheritFromSlice.
func InheritFromSliceWithSkipCallerf(errs []error, skipCaller int, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return nil
	}
	e := NewWithSkipCallerf(skipCaller+1, format, msg...)
	e.parent = inheritChainFromSlice(errs)
	return e
}

// InheritFromSliceWithSkipCallerAndCode is the skipCaller and code variant of InheritFromSlice.
// When msg is omitted, the chain is returned with the code applied to the root node.
func InheritFromSliceWithSkipCallerAndCode(errs []error, skipCaller int, code string, msg ...any) *Err {
	if len(errs) == 0 {
		return nil
	}
	chain := inheritChainFromSlice(errs)
	if buildMessage(msg...) == "<empty>" {
		chain.code = code
		return chain
	}
	e := NewWithSkipCallerAndCode(skipCaller+1, code, msg...)
	e.parent = chain
	return e
}

// InheritFromSliceWithSkipCallerAndCodef is the skipCaller, code and format-string variant of InheritFromSlice.
func InheritFromSliceWithSkipCallerAndCodef(errs []error, skipCaller int, code, format string, msg ...any) *Err {
	if len(errs) == 0 {
		return nil
	}
	e := NewWithSkipCallerAndCodef(skipCaller+1, code, format, msg...)
	e.parent = inheritChainFromSlice(errs)
	return e
}

func Join(errs []error, sep string) *Err {
	if len(errs) == 0 {
		return nil
	}
	file, line, funcName := callerInfos(1)
	return &Err{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  JoinToString(errs, sep),
		stack:    buildDebugStack(),
	}
}

func JoinInherit(errs []error, sep string, msg ...any) error {
	if errs == nil {
		return nil
	}
	return Inherit(Join(errs, sep), msg...)
}

func JoinInheritf(errs []error, sep, format string, msg ...any) *Err {
	if errs == nil {
		return nil
	}
	return Inheritf(Join(errs, sep), format, msg...)
}

func JoinInheritAsSlice(errs []error, sep string, msg ...any) []error {
	if errs == nil {
		return nil
	}
	return InheritAsSlice(Join(errs, sep), msg...)
}

func JoinInheritAsSlicef(errs []error, sep, format string, msg ...any) []error {
	if errs == nil {
		return nil
	}
	return InheritAsSlicef(Join(errs, sep), format, msg...)
}

func JoinInheritWithSkipCaller(errs []error, sep string, skipCaller int, msg ...any) *Err {
	if errs == nil {
		return nil
	}

	return InheritWithSkipCaller(Join(errs, sep), skipCaller, msg...)
}

func JoinInheritWithCode(errs []error, sep, code string, msg ...any) *Err {
	if errs == nil {
		return nil
	}

	return InheritWithCode(Join(errs, sep), code, msg...)
}

func JoinInheritWithSkipCallerf(errs []error, sep string, skipCaller int, format string, msg ...any) *Err {
	if errs == nil {
		return nil
	}

	return InheritWithSkipCallerf(Join(errs, sep), skipCaller, format, msg...)
}

func JoinInheritWithAll(errs []error, sep string, skipCaller int, code string, metadata map[string]any, msg ...any) *Err {
	if errs == nil {
		return nil
	}

	return InheritWithAll(Join(errs, sep), skipCaller, code, metadata, msg...)
}

func JoinInheritWithSkipCallerAndCode(errs []error, sep string, skipCaller int, code string, msg ...any) *Err {
	if errs == nil {
		return nil
	}

	return InheritWithSkipCallerAndCode(Join(errs, sep), skipCaller, code, msg...)
}

// JoinToString concatena as mensagens dos erros da fatia separadas por sep.
// Erros nil ou sem mensagem são ignorados silenciosamente.
func JoinToString(errs []error, sep string) (result string) {
	if len(errs) == 0 {
		return ""
	}
	for i, err := range errs {
		dt := Wrap(err)
		if dt != nil {
			result += dt.message
		}
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

	rg := compiledRegex
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

	stackRaw := matches[7]

	// Separate the own stack from any [INHERITED BY] blocks
	ownStack, inheritedParts := splitInheritedBlocks(stackRaw)

	e = &Err{
		file:     matches[2],
		line:     matches[3],
		funcName: matches[4],
		message:  matches[5],
		stack:    []byte(ownStack),
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

	// Try to extract a nested error from the message (snapshot format: [CODE]: X [MESSAGE]: Y)
	if nested := parseSnapshotFromMessage(e.message); nested != nil {
		e.message = nested.message
		if nested.code != "" && e.code == "" {
			e.code = nested.code
		}
		if nested.code != "" && e.code != nested.code {
			// The message contained a different code — preserve it as a parent
			nested.file = e.file
			nested.line = e.line
			nested.funcName = e.funcName
			nested.stack = e.stack
			e.parent = chainParents(nested, inheritedParts)
			return true, e
		}
	}

	// Build parent chain from [INHERITED BY] blocks
	e.parent = chainParents(nil, inheritedParts)

	return true, e
}

// parseSnapshotFromMessage checks if a message contains a snapshot pattern
// like "[CODE]: X [MESSAGE]: Y" and extracts it into an Err with code and message.
func parseSnapshotFromMessage(msg string) *Err {
	m := compiledSnapshotRegex.FindStringSubmatch(msg)
	if len(m) == 0 {
		return nil
	}
	return &Err{
		code:    m[1],
		message: m[2],
	}
}

// splitInheritedBlocks splits a stack string at each [INHERITED BY] separator,
// returning the own stack (before the first separator) and a slice of the
// inherited error strings (each one is a full Error() representation).
func splitInheritedBlocks(stack string) (string, []string) {
	sep := "----------------\n\t\t|\t[INHERITED BY]: "
	parts := strings.SplitN(stack, sep, 2)
	if len(parts) == 1 {
		return stack, nil
	}

	ownStack := parts[0]
	rest := parts[1]

	var inherited []string
	for {
		idx := strings.Index(rest, sep)
		if idx < 0 {
			inherited = append(inherited, rest)
			break
		}
		inherited = append(inherited, rest[:idx])
		rest = rest[idx+len(sep):]
	}

	return ownStack, inherited
}

// chainParents builds a linked parent chain from an optional head Err and
// a slice of inherited error strings. Each string is parsed via extract.
func chainParents(head *Err, inheritedParts []string) *Err {
	// Parse all inherited parts into Err nodes
	var nodes []*Err
	if head != nil {
		nodes = append(nodes, head)
	}
	for _, part := range inheritedParts {
		_, parsed := extract(errors.New(part), 3)
		if parsed != nil {
			nodes = append(nodes, parsed)
		}
	}

	if len(nodes) == 0 {
		return nil
	}

	// Link them: first is the direct parent, last is the deepest ancestor
	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].parent = nodes[i+1]
	}

	return nodes[0]
}

func isValidAsTarget(target any) bool {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false
	}
	elem := rv.Type().Elem()
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	return elem.Implements(errorType) || elem.Kind() == reflect.Interface
}

// inheritChainFromSlice builds a parent chain from a slice of errors.
// The first error in the slice sits at the top of the chain; the last is the deepest parent.
// Returns nil if the slice is empty.
func inheritChainFromSlice(errs []error) *Err {
	if len(errs) == 0 {
		return nil
	}
	_, chain := extract(errs[len(errs)-1], 3)
	for i := len(errs) - 2; i >= 0; i-- {
		_, node := extract(errs[i], 3)
		node.parent = chain
		chain = node
	}
	return chain
}
