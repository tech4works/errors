package errors

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("something went wrong")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Message() != "something went wrong" {
		t.Errorf("expected message 'something went wrong', got '%s'", err.Message())
	}
	if err.File() == "" {
		t.Error("expected non-empty file")
	}
	if err.Line() == 0 {
		t.Error("expected non-zero line")
	}
}

func TestNewEmpty(t *testing.T) {
	err := New()
	if err.Message() != "<empty>" {
		t.Errorf("expected '<empty>', got '%s'", err.Message())
	}
}

func TestNewf(t *testing.T) {
	err := Newf("error %s", "42")
	if err.Message() != "error 42" {
		t.Errorf("expected 'error 42', got '%s'", err.Message())
	}
}

func TestNewWithCode(t *testing.T) {
	err := NewWithCode("ERR001", "bad request")
	if err.Code() != "ERR001" {
		t.Errorf("expected code 'ERR001', got '%s'", err.Code())
	}
	if !err.HasCode() {
		t.Error("expected HasCode to be true")
	}
}

func TestNewWithMetadata(t *testing.T) {
	meta := map[string]any{"key": "value"}
	err := NewWithMetadata(meta, "with meta")
	if !err.HasMetadata() {
		t.Error("expected HasMetadata to be true")
	}
	if err.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", err.Metadata()["key"])
	}
}

func TestNewWithCodeAndMetadata(t *testing.T) {
	meta := map[string]any{"x": 1}
	err := NewWithCodeAndMetadata("CODE", meta, "msg")
	if err.Code() != "CODE" {
		t.Errorf("expected code 'CODE', got '%s'", err.Code())
	}
	if err.Metadata()["x"] != 1 {
		t.Errorf("expected metadata x=1")
	}
}

func TestNewWithCodef(t *testing.T) {
	err := NewWithCodef("ERR", "value is %s", "bad")
	if err.Code() != "ERR" {
		t.Errorf("expected code 'ERR', got '%s'", err.Code())
	}
	if err.Message() != "value is bad" {
		t.Errorf("expected 'value is bad', got '%s'", err.Message())
	}
}

func TestNewWithMetadataf(t *testing.T) {
	meta := map[string]any{"k": "v"}
	err := NewWithMetadataf(meta, "msg %s", "here")
	if err.Message() != "msg here" {
		t.Errorf("expected 'msg here', got '%s'", err.Message())
	}
}

func TestNewWithCodeAndMetadataf(t *testing.T) {
	meta := map[string]any{"k": "v"}
	err := NewWithCodeAndMetadataf("C1", meta, "msg %s", "x")
	if err.Code() != "C1" {
		t.Errorf("expected code 'C1', got '%s'", err.Code())
	}
}

func TestNewWithSkipCaller(t *testing.T) {
	err := NewWithSkipCaller(1, "skip caller")
	if err.Message() != "skip caller" {
		t.Errorf("expected 'skip caller', got '%s'", err.Message())
	}
}

func TestNewWithSkipCallerAndCode(t *testing.T) {
	err := NewWithSkipCallerAndCode(1, "SC01", "msg")
	if err.Code() != "SC01" {
		t.Errorf("expected 'SC01', got '%s'", err.Code())
	}
}

func TestNewWithAll(t *testing.T) {
	meta := map[string]any{"a": "b"}
	err := NewWithAll(1, "ALL01", meta, "all msg")
	if err.Code() != "ALL01" {
		t.Errorf("expected 'ALL01', got '%s'", err.Code())
	}
	if err.Metadata()["a"] != "b" {
		t.Error("expected metadata a=b")
	}
}

func TestNewWithSkipCallerAndCodef(t *testing.T) {
	err := NewWithSkipCallerAndCodef(1, "FMT01", "fmt %s", "val")
	if err.Code() != "FMT01" {
		t.Errorf("expected 'FMT01', got '%s'", err.Code())
	}
	if err.Message() != "fmt val" {
		t.Errorf("expected 'fmt val', got '%s'", err.Message())
	}
}

func TestNewWithSkipCallerf(t *testing.T) {
	err := NewWithSkipCallerf(1, "fmt %s", "99")
	if err.Message() != "fmt 99" {
		t.Errorf("expected 'fmt 99', got '%s'", err.Message())
	}
}

func TestNewWithAllf(t *testing.T) {
	meta := map[string]any{"z": "w"}
	err := NewWithAllf(1, "ALLF", meta, "msg %s", "x")
	if err.Code() != "ALLF" {
		t.Errorf("expected 'ALLF', got '%s'", err.Code())
	}
}

func TestNewWithRawData(t *testing.T) {
	err := NewWithRawData("file.go", "10", "myFunc", "RAW01", "raw msg", nil, []byte("stack"))
	if err.File() != "file.go" {
		t.Errorf("expected 'file.go', got '%s'", err.File())
	}
	if err.Line() != 10 {
		t.Errorf("expected line 10, got %d", err.Line())
	}
	if err.Func() != "myFunc" {
		t.Errorf("expected 'myFunc', got '%s'", err.Func())
	}
	if err.Code() != "RAW01" {
		t.Errorf("expected 'RAW01', got '%s'", err.Code())
	}
	if err.Message() != "raw msg" {
		t.Errorf("expected 'raw msg', got '%s'", err.Message())
	}
}

func TestNewAsSlice(t *testing.T) {
	errs := NewAsSlice("slice error")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewAsSlicef(t *testing.T) {
	errs := NewAsSlicef("error %d", 1)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithCodeAsSlice(t *testing.T) {
	errs := NewWithCodeAsSlice("C1", "msg")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithMetadataAsSlice(t *testing.T) {
	errs := NewWithMetadataAsSlice(map[string]any{"k": "v"}, "msg")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithCodeAndMetadataAsSlice(t *testing.T) {
	errs := NewWithCodeAndMetadataAsSlice("C1", map[string]any{"k": "v"}, "msg")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithCodeAsSlicef(t *testing.T) {
	errs := NewWithCodeAsSlicef("C1", "msg %s", "x")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithMetadataAsSlicef(t *testing.T) {
	errs := NewWithMetadataAsSlicef(map[string]any{"k": "v"}, "msg %s", "x")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithCodeAndMetadataAsSlicef(t *testing.T) {
	errs := NewWithCodeAndMetadataAsSlicef("C1", map[string]any{"k": "v"}, "msg %s", "x")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithSkipCallerAsSlice(t *testing.T) {
	errs := NewWithSkipCallerAsSlice(1, "msg")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithSkipCallerAndCodeAsSlice(t *testing.T) {
	errs := NewWithSkipCallerAndCodeAsSlice(1, "C1", "msg")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithAllAsSlice(t *testing.T) {
	errs := NewWithAllAsSlice(1, "C1", map[string]any{"k": "v"}, "msg")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithSkipCallerAndCodeAsSlicef(t *testing.T) {
	errs := NewWithSkipCallerAndCodeAsSlicef(1, "C1", "msg %s", "x")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithSkipCallerAsSlicef(t *testing.T) {
	errs := NewWithSkipCallerAsSlicef(1, "msg %s", "x")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithAllAsSlicef(t *testing.T) {
	errs := NewWithAllAsSlicef(1, "C1", map[string]any{"k": "v"}, "msg %s", "x")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithRawDataAsSlice(t *testing.T) {
	errs := NewWithRawDataAsSlice("f.go", "5", "fn", "C1", "msg", nil, []byte("stack"))
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestTargetWithCode(t *testing.T) {
	target := TargetWithCode("T01")
	if !target.IsTarget() {
		t.Error("expected IsTarget to be true")
	}
	if target.Code() != "T01" {
		t.Errorf("expected code 'T01', got '%s'", target.Code())
	}
}

func TestTargetWithMessage(t *testing.T) {
	target := TargetWithMessage("target msg")
	if !target.IsTarget() {
		t.Error("expected IsTarget to be true")
	}
	if target.Message() != "target msg" {
		t.Errorf("expected 'target msg', got '%s'", target.Message())
	}
}

func TestErrError(t *testing.T) {
	err := New("test error")
	s := err.Error()
	if s == "" {
		t.Error("expected non-empty Error() string")
	}
}

func TestErrStackAsSlice(t *testing.T) {
	err := New("test error")
	stack := err.StackAsSlice()
	if len(stack) == 0 {
		t.Error("expected non-empty StackAsSlice()")
	}
	// Verifica que não contém \t ou \n
	for _, line := range stack {
		if strings.Contains(line, "\t") || strings.Contains(line, "\n") {
			t.Errorf("expected StackAsSlice() line to not contain tabs or newlines, got: %s", line)
		}
	}
}

func TestErrErrorWithCode(t *testing.T) {
	err := NewWithCode("C1", "msg")
	s := err.Error()
	if s == "" {
		t.Error("expected non-empty Error() string")
	}
}

func TestErrErrorWithMetadata(t *testing.T) {
	err := NewWithMetadata(map[string]any{"k": "v"}, "msg")
	s := err.Error()
	if s == "" {
		t.Error("expected non-empty Error() string")
	}
}

func TestErrTargetError(t *testing.T) {
	target := TargetWithCode("T01")
	s := target.Error()
	if s == "" {
		t.Error("expected non-empty Error() string for target")
	}
}

func TestErrTargetErrorWithMessage(t *testing.T) {
	target := TargetWithMessage("msg")
	s := target.Error()
	if s == "" {
		t.Error("expected non-empty Error() string for target with message")
	}
}

func TestErrCause(t *testing.T) {
	err := New("cause test")
	cause := err.Cause()
	if cause == nil {
		t.Error("expected non-nil cause")
	}
}

func TestErrSnapshot(t *testing.T) {
	err := New("snap")
	if err.Snapshot() != "snap" {
		t.Errorf("expected 'snap', got '%s'", err.Snapshot())
	}
}

func TestErrSnapshotWithCode(t *testing.T) {
	err := NewWithCode("C1", "snap")
	snap := err.Snapshot()
	if snap == "" {
		t.Error("expected non-empty snapshot")
	}
}

func TestErrString(t *testing.T) {
	err := New("string test")
	if err.String() != err.Error() {
		t.Error("String() should equal Error()")
	}
}

func TestErrRaw(t *testing.T) {
	err := New("raw test")
	raw := err.Raw()
	if raw == "" {
		t.Error("expected non-empty Raw()")
	}
}

func TestErrRawWithCode(t *testing.T) {
	err := NewWithCode("C1", "raw test")
	raw := err.Raw()
	if raw == "" {
		t.Error("expected non-empty Raw() with code")
	}
}

func TestErrRawWithMetadata(t *testing.T) {
	err := NewWithMetadata(map[string]any{"k": "v"}, "raw test")
	raw := err.Raw()
	if raw == "" {
		t.Error("expected non-empty Raw() with metadata")
	}
}

func TestErrStack(t *testing.T) {
	err := New("stack test")
	stack := err.Stack()
	if stack == "" {
		t.Error("expected non-empty Stack()")
	}
}

func TestErrStackFormatted(t *testing.T) {
	err := New("stack test")
	stack := err.Stack()
	if stack == "" {
		t.Error("expected non-empty Stack()")
	}
}

func TestErrHasMessage(t *testing.T) {
	err := New("msg")
	if !err.HasMessage() {
		t.Error("expected HasMessage to be true")
	}
}

func TestSetPolicy(t *testing.T) {
	SetPolicy(PolicyDetailed)
	if getPolicy() != PolicyDetailed {
		t.Error("expected PolicyDetailed")
	}
	SetPolicy(PolicyNormal)
	if getPolicy() != PolicyNormal {
		t.Error("expected PolicyNormal")
	}
	SetPolicy(PolicyNative)
	if getPolicy() != PolicyNative {
		t.Error("expected PolicyNative")
	}
	SetPolicy(PolicyNormal) // restore
}

func TestErrWithParent(t *testing.T) {
	parent := New("parent error")
	child := Inherit(parent, "child error")
	s := child.Error()
	if s == "" {
		t.Error("expected non-empty error with parent")
	}
}

func TestErrWithParentFormatted(t *testing.T) {
	parent := New("parent error")
	child := Inherit(parent, "child error")
	s := child.Error()
	if s == "" {
		t.Error("expected non-empty error with parent")
	}
}

func TestNewByParent(t *testing.T) {
	parent := New("parent error")
	child := NewByParent(parent, "child error")
	if child.parent != parent {
		t.Error("expected parent to be set")
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if child.HasCode() {
		t.Error("expected no code to be inherited")
	}
	if child.HasMetadata() {
		t.Error("expected no metadata to be inherited")
	}
}

func TestNewByParentf(t *testing.T) {
	parent := New("parent error")
	child := NewByParentf(parent, "child %s", "error")
	if child.parent != parent {
		t.Error("expected parent to be set")
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if child.HasCode() {
		t.Error("expected no code to be inherited")
	}
	if child.HasMetadata() {
		t.Error("expected no metadata to be inherited")
	}
}

func TestNewByParentWithCode(t *testing.T) {
	parent := New("parent error")
	child := NewByParentWithCode(parent, "ERR_CODE", "child error")
	if child.parent != parent {
		t.Error("expected parent to be set")
	}
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if child.HasMetadata() {
		t.Error("expected no metadata to be inherited")
	}
}

func TestNewByParentWithCodef(t *testing.T) {
	parent := New("parent error")
	child := NewByParentWithCodef(parent, "ERR_CODE", "child %s", "error")
	if child.parent != parent {
		t.Error("expected parent to be set")
	}
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if child.HasMetadata() {
		t.Error("expected no metadata to be inherited")
	}
}

func TestNewByParentWithMetadata(t *testing.T) {
	parent := New("parent error")
	metadata := map[string]any{"key": "value"}
	child := NewByParentWithMetadata(parent, metadata, "child error")
	if child.parent != parent {
		t.Error("expected parent to be set")
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", child.Metadata()["key"])
	}
	if child.HasCode() {
		t.Error("expected no code to be inherited")
	}
}

func TestNewByParentWithMetadataf(t *testing.T) {
	parent := New("parent error")
	metadata := map[string]any{"key": "value"}
	child := NewByParentWithMetadataf(parent, metadata, "child %s", "error")
	if child.parent != parent {
		t.Error("expected parent to be set")
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", child.Metadata()["key"])
	}
	if child.HasCode() {
		t.Error("expected no code to be inherited")
	}
}

func TestNewByParentWithCodeAndMetadata(t *testing.T) {
	parent := New("parent error")
	metadata := map[string]any{"key": "value"}
	child := NewByParentWithCodeAndMetadata(parent, "ERR_CODE", metadata, "child error")
	if child.parent != parent {
		t.Error("expected parent to be set")
	}
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", child.Metadata()["key"])
	}
}

func TestNewByParentWithCodeAndMetadataf(t *testing.T) {
	parent := New("parent error")
	metadata := map[string]any{"key": "value"}
	child := NewByParentWithCodeAndMetadataf(parent, "ERR_CODE", metadata, "child %s", "error")
	if child.parent != parent {
		t.Error("expected parent to be set")
	}
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", child.Metadata()["key"])
	}
}

func TestNewByChain(t *testing.T) {
	errs := []error{New("err1"), New("err2"), New("err3")}
	child := NewByChain(errs, "child error")
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasParent() {
		t.Error("expected parent chain to be set")
	}
	if child.ChainLen() != 3 {
		t.Errorf("expected chain length 3, got %d", child.ChainLen())
	}
	if child.HasCode() {
		t.Error("expected no code")
	}
	if child.HasMetadata() {
		t.Error("expected no metadata")
	}
}

func TestNewByChain_EmptySlice(t *testing.T) {
	child := NewByChain(nil, "child error")
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if child.HasParent() {
		t.Error("expected no parent when slice is empty")
	}
}

func TestNewByChainf(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	child := NewByChainf(errs, "child %s", "error")
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasParent() {
		t.Error("expected parent chain to be set")
	}
	if child.ChainLen() != 2 {
		t.Errorf("expected chain length 2, got %d", child.ChainLen())
	}
	if child.HasCode() {
		t.Error("expected no code")
	}
	if child.HasMetadata() {
		t.Error("expected no metadata")
	}
}

func TestNewByChainf_EmptySlice(t *testing.T) {
	child := NewByChainf(nil, "child %s", "error")
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if child.HasParent() {
		t.Error("expected no parent when slice is empty")
	}
}

func TestNewByChainWithCode(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	child := NewByChainWithCode(errs, "ERR_CODE", "child error")
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasParent() {
		t.Error("expected parent chain to be set")
	}
	if child.ChainLen() != 2 {
		t.Errorf("expected chain length 2, got %d", child.ChainLen())
	}
	if child.HasMetadata() {
		t.Error("expected no metadata")
	}
}

func TestNewByChainWithCode_EmptySlice(t *testing.T) {
	child := NewByChainWithCode(nil, "ERR_CODE", "child error")
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if child.HasParent() {
		t.Error("expected no parent when slice is empty")
	}
}

func TestNewByChainWithCodef(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	child := NewByChainWithCodef(errs, "ERR_CODE", "child %s", "error")
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasParent() {
		t.Error("expected parent chain to be set")
	}
	if child.HasMetadata() {
		t.Error("expected no metadata")
	}
}

func TestNewByChainWithCodef_EmptySlice(t *testing.T) {
	child := NewByChainWithCodef(nil, "ERR_CODE", "child %s", "error")
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if child.HasParent() {
		t.Error("expected no parent when slice is empty")
	}
}

func TestNewByChainWithMetadata(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	metadata := map[string]any{"key": "value"}
	child := NewByChainWithMetadata(errs, metadata, "child error")
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasParent() {
		t.Error("expected parent chain to be set")
	}
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", child.Metadata()["key"])
	}
	if child.HasCode() {
		t.Error("expected no code")
	}
}

func TestNewByChainWithMetadata_EmptySlice(t *testing.T) {
	metadata := map[string]any{"key": "value"}
	child := NewByChainWithMetadata(nil, metadata, "child error")
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.HasParent() {
		t.Error("expected no parent when slice is empty")
	}
}

func TestNewByChainWithMetadataf(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	metadata := map[string]any{"key": "value"}
	child := NewByChainWithMetadataf(errs, metadata, "child %s", "error")
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasParent() {
		t.Error("expected parent chain to be set")
	}
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", child.Metadata()["key"])
	}
	if child.HasCode() {
		t.Error("expected no code")
	}
}

func TestNewByChainWithMetadataf_EmptySlice(t *testing.T) {
	metadata := map[string]any{"key": "value"}
	child := NewByChainWithMetadataf(nil, metadata, "child %s", "error")
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.HasParent() {
		t.Error("expected no parent when slice is empty")
	}
}

func TestNewByChainWithCodeAndMetadata(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	metadata := map[string]any{"key": "value"}
	child := NewByChainWithCodeAndMetadata(errs, "ERR_CODE", metadata, "child error")
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasParent() {
		t.Error("expected parent chain to be set")
	}
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", child.Metadata()["key"])
	}
}

func TestNewByChainWithCodeAndMetadata_EmptySlice(t *testing.T) {
	metadata := map[string]any{"key": "value"}
	child := NewByChainWithCodeAndMetadata(nil, "ERR_CODE", metadata, "child error")
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.HasParent() {
		t.Error("expected no parent when slice is empty")
	}
}

func TestNewByChainWithCodeAndMetadataf(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	metadata := map[string]any{"key": "value"}
	child := NewByChainWithCodeAndMetadataf(errs, "ERR_CODE", metadata, "child %s", "error")
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if child.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", child.Message())
	}
	if !child.HasParent() {
		t.Error("expected parent chain to be set")
	}
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", child.Metadata()["key"])
	}
}

func TestNewByChainWithCodeAndMetadataf_EmptySlice(t *testing.T) {
	metadata := map[string]any{"key": "value"}
	child := NewByChainWithCodeAndMetadataf(nil, "ERR_CODE", metadata, "child %s", "error")
	if child.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", child.Code())
	}
	if !child.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if child.HasParent() {
		t.Error("expected no parent when slice is empty")
	}
}

func TestNewByChain_ChainOrder(t *testing.T) {
	err1 := NewWithCode("CODE1", "first")
	err2 := NewWithCode("CODE2", "second")
	err3 := NewWithCode("CODE3", "third")
	child := NewByChain([]error{err1, err2, err3}, "wrapper")

	chain := child.Chain()
	if len(chain) != 3 {
		t.Fatalf("expected 3 parents in chain, got %d", len(chain))
	}
	if chain[0].Message() != "first" {
		t.Errorf("expected first parent message 'first', got '%s'", chain[0].Message())
	}
	if chain[1].Message() != "second" {
		t.Errorf("expected second parent message 'second', got '%s'", chain[1].Message())
	}
	if chain[2].Message() != "third" {
		t.Errorf("expected third parent message 'third', got '%s'", chain[2].Message())
	}
}

func TestNewByChain_IsRecognizesChainNode(t *testing.T) {
	target := NewWithCode("DEEP_ERR", "deep error")
	errs := []error{New("shallow"), target}
	child := NewByChain(errs, "wrapper")

	if !Is(child, TargetWithCode("DEEP_ERR")) {
		t.Error("expected Is to recognize code from chain node")
	}
}

func TestNewByChainAsSlice(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	result := NewByChainAsSlice(errs, "child error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", e.Message())
	}
	if !e.HasParent() {
		t.Error("expected parent chain to be set")
	}
}

func TestNewByChainAsSlice_EmptySlice(t *testing.T) {
	result := NewByChainAsSlice(nil, "child error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.HasParent() {
		t.Error("expected no parent when slice is empty")
	}
}

func TestNewByChainAsSlicef(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	result := NewByChainAsSlicef(errs, "child %s", "error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", e.Message())
	}
	if !e.HasParent() {
		t.Error("expected parent chain to be set")
	}
}

func TestNewByChainWithCodeAsSlice(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	result := NewByChainWithCodeAsSlice(errs, "ERR_CODE", "child error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", e.Code())
	}
	if !e.HasParent() {
		t.Error("expected parent chain to be set")
	}
}

func TestNewByChainWithCodeAsSlicef(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	result := NewByChainWithCodeAsSlicef(errs, "ERR_CODE", "child %s", "error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", e.Code())
	}
	if e.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", e.Message())
	}
	if !e.HasParent() {
		t.Error("expected parent chain to be set")
	}
}

func TestNewByChainWithMetadataAsSlice(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	metadata := map[string]any{"key": "value"}
	result := NewByChainWithMetadataAsSlice(errs, metadata, "child error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if !e.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if e.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", e.Metadata()["key"])
	}
	if !e.HasParent() {
		t.Error("expected parent chain to be set")
	}
}

func TestNewByChainWithMetadataAsSlicef(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	metadata := map[string]any{"key": "value"}
	result := NewByChainWithMetadataAsSlicef(errs, metadata, "child %s", "error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", e.Message())
	}
	if !e.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if !e.HasParent() {
		t.Error("expected parent chain to be set")
	}
}

func TestNewByChainWithCodeAndMetadataAsSlice(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	metadata := map[string]any{"key": "value"}
	result := NewByChainWithCodeAndMetadataAsSlice(errs, "ERR_CODE", metadata, "child error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", e.Code())
	}
	if !e.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if !e.HasParent() {
		t.Error("expected parent chain to be set")
	}
}

func TestNewByChainWithCodeAndMetadataAsSlicef(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	metadata := map[string]any{"key": "value"}
	result := NewByChainWithCodeAndMetadataAsSlicef(errs, "ERR_CODE", metadata, "child %s", "error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", e.Code())
	}
	if e.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", e.Message())
	}
	if !e.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if !e.HasParent() {
		t.Error("expected parent chain to be set")
	}
}

func TestNewByParentAsSlice(t *testing.T) {
	parent := New("parent error")
	result := NewByParentAsSlice(parent, "child error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", e.Message())
	}
	if !e.HasParent() {
		t.Error("expected parent to be set")
	}
}

func TestNewByParentAsSlice_NilParent(t *testing.T) {
	result := NewByParentAsSlice(nil, "child error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.HasParent() {
		t.Error("expected no parent when parent is nil")
	}
}

func TestNewByParentAsSlicef(t *testing.T) {
	parent := New("parent error")
	result := NewByParentAsSlicef(parent, "child %s", "error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", e.Message())
	}
	if !e.HasParent() {
		t.Error("expected parent to be set")
	}
}

func TestNewByParentWithCodeAsSlice(t *testing.T) {
	parent := New("parent error")
	result := NewByParentWithCodeAsSlice(parent, "ERR_CODE", "child error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", e.Code())
	}
	if !e.HasParent() {
		t.Error("expected parent to be set")
	}
}

func TestNewByParentWithCodeAsSlicef(t *testing.T) {
	parent := New("parent error")
	result := NewByParentWithCodeAsSlicef(parent, "ERR_CODE", "child %s", "error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", e.Code())
	}
	if e.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", e.Message())
	}
	if !e.HasParent() {
		t.Error("expected parent to be set")
	}
}

func TestNewByParentWithMetadataAsSlice(t *testing.T) {
	parent := New("parent error")
	metadata := map[string]any{"key": "value"}
	result := NewByParentWithMetadataAsSlice(parent, metadata, "child error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if !e.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if e.Metadata()["key"] != "value" {
		t.Errorf("expected metadata key 'value', got '%v'", e.Metadata()["key"])
	}
	if !e.HasParent() {
		t.Error("expected parent to be set")
	}
}

func TestNewByParentWithMetadataAsSlicef(t *testing.T) {
	parent := New("parent error")
	metadata := map[string]any{"key": "value"}
	result := NewByParentWithMetadataAsSlicef(parent, metadata, "child %s", "error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", e.Message())
	}
	if !e.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if !e.HasParent() {
		t.Error("expected parent to be set")
	}
}

func TestNewByParentWithCodeAndMetadataAsSlice(t *testing.T) {
	parent := New("parent error")
	metadata := map[string]any{"key": "value"}
	result := NewByParentWithCodeAndMetadataAsSlice(parent, "ERR_CODE", metadata, "child error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", e.Code())
	}
	if !e.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if !e.HasParent() {
		t.Error("expected parent to be set")
	}
}

func TestNewByParentWithCodeAndMetadataAsSlicef(t *testing.T) {
	parent := New("parent error")
	metadata := map[string]any{"key": "value"}
	result := NewByParentWithCodeAndMetadataAsSlicef(parent, "ERR_CODE", metadata, "child %s", "error")
	if len(result) != 1 {
		t.Fatalf("expected slice of length 1, got %d", len(result))
	}
	var e *Err
	if !As(result[0], &e) {
		t.Fatal("expected result to be *Err")
	}
	if e.Code() != "ERR_CODE" {
		t.Errorf("expected code 'ERR_CODE', got '%s'", e.Code())
	}
	if e.Message() != "child error" {
		t.Errorf("expected message 'child error', got '%s'", e.Message())
	}
	if !e.HasMetadata() {
		t.Error("expected metadata to be set")
	}
	if !e.HasParent() {
		t.Error("expected parent to be set")
	}
}
