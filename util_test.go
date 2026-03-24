package errors

import (
	"fmt"
	"testing"
)

// --- renderStackByPolicy ---

func TestRenderStackByPolicy_Detailed(t *testing.T) {
	SetPolicy(PolicyDetailed)
	raw := []byte("goroutine 1 [running]:\nsome.func()\n\t/path/file.go:10 +0x1")
	result := renderStackByPolicy(raw)
	if result == "" {
		t.Error("expected non-empty stack")
	}
	SetPolicy(PolicyNormal)
}

func TestRenderStackByPolicy_Normal(t *testing.T) {
	SetPolicy(PolicyNormal)
	raw := []byte("goroutine 1 [running]:\nsome.func()\n\t/path/file.go:10 +0x1")
	result := renderStackByPolicy(raw)
	if result == "" {
		t.Error("expected non-empty stack")
	}
}

func TestRenderStackByPolicy_Native(t *testing.T) {
	SetPolicy(PolicyNative)
	raw := []byte("goroutine 1 [running]:\nsome.func()\n\t/path/file.go:10 +0x1")
	result := renderStackByPolicy(raw)
	_ = result // may be empty after filtering
	SetPolicy(PolicyNormal)
}

func TestRenderStackByPolicy_Default(t *testing.T) {
	// Use an invalid policy value to hit the default branch
	policy.Store(999)
	raw := []byte("some stack")
	result := renderStackByPolicy(raw)
	if result == "" {
		t.Error("expected non-empty stack for default policy")
	}
	SetPolicy(PolicyNormal)
}

// --- normalizeStack ---

func TestNormalizeStack_Empty(t *testing.T) {
	result := normalizeStack([]byte("   "))
	if result != "<empty>" {
		t.Errorf("expected '<empty>', got '%s'", result)
	}
}

func TestNormalizeStack_NonEmpty(t *testing.T) {
	result := normalizeStack([]byte("some stack trace"))
	if result != "some stack trace" {
		t.Errorf("expected 'some stack trace', got '%s'", result)
	}
}

// --- filterBoilerplateFrames ---

func TestFilterBoilerplateFrames_NoFrames(t *testing.T) {
	raw := []byte("no frames here")
	result := filterBoilerplateFrames(raw, 0)
	// fallback: returns raw
	if string(result) != "no frames here" {
		t.Errorf("expected raw fallback, got '%s'", result)
	}
}

func TestFilterBoilerplateFrames_WithBoilerplate(t *testing.T) {
	raw := []byte("runtime.goexit()\n\t/usr/local/go/src/runtime/asm.s:1650 +0x1\nmyapp.myFunc()\n\t/app/main.go:42 +0x1")
	result := filterBoilerplateFrames(raw, 0)
	_ = result
}

func TestFilterBoilerplateFrames_KeepLast(t *testing.T) {
	raw := []byte("myapp.funcA()\n\t/app/a.go:10 +0x1\nmyapp.funcB()\n\t/app/b.go:20 +0x1\nmyapp.funcC()\n\t/app/c.go:30 +0x1")
	result := filterBoilerplateFrames(raw, 1)
	_ = result
}

// --- looksLikeFileLine ---

func TestLooksLikeFileLine_True(t *testing.T) {
	if !looksLikeFileLine("\t/path/to/file.go:42 +0x1") {
		t.Error("expected true")
	}
}

func TestLooksLikeFileLine_False(t *testing.T) {
	if looksLikeFileLine("some.func()") {
		t.Error("expected false")
	}
}

// --- isBoilerplate ---

func TestIsBoilerplate_Runtime(t *testing.T) {
	if !isBoilerplate("runtime.goexit()", "\t/usr/local/go/src/runtime/asm.s:1650") {
		t.Error("expected true for runtime boilerplate")
	}
}

func TestIsBoilerplate_UserCode(t *testing.T) {
	if isBoilerplate("myapp.myFunc()", "\t/app/main.go:42") {
		t.Error("expected false for user code")
	}
}

func TestIsBoilerplate_Testing(t *testing.T) {
	if !isBoilerplate("testing.tRunner()", "\t/usr/local/go/src/testing/testing.go:1234") {
		t.Error("expected true for testing boilerplate")
	}
}

// --- escapeSpecialChars ---

func TestEscapeSpecialChars(t *testing.T) {
	input := "line1\nline2\r\nline3\ttab"
	result := escapeSpecialChars(input)
	expected := `line1\nline2\r\nline3\ttab`
	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}

// --- callerInfos ---

func TestCallerInfos(t *testing.T) {
	file, line, funcName := callerInfos(1)
	if file == "" {
		t.Error("expected non-empty file")
	}
	if line == "" {
		t.Error("expected non-empty line")
	}
	if funcName == "" {
		t.Error("expected non-empty funcName")
	}
}

// --- buildMessage ---

func TestBuildMessage_Normal(t *testing.T) {
	result := buildMessage("hello", "world")
	if result == "" {
		t.Error("expected non-empty message")
	}
}

func TestBuildMessage_Empty(t *testing.T) {
	result := buildMessage()
	if result != "<empty>" {
		t.Errorf("expected '<empty>', got '%s'", result)
	}
}

func TestBuildMessage_WithError(t *testing.T) {
	err := New("inner error")
	result := buildMessage(err)
	if result == "" {
		t.Error("expected non-empty message from error")
	}
}

// --- buildMessageByFormat ---

func TestBuildMessageByFormat(t *testing.T) {
	result := buildMessageByFormat("value is %s", "42")
	if result != "value is 42" {
		t.Errorf("expected 'value is 42', got '%s'", result)
	}
}

// --- cleanMessage ---

func TestCleanMessage_RemovesKeywords(t *testing.T) {
	input := "error [CODE] [METADATA] [STACK] [CAUSE] here"
	result := cleanMessage(input)
	for _, kw := range []string{"[CODE]", "[METADATA]", "[STACK]", "[CAUSE]"} {
		if contains(result, kw) {
			t.Errorf("expected keyword '%s' to be removed", kw)
		}
	}
}

func TestCleanMessage_ReplacesNewlines(t *testing.T) {
	result := cleanMessage("line1\nline2")
	if contains(result, "\n") {
		t.Error("expected newlines to be replaced")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// --- filterMsg ---

func TestFilterMsg_WithString(t *testing.T) {
	result := filterMsg("hello")
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
}

func TestFilterMsg_WithError(t *testing.T) {
	err := New("filter error")
	result := filterMsg(err)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
}

func TestFilterMsg_WithNilError(t *testing.T) {
	var err error
	result := filterMsg(err)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
}

// --- formatFuncName ---

func TestFormatFuncName(t *testing.T) {
	result := formatFuncName("github.com/user/pkg.MyFunc")
	if result != "MyFunc" {
		t.Errorf("expected 'MyFunc', got '%s'", result)
	}
}

// --- toString / toStringWithErr ---

func TestToString_String(t *testing.T) {
	result := toString("hello")
	if result != "hello" {
		t.Errorf("expected 'hello', got '%s'", result)
	}
}

func TestToString_Int(t *testing.T) {
	result := toString(42)
	if result != "42" {
		t.Errorf("expected '42', got '%s'", result)
	}
}

func TestToString_Float(t *testing.T) {
	result := toString(3.14)
	if result == "" {
		t.Error("expected non-empty float string")
	}
}

func TestToString_Bool(t *testing.T) {
	result := toString(true)
	if result != "true" {
		t.Errorf("expected 'true', got '%s'", result)
	}
}

func TestToString_Bytes(t *testing.T) {
	result := toString([]byte("bytes"))
	if result != "bytes" {
		t.Errorf("expected 'bytes', got '%s'", result)
	}
}

func TestToString_Map(t *testing.T) {
	result := toString(map[string]int{"a": 1})
	if result == "" {
		t.Error("expected non-empty map string")
	}
}

func TestToString_Struct(t *testing.T) {
	type S struct{ X int }
	result := toString(S{X: 1})
	if result == "" {
		t.Error("expected non-empty struct string")
	}
}

func TestToString_Nil(t *testing.T) {
	result := toString(nil)
	_ = result // should not panic
}

func TestToString_Stringer(t *testing.T) {
	type myStringer struct{}
	// use fmt.Stringer via a known type
	result := toString(fmt.Errorf("stringer error"))
	if result != "stringer error" {
		t.Errorf("expected 'stringer error', got '%s'", result)
	}
}

func TestToString_NilPointer(t *testing.T) {
	var p *int
	result := toString(p)
	_ = result // should not panic
}

func TestToString_IntVariants(t *testing.T) {
	cases := []any{int8(1), int16(2), int32(3), int64(4)}
	for _, c := range cases {
		result := toString(c)
		if result == "" {
			t.Errorf("expected non-empty string for %T", c)
		}
	}
}

func TestToString_UintVariants(t *testing.T) {
	cases := []any{uint(1), uint8(2), uint16(3), uint32(4), uint64(5), uintptr(6)}
	for _, c := range cases {
		result := toString(c)
		if result == "" {
			t.Errorf("expected non-empty string for %T", c)
		}
	}
}

func TestToString_Complex(t *testing.T) {
	result := toString(complex(1, 2))
	if result == "" {
		t.Error("expected non-empty complex string")
	}
}

func TestToString_Array(t *testing.T) {
	result := toString([3]int{1, 2, 3})
	if result == "" {
		t.Error("expected non-empty array string")
	}
}

func TestToString_ByteArray(t *testing.T) {
	result := toString([3]byte{'a', 'b', 'c'})
	if result != "abc" {
		t.Errorf("expected 'abc', got '%s'", result)
	}
}

func TestToStringWithErr_Nil(t *testing.T) {
	_, err := toStringWithErr(nil)
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestToStringWithErr_NilPointer(t *testing.T) {
	var p *int
	_, err := toStringWithErr(p)
	if err == nil {
		t.Error("expected error for nil pointer")
	}
}

func TestToStringWithErr_UnsupportedType(t *testing.T) {
	// channel is unsupported
	ch := make(chan int)
	_, err := toStringWithErr(ch)
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

// --- implementsStringer / implementsError ---

func TestImplementsStringer_Nil(t *testing.T) {
	if implementsStringer(nil) {
		t.Error("expected false for nil type")
	}
}

func TestImplementsError_Nil(t *testing.T) {
	if implementsError(nil) {
		t.Error("expected false for nil type")
	}
}

// --- optionalField ---

func TestOptionalField_Empty(t *testing.T) {
	result := optionalField("label: ", "")
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestOptionalField_NonEmpty(t *testing.T) {
	result := optionalField("label: ", "value")
	if result != "label: value" {
		t.Errorf("expected 'label: value', got '%s'", result)
	}
}
