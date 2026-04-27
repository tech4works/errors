package errors

import (
	"errors"
	"strings"
	"testing"
)

// --- Is / IsNot ---

func TestIs_NilErr(t *testing.T) {
	if Is(nil, New("x")) {
		t.Error("Is(nil, target) should be false")
	}
}

func TestIs_NilTarget(t *testing.T) {
	if Is(New("x"), nil) {
		t.Error("Is(err, nil) should be false")
	}
}

func TestIs_SameMessage(t *testing.T) {
	err := New("same")
	target := New("same")
	// snapshot comparison
	if !Is(err, target) {
		t.Error("expected Is to be true for same message")
	}
}

func TestIs_DifferentMessage(t *testing.T) {
	err := New("foo")
	target := New("bar")
	if Is(err, target) {
		t.Error("expected Is to be false for different messages")
	}
}

func TestIs_TargetIsTarget(t *testing.T) {
	target := TargetWithCode("T01")
	err := NewWithCode("T01", "some error")
	if !Is(err, target) {
		t.Error("expected Is to be true when err contains target code")
	}
}

func TestIs_TargetWithMessage(t *testing.T) {
	target := TargetWithMessage("specific error")
	err := New("specific error")
	if !Is(err, target) {
		t.Error("expected Is to be true when err contains target message")
	}
}

func TestIs_ErrIsTarget(t *testing.T) {
	target := TargetWithCode("T01")
	// err itself is a target — should return false
	if Is(target, New("x")) {
		t.Error("Is(target, x) should be false when err is a target")
	}
}

func TestIs_PlainErrors(t *testing.T) {
	err := errors.New("plain error")
	target := errors.New("plain error")
	if !Is(err, target) {
		t.Error("expected Is to be true for identical plain errors")
	}
}

func TestIsNot(t *testing.T) {
	err := New("foo")
	target := New("bar")
	if !IsNot(err, target) {
		t.Error("expected IsNot to be true")
	}
}

// --- ContainsCode / NotContainsCode ---

func TestContainsCode_True(t *testing.T) {
	err := NewWithCode("ERR001", "msg")
	if !ContainsCode(err, "ERR001") {
		t.Error("expected ContainsCode to be true")
	}
}

func TestContainsCode_False(t *testing.T) {
	err := New("msg")
	if ContainsCode(err, "ERR001") {
		t.Error("expected ContainsCode to be false")
	}
}

func TestContainsCode_NilErr(t *testing.T) {
	if ContainsCode(nil, "ERR001") {
		t.Error("expected ContainsCode(nil, ...) to be false")
	}
}

func TestContainsCode_EmptyCode(t *testing.T) {
	err := New("msg")
	if ContainsCode(err, "") {
		t.Error("expected ContainsCode(err, '') to be false")
	}
}

func TestNotContainsCode(t *testing.T) {
	err := New("msg")
	if !NotContainsCode(err, "ERR001") {
		t.Error("expected NotContainsCode to be true")
	}
}

// --- Contains / NotContains ---

func TestContains_True(t *testing.T) {
	target := New("target")
	errs := []error{New("other"), target}
	if !Contains(errs, target) {
		t.Error("expected Contains to be true")
	}
}

func TestContains_False(t *testing.T) {
	errs := []error{New("other")}
	if Contains(errs, New("missing")) {
		t.Error("expected Contains to be false")
	}
}

func TestContains_EmptySlice(t *testing.T) {
	if Contains([]error{}, New("x")) {
		t.Error("expected Contains(empty, ...) to be false")
	}
}

func TestContains_NilTarget(t *testing.T) {
	if Contains([]error{New("x")}, nil) {
		t.Error("expected Contains(..., nil) to be false")
	}
}

func TestNotContains(t *testing.T) {
	errs := []error{New("other")}
	if !NotContains(errs, New("missing")) {
		t.Error("expected NotContains to be true")
	}
}

// --- Only / NotOnly ---

func TestOnly_True(t *testing.T) {
	target := TargetWithMessage("same")
	errs := []error{New("same"), New("same")}
	if !Only(errs, target) {
		t.Error("expected Only to be true")
	}
}

func TestOnly_False(t *testing.T) {
	target := TargetWithMessage("same")
	errs := []error{New("same"), New("different")}
	if Only(errs, target) {
		t.Error("expected Only to be false")
	}
}

func TestOnly_EmptySlice(t *testing.T) {
	if Only([]error{}, New("x")) {
		t.Error("expected Only(empty, ...) to be false")
	}
}

func TestOnly_NilTarget(t *testing.T) {
	if Only([]error{New("x")}, nil) {
		t.Error("expected Only(..., nil) to be false")
	}
}

func TestNotOnly(t *testing.T) {
	target := TargetWithMessage("same")
	errs := []error{New("same"), New("different")}
	if !NotOnly(errs, target) {
		t.Error("expected NotOnly to be true")
	}
}

// --- As ---

func TestAs_WithErrPointerPointer(t *testing.T) {
	err := New("as test")
	var target *Err
	if !As(err, &target) {
		t.Error("expected As to return true")
	}
	if target == nil {
		t.Error("expected target to be non-nil")
	}
}

func TestAs_WithErrPointer(t *testing.T) {
	err := New("as test")
	var target Err
	if !As(err, &target) {
		t.Error("expected As to return true")
	}
}

func TestAs_NilErr(t *testing.T) {
	var target *Err
	if As(nil, &target) {
		t.Error("expected As(nil, ...) to be false")
	}
}

func TestAs_NilTarget(t *testing.T) {
	err := New("x")
	if As(err, nil) {
		t.Error("expected As(err, nil) to be false")
	}
}

func TestAs_PlainError(t *testing.T) {
	err := errors.New("plain")
	var target *Err
	// plain error won't be extracted as *Err
	result := As(err, &target)
	// result may be false since it's not an *Err
	_ = result
}

func TestAs_StandardError(t *testing.T) {
	err := errors.New("std")
	var target error
	if !As(err, &target) {
		t.Error("expected As to work with standard error interface")
	}
}

// --- Wrap ---

func TestWrap_Nil(t *testing.T) {
	if Wrap(nil) != nil {
		t.Error("expected Wrap(nil) to be nil")
	}
}

func TestWrap_Err(t *testing.T) {
	err := New("wrap me")
	wrapped := Wrap(err)
	if wrapped == nil {
		t.Error("expected non-nil wrapped error")
	}
	if wrapped.Message() != "wrap me" {
		t.Errorf("expected 'wrap me', got '%s'", wrapped.Message())
	}
}

func TestWrap_PlainError(t *testing.T) {
	err := errors.New("plain")
	wrapped := Wrap(err)
	if wrapped == nil {
		t.Error("expected non-nil wrapped plain error")
	}
}

// --- Inherit ---

func TestInherit_Nil(t *testing.T) {
	if Inherit(nil, "msg") != nil {
		t.Error("expected Inherit(nil) to be nil")
	}
}

func TestInherit_WithMessage(t *testing.T) {
	parent := New("parent")
	child := Inherit(parent, "child")
	if child.Message() != "child" {
		t.Errorf("expected 'child', got '%s'", child.Message())
	}
}

func TestInherit_EmptyMessage(t *testing.T) {
	parent := New("parent msg")
	child := Inherit(parent)
	// empty message inherits parent's message
	if child.Message() != "parent msg" {
		t.Errorf("expected 'parent msg', got '%s'", child.Message())
	}
}

func TestInherit_InheritsCode(t *testing.T) {
	parent := NewWithCode("P01", "parent")
	child := Inherit(parent, "child")
	if child.Code() != "P01" {
		t.Errorf("expected inherited code 'P01', got '%s'", child.Code())
	}
}

func TestInheritf_Nil(t *testing.T) {
	if Inheritf(nil, "fmt %s", "x") != nil {
		t.Error("expected Inheritf(nil) to be nil")
	}
}

func TestInheritf(t *testing.T) {
	parent := New("parent")
	child := Inheritf(parent, "child %s", "1")
	if child.Message() != "child 1" {
		t.Errorf("expected 'child 1', got '%s'", child.Message())
	}
}

func TestInheritWithSkipCaller_Nil(t *testing.T) {
	if InheritWithSkipCaller(nil, 1, "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritWithSkipCaller(t *testing.T) {
	parent := New("parent")
	child := InheritWithSkipCaller(parent, 1, "child")
	if child.Message() != "child" {
		t.Errorf("expected 'child', got '%s'", child.Message())
	}
}

func TestInheritWithCode_Nil(t *testing.T) {
	if InheritWithCode(nil, "C1", "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritWithCode(t *testing.T) {
	parent := New("parent")
	child := InheritWithCode(parent, "C1", "child")
	if child.Code() != "C1" {
		t.Errorf("expected 'C1', got '%s'", child.Code())
	}
}

func TestInheritWithCode_EmptyMessage(t *testing.T) {
	parent := New("parent msg")
	child := InheritWithCode(parent, "C1")
	if child.Message() != "parent msg" {
		t.Errorf("expected 'parent msg', got '%s'", child.Message())
	}
}

func TestInheritWithSkipCallerf_Nil(t *testing.T) {
	if InheritWithSkipCallerf(nil, 1, "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestInheritWithSkipCallerf(t *testing.T) {
	parent := New("parent")
	child := InheritWithSkipCallerf(parent, 1, "child %s", "2")
	if child.Message() != "child 2" {
		t.Errorf("expected 'child 2', got '%s'", child.Message())
	}
}

func TestInheritWithAll_Nil(t *testing.T) {
	if InheritWithAll(nil, 1, "C1", nil, "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritWithAll(t *testing.T) {
	parent := NewWithMetadata(map[string]any{"p": "v"}, "parent")
	child := InheritWithAll(parent, 1, "C1", map[string]any{"c": "x"}, "child")
	if child.Code() != "C1" {
		t.Errorf("expected 'C1', got '%s'", child.Code())
	}
	if child.Metadata()["c"] != "x" {
		t.Error("expected metadata c=x")
	}
}

func TestInheritWithAll_EmptyMessage(t *testing.T) {
	parent := New("parent msg")
	child := InheritWithAll(parent, 1, "C1", nil)
	if child.Message() != "parent msg" {
		t.Errorf("expected 'parent msg', got '%s'", child.Message())
	}
}

func TestInheritWithSkipCallerAndCode_Nil(t *testing.T) {
	if InheritWithSkipCallerAndCode(nil, 1, "C1", "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritWithSkipCallerAndCode(t *testing.T) {
	parent := NewWithCode("SC01", "parent")
	child := InheritWithSkipCallerAndCode(parent, 1, "SC01", "child")
	if child.Code() != "SC01" {
		t.Errorf("expected 'SC01', got '%s'", child.Code())
	}
}

// --- InheritAsSlice variants ---

func TestInheritAsSlice_Nil(t *testing.T) {
	if InheritAsSlice(nil, "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritAsSlice(t *testing.T) {
	parent := New("parent")
	errs := InheritAsSlice(parent, "child")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestInheritAsSlicef_Nil(t *testing.T) {
	if InheritAsSlicef(nil, "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestInheritAsSlicef(t *testing.T) {
	parent := New("parent")
	errs := InheritAsSlicef(parent, "child %d", 1)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestInheritWithSkipCallerAsSlice_Nil(t *testing.T) {
	if InheritWithSkipCallerAsSlice(nil, 1, "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritWithSkipCallerAsSlice(t *testing.T) {
	parent := New("parent")
	errs := InheritWithSkipCallerAsSlice(parent, 1, "child")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestInheritWithCodeAsSlice_Nil(t *testing.T) {
	if InheritWithCodeAsSlice(nil, "C1", "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritWithCodeAsSlice(t *testing.T) {
	parent := New("parent")
	errs := InheritWithCodeAsSlice(parent, "C1", "child")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestInheritWithSkipCallerAsSlicef_Nil(t *testing.T) {
	if InheritWithSkipCallerAsSlicef(nil, 1, "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestInheritWithSkipCallerAsSlicef(t *testing.T) {
	parent := New("parent")
	errs := InheritWithSkipCallerAsSlicef(parent, 1, "child %d", 1)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestInheritWithAllAsSlice_Nil(t *testing.T) {
	if InheritWithAllAsSlice(nil, 1, "C1", nil, "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritWithAllAsSlice(t *testing.T) {
	parent := New("parent")
	errs := InheritWithAllAsSlice(parent, 1, "C1", nil, "child")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestInheritWithSkipCallerAndCodeAsSlice_Nil(t *testing.T) {
	if InheritWithSkipCallerAndCodeAsSlice(nil, 1, "C1", "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritWithSkipCallerAndCodeAsSlice(t *testing.T) {
	parent := New("parent")
	errs := InheritWithSkipCallerAndCodeAsSlice(parent, 1, "C1", "child")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

// --- Join ---

func TestJoin(t *testing.T) {
	errs := []error{New("a"), New("b")}
	joined := Join(errs, ", ")
	if joined == nil {
		t.Error("expected non-nil joined error")
	}
	// deve retornar *Err, não errors.New da stdlib
	var e *Err
	if !As(joined, &e) {
		t.Error("expected Join to return *Err")
	}
}

func TestJoin_Empty(t *testing.T) {
	if Join([]error{}, ", ") != nil {
		t.Error("expected Join(empty) to be nil")
	}
}

func TestJoin_IsRecognized(t *testing.T) {
	target := TargetWithMessage("a")
	errs := []error{New("a"), New("b")}
	joined := Join(errs, ", ")
	// Is deve conseguir identificar o erro dentro do join
	if !Is(joined, target) {
		t.Error("expected Is to find target message inside joined error")
	}
}

func TestJoin_WrapPreservesStructure(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	joined := Join(errs, " | ")
	wrapped := Wrap(joined)
	if wrapped == nil {
		t.Error("expected Wrap of joined *Err to be non-nil")
	}
	if wrapped.Message() == "" {
		t.Error("expected non-empty message after Wrap")
	}
}

func TestJoinToString_Empty(t *testing.T) {
	result := JoinToString([]error{}, ", ")
	if result != "" {
		t.Errorf("expected empty string, got '%s'", result)
	}
}

func TestJoinToString(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := JoinToString(errs, " | ")
	if result == "" {
		t.Error("expected non-empty JoinToString result")
	}
}

func TestJoinInherit_Nil(t *testing.T) {
	if JoinInherit(nil, ", ", "msg") != nil {
		t.Error("expected nil")
	}
}

func TestJoinInherit(t *testing.T) {
	errs := []error{New("a"), New("b")}
	err := JoinInherit(errs, ", ", "joined")
	if err == nil {
		t.Error("expected non-nil error")
	}
}

func TestJoinInheritf_Nil(t *testing.T) {
	if JoinInheritf(nil, ", ", "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestJoinInheritf(t *testing.T) {
	errs := []error{New("a"), New("b")}
	err := JoinInheritf(errs, ", ", "joined %d", 1)
	if err == nil {
		t.Error("expected non-nil error")
	}
}

func TestJoinInheritAsSlice_Nil(t *testing.T) {
	if JoinInheritAsSlice(nil, ", ", "msg") != nil {
		t.Error("expected nil")
	}
}

func TestJoinInheritAsSlice(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := JoinInheritAsSlice(errs, ", ", "joined")
	if len(result) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result))
	}
}

func TestJoinInheritAsSlicef_Nil(t *testing.T) {
	if JoinInheritAsSlicef(nil, ", ", "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestJoinInheritAsSlicef(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := JoinInheritAsSlicef(errs, ", ", "joined %d", 1)
	if len(result) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result))
	}
}

func TestJoinInheritWithSkipCaller_Nil(t *testing.T) {
	if JoinInheritWithSkipCaller(nil, ", ", 1, "msg") != nil {
		t.Error("expected nil")
	}
}

func TestJoinInheritWithSkipCaller(t *testing.T) {
	errs := []error{New("a"), New("b")}
	err := JoinInheritWithSkipCaller(errs, ", ", 1, "joined")
	if err == nil {
		t.Error("expected non-nil error")
	}
}

func TestJoinInheritWithCode_Nil(t *testing.T) {
	if JoinInheritWithCode(nil, ", ", "C1", "msg") != nil {
		t.Error("expected nil")
	}
}

func TestJoinInheritWithCode(t *testing.T) {
	errs := []error{New("a"), New("b")}
	err := JoinInheritWithCode(errs, ", ", "C1", "joined")
	if err == nil {
		t.Error("expected non-nil error")
	}
}

func TestJoinInheritWithSkipCallerf_Nil(t *testing.T) {
	if JoinInheritWithSkipCallerf(nil, ", ", 1, "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestJoinInheritWithSkipCallerf(t *testing.T) {
	errs := []error{New("a"), New("b")}
	err := JoinInheritWithSkipCallerf(errs, ", ", 1, "joined %d", 1)
	if err == nil {
		t.Error("expected non-nil error")
	}
}

func TestJoinInheritWithAll_Nil(t *testing.T) {
	if JoinInheritWithAll(nil, ", ", 1, "C1", nil, "msg") != nil {
		t.Error("expected nil")
	}
}

func TestJoinInheritWithAll(t *testing.T) {
	errs := []error{New("a"), New("b")}
	err := JoinInheritWithAll(errs, ", ", 1, "C1", nil, "joined")
	if err == nil {
		t.Error("expected non-nil error")
	}
}

func TestJoinInheritWithSkipCallerAndCode_Nil(t *testing.T) {
	if JoinInheritWithSkipCallerAndCode(nil, ", ", 1, "C1", "msg") != nil {
		t.Error("expected nil")
	}
}

func TestJoinInheritWithSkipCallerAndCode(t *testing.T) {
	errs := []error{New("a"), New("b")}
	err := JoinInheritWithSkipCallerAndCode(errs, ", ", 1, "C1", "joined")
	if err == nil {
		t.Error("expected non-nil error")
	}
}

// --- Chain / ChainFromSlice ---

func TestChain_Empty(t *testing.T) {
	if Chain() != nil {
		t.Error("expected nil for no args")
	}
}

func TestChain_Single(t *testing.T) {
	result := Chain(New("only"))
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Message() != "only" {
		t.Errorf("expected 'only', got '%s'", result.Message())
	}
	if result.parent != nil {
		t.Error("expected no parent for single error")
	}
}

func TestChain_Order(t *testing.T) {
	result := Chain(New("first"), New("second"), New("third"))
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Message() != "first" {
		t.Errorf("expected root 'first', got '%s'", result.Message())
	}
	full := result.Error()
	idx1 := strings.Index(full, "first")
	idx2 := strings.Index(full, "second")
	idx3 := strings.Index(full, "third")
	if idx1 < 0 || idx2 < 0 || idx3 < 0 {
		t.Fatalf("expected all messages in chain output, got: %s", full)
	}
	if !(idx1 < idx2 && idx2 < idx3) {
		t.Errorf("expected order first < second < third, got positions %d %d %d", idx1, idx2, idx3)
	}
}

func TestChain_IsRecognizesAnyNode(t *testing.T) {
	target := TargetWithMessage("second")
	result := Chain(New("first"), New("second"), New("third"))
	if !Is(result, target) {
		t.Error("expected Is to find 'second' anywhere in the chain")
	}
}

func TestChainFromSlice_Nil(t *testing.T) {
	if ChainFromSlice(nil) != nil {
		t.Error("expected nil for nil slice")
	}
}

func TestChainFromSlice_Empty(t *testing.T) {
	if ChainFromSlice([]error{}) != nil {
		t.Error("expected nil for empty slice")
	}
}

func TestChainFromSlice_Order(t *testing.T) {
	errs := []error{New("a"), New("b"), New("c")}
	result := ChainFromSlice(errs)
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Message() != "a" {
		t.Errorf("expected root 'a', got '%s'", result.Message())
	}
	full := result.Error()
	if !strings.Contains(full, "b") || !strings.Contains(full, "c") {
		t.Errorf("expected 'b' and 'c' in chain output, got: %s", full)
	}
}

func TestChainFromSlice_CompatibleWithInheritFromSlice(t *testing.T) {
	errs := []error{New("x"), New("y")}
	fromChain := ChainFromSlice(errs)
	fromInherit := InheritFromSlice(errs)
	// both should produce the same root message and chain depth
	if fromChain.Message() != fromInherit.Message() {
		t.Errorf("expected same root message: chain=%s inherit=%s", fromChain.Message(), fromInherit.Message())
	}
}

// --- InheritFromSlice ---

func TestInheritFromSlice_Nil(t *testing.T) {
	if InheritFromSlice(nil, "msg") != nil {
		t.Error("expected nil for nil slice")
	}
}

func TestInheritFromSlice_Empty(t *testing.T) {
	if InheritFromSlice([]error{}, "msg") != nil {
		t.Error("expected nil for empty slice")
	}
}

func TestInheritFromSlice_MessageSet(t *testing.T) {
	errs := []error{New("err1"), New("err2")}
	result := InheritFromSlice(errs, "root message")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Message() != "root message" {
		t.Errorf("expected 'root message', got '%s'", result.Message())
	}
}

func TestInheritFromSlice_EmptyMsgFallsBackToChain(t *testing.T) {
	errs := []error{New("first error"), New("second error")}
	result := InheritFromSlice(errs)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// without msg, the chain itself is returned — root is the first error
	if result.Message() != "first error" {
		t.Errorf("expected 'first error', got '%s'", result.Message())
	}
	// no synthetic wrapper: parent holds the remaining chain
	if result.parent == nil {
		t.Error("expected parent to be set (second error in chain)")
	}
}

func TestInheritFromSlice_ChainPreservesOrder(t *testing.T) {
	errs := []error{New("err1"), New("err2"), New("err3")}
	result := InheritFromSlice(errs, "root")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	// a cadeia de parent deve conter todos os erros na ordem
	full := result.Error()
	idx1 := strings.Index(full, "err1")
	idx2 := strings.Index(full, "err2")
	idx3 := strings.Index(full, "err3")
	if idx1 < 0 || idx2 < 0 || idx3 < 0 {
		t.Errorf("expected all errors in chain, got: %s", full)
	}
	if !(idx1 < idx2 && idx2 < idx3) {
		t.Errorf("expected order err1 < err2 < err3 in chain, got positions %d %d %d", idx1, idx2, idx3)
	}
}

func TestInheritFromSlice_SingleError(t *testing.T) {
	errs := []error{New("only one")}
	result := InheritFromSlice(errs, "root")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Message() != "root" {
		t.Errorf("expected 'root', got '%s'", result.Message())
	}
	if result.parent == nil {
		t.Error("expected parent to be set")
	}
}

func TestInheritFromSlice_IsRecognizesChildError(t *testing.T) {
	target := TargetWithMessage("err2")
	errs := []error{New("err1"), New("err2"), New("err3")}
	result := InheritFromSlice(errs, "root")
	// Is deve encontrar err2 na cadeia via Error()
	if !Is(result, target) {
		t.Error("expected Is to find 'err2' inside the inherited chain")
	}
}

func TestInheritFromSlice_WrapReturnsErr(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := InheritFromSlice(errs, "root")
	wrapped := Wrap(result)
	if wrapped == nil {
		t.Error("expected Wrap to return non-nil")
	}
	if wrapped.Message() != "root" {
		t.Errorf("expected 'root', got '%s'", wrapped.Message())
	}
}

// --- InheritFromSlicef ---

func TestInheritFromSlicef_Nil(t *testing.T) {
	if InheritFromSlicef(nil, "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestInheritFromSlicef(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := InheritFromSlicef(errs, "root %d errors", 2)
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Message() != "root 2 errors" {
		t.Errorf("expected 'root 2 errors', got '%s'", result.Message())
	}
}

// --- InheritFromSliceWithCode ---

func TestInheritFromSliceWithCode_Nil(t *testing.T) {
	if InheritFromSliceWithCode(nil, "C1", "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritFromSliceWithCode(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := InheritFromSliceWithCode(errs, "UPLOAD_ERR", "upload failed")
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Code() != "UPLOAD_ERR" {
		t.Errorf("expected code 'UPLOAD_ERR', got '%s'", result.Code())
	}
	if result.Message() != "upload failed" {
		t.Errorf("expected 'upload failed', got '%s'", result.Message())
	}
}

func TestInheritFromSliceWithCode_IsRecognizedByCode(t *testing.T) {
	target := TargetWithCode("UPLOAD_ERR")
	errs := []error{New("a"), New("b")}
	result := InheritFromSliceWithCode(errs, "UPLOAD_ERR", "upload failed")
	if !Is(result, target) {
		t.Error("expected Is to match by code")
	}
}

// --- InheritFromSliceWithCodef ---

func TestInheritFromSliceWithCodef_Nil(t *testing.T) {
	if InheritFromSliceWithCodef(nil, "C1", "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestInheritFromSliceWithCodef(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := InheritFromSliceWithCodef(errs, "C1", "failed %d files", 2)
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Code() != "C1" {
		t.Errorf("expected code 'C1', got '%s'", result.Code())
	}
	if result.Message() != "failed 2 files" {
		t.Errorf("expected 'failed 2 files', got '%s'", result.Message())
	}
}

// --- InheritFromSliceWithAll ---

func TestInheritFromSliceWithAll_Nil(t *testing.T) {
	if InheritFromSliceWithAll(nil, 1, "C1", nil, "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritFromSliceWithAll(t *testing.T) {
	meta := map[string]any{"bucket": "manas-social"}
	errs := []error{New("a"), New("b")}
	result := InheritFromSliceWithAll(errs, 1, "C1", meta, "upload failed")
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Code() != "C1" {
		t.Errorf("expected code 'C1', got '%s'", result.Code())
	}
	if result.Metadata()["bucket"] != "manas-social" {
		t.Error("expected metadata bucket=manas-social")
	}
}

// --- InheritFromSliceWithAllf ---

func TestInheritFromSliceWithAllf_Nil(t *testing.T) {
	if InheritFromSliceWithAllf(nil, 1, "C1", nil, "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestInheritFromSliceWithAllf(t *testing.T) {
	meta := map[string]any{"count": 3}
	errs := []error{New("a"), New("b"), New("c")}
	result := InheritFromSliceWithAllf(errs, 1, "C1", meta, "failed %d uploads", 3)
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Message() != "failed 3 uploads" {
		t.Errorf("expected 'failed 3 uploads', got '%s'", result.Message())
	}
}

// --- InheritFromSliceWithSkipCaller ---

func TestInheritFromSliceWithSkipCaller_Nil(t *testing.T) {
	if InheritFromSliceWithSkipCaller(nil, 1, "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritFromSliceWithSkipCaller(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := InheritFromSliceWithSkipCaller(errs, 1, "root")
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Message() != "root" {
		t.Errorf("expected 'root', got '%s'", result.Message())
	}
}

// --- InheritFromSliceWithSkipCallerf ---

func TestInheritFromSliceWithSkipCallerf_Nil(t *testing.T) {
	if InheritFromSliceWithSkipCallerf(nil, 1, "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestInheritFromSliceWithSkipCallerf(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := InheritFromSliceWithSkipCallerf(errs, 1, "root %s", "msg")
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Message() != "root msg" {
		t.Errorf("expected 'root msg', got '%s'", result.Message())
	}
}

// --- InheritFromSliceWithSkipCallerAndCode ---

func TestInheritFromSliceWithSkipCallerAndCode_Nil(t *testing.T) {
	if InheritFromSliceWithSkipCallerAndCode(nil, 1, "C1", "msg") != nil {
		t.Error("expected nil")
	}
}

func TestInheritFromSliceWithSkipCallerAndCode(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := InheritFromSliceWithSkipCallerAndCode(errs, 1, "SC01", "root")
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Code() != "SC01" {
		t.Errorf("expected code 'SC01', got '%s'", result.Code())
	}
}

// --- InheritFromSliceWithSkipCallerAndCodef ---

func TestInheritFromSliceWithSkipCallerAndCodef_Nil(t *testing.T) {
	if InheritFromSliceWithSkipCallerAndCodef(nil, 1, "C1", "fmt %s", "x") != nil {
		t.Error("expected nil")
	}
}

func TestInheritFromSliceWithSkipCallerAndCodef(t *testing.T) {
	errs := []error{New("a"), New("b")}
	result := InheritFromSliceWithSkipCallerAndCodef(errs, 1, "SC01", "root %d", 42)
	if result == nil {
		t.Fatal("expected non-nil")
	}
	if result.Code() != "SC01" {
		t.Errorf("expected code 'SC01', got '%s'", result.Code())
	}
	if result.Message() != "root 42" {
		t.Errorf("expected 'root 42', got '%s'", result.Message())
	}
}

func TestExtract_PlainError(t *testing.T) {
	plain := errors.New("plain error")
	extracted, e := extract(plain, 1)
	if extracted {
		t.Error("expected extracted=false for plain error")
	}
	if e == nil {
		t.Error("expected non-nil Err for plain error")
	}
}

func TestExtract_Nil(t *testing.T) {
	extracted, e := extract(nil, 1)
	if extracted || e != nil {
		t.Error("expected false, nil for nil error")
	}
}

func TestExtract_FormattedError(t *testing.T) {
	// Create a formatted error string that matches the regex
	original := New("formatted")
	formatted := errors.New(original.Error())
	extracted, e := extract(formatted, 1)
	if !extracted {
		t.Error("expected extracted=true for formatted error string")
	}
	if e == nil {
		t.Error("expected non-nil Err")
	}
}

// --- isValidAsTarget ---

func TestIsValidAsTarget_ValidPointer(t *testing.T) {
	var e error
	if !isValidAsTarget(&e) {
		t.Error("expected true for *error")
	}
}

func TestIsValidAsTarget_NilPointer(t *testing.T) {
	var p *int
	if isValidAsTarget(p) {
		t.Error("expected false for nil pointer")
	}
}

func TestIsValidAsTarget_NonPointer(t *testing.T) {
	if isValidAsTarget(42) {
		t.Error("expected false for non-pointer")
	}
}
