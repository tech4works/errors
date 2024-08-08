package errors

import (
	"strings"
	"testing"
)

func TestDetail_Error(t *testing.T) {
	e := &Detail{file: "file.go", line: "10", funcName: "function", message: "message", stack: "stack trace "}
	if got := e.Error(); !strings.Contains(got, "[CAUSE]: ") && strings.Contains(got, "[STACK]: ") {
		t.Errorf("Detail.Error() = %v, want %v", got, "[CAUSE]: (file.go:10) function: message [STACK]: stack trace")
	}
}

func TestDetail_PrintStackTrace(t *testing.T) {
	e := New("test")
	Details(e).PrintStackTrace()
}

func TestDetail_PrintCause(t *testing.T) {
	e := New("test")
	Details(e).PrintCause()
}

func TestDetail_Cause(t *testing.T) {
	e := &Detail{file: "file.go", line: "10", funcName: "function", message: "message", stack: "stack trace "}
	if got := e.Cause(); !strings.Contains(got, "(file.go:10) function: message") {
		t.Errorf("Detail.Cause() = %v, want %v", got, "(file.go:10) function: message")
	}
}

func TestDetail_Message(t *testing.T) {
	e := &Detail{message: "message"}
	if got := e.Message(); got != "message" {
		t.Errorf("Detail.Message() = %v, want %v", got, "message")
	}
}

func TestDetail_File(t *testing.T) {
	e := &Detail{file: "file.go"}
	if got := e.File(); got != "file.go" {
		t.Errorf("Detail.File() = %v, want %v", got, "file.go")
	}
}

func TestDetail_Line(t *testing.T) {
	e := &Detail{line: "10"}
	if got := e.Line(); got != 10 {
		t.Errorf("Detail.Line() = %v, want %v", got, 10)
	}
}

func TestDetail_Func(t *testing.T) {
	e := &Detail{funcName: "function"}
	if got := e.Func(); got != "function" {
		t.Errorf("Detail.Func() = %v, want %v", got, "function")
	}
}

func TestDetail_Stack(t *testing.T) {
	e := &Detail{stack: "stack trace "}
	if got := e.Stack(); got != "stack trace " {
		t.Errorf("Detail.Stack() = %v, want %v", got, "stack trace ")
	}
}

func TestNew(t *testing.T) {
	err := New("some error")
	if err == nil {
		t.Error("New() should not return nil")
	}
}

func TestNewf(t *testing.T) {
	err := Newf("%s error", "some")
	if err == nil {
		t.Error("Newf() should not return nil")
	}
}

func TestNewSkipCaller(t *testing.T) {
	err := NewSkipCaller(1, "some error")
	if err == nil {
		t.Error("NewSkipCaller() should not return nil")
	}
}

func TestNewSkipCallerf(t *testing.T) {
	err := NewSkipCallerf(1, "%s error", "some")
	if err == nil {
		t.Error("NewSkipCallerf() should not return nil")
	}
}
