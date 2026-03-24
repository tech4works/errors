package errors

import (
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
	err := NewWithCodef("value is %s", "ERR", "bad")
	if err.Code() != "ERR" {
		t.Errorf("expected code 'ERR', got '%s'", err.Code())
	}
	if err.Message() != "value is bad" {
		t.Errorf("expected 'value is bad', got '%s'", err.Message())
	}
}

func TestNewWithMetadataf(t *testing.T) {
	meta := map[string]any{"k": "v"}
	err := NewWithMetadataf("msg %s", meta, "here")
	if err.Message() != "msg here" {
		t.Errorf("expected 'msg here', got '%s'", err.Message())
	}
}

func TestNewWithCodeAndMetadataf(t *testing.T) {
	meta := map[string]any{"k": "v"}
	err := NewWithCodeAndMetadataf("msg %s", "C1", meta, "x")
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
	err := NewWithSkipCallerAndCodef(1, "fmt %s", "FMT01", "val")
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
	err := NewWithAllf(1, "msg %s", "ALLF", meta, "x")
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
	errs := NewWithCodeAsSlicef("msg %s", "C1", "x")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithMetadataAsSlicef(t *testing.T) {
	errs := NewWithMetadataAsSlicef("msg %s", map[string]any{"k": "v"}, "x")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestNewWithCodeAndMetadataAsSlicef(t *testing.T) {
	errs := NewWithCodeAndMetadataAsSlicef("msg %s", "C1", map[string]any{"k": "v"}, "x")
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
	errs := NewWithSkipCallerAndCodeAsSlicef(1, "msg %s", "C1", "x")
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
	errs := NewWithAllAsSlicef(1, "msg %s", "C1", map[string]any{"k": "v"}, "x")
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
