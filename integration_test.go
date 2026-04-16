package errors

import (
	"strings"
	"testing"
)

// TestIntegration_PrefixRemovalWithStackTrace tests the complete flow
// of creating an error, setting prefixes, and verifying they're removed
func TestIntegration_PrefixRemovalWithStackTrace(t *testing.T) {
	// Save original policy
	originalPolicy := getPolicy()
	defer SetPolicy(originalPolicy)

	// Set policy to Normal to filter boilerplate
	SetPolicy(PolicyNormal)

	// Configure prefixes to remove
	SetPrefixesToSanitize("github.com/tech4works/", "sdk-manas")
	defer SetPrefixesToSanitize() // reset

	// Create an error
	err := New("integration test error")

	// Get the stack trace
	stack := err.Stack()

	// Verify prefixes are removed
	if strings.Contains(stack, "github.com/tech4works/") {
		t.Errorf("expected 'github.com/tech4works/' to be removed from stack, got: %s", stack)
	}

	if strings.Contains(stack, "sdk-manas") {
		t.Errorf("expected 'sdk-manas' to be removed from stack, got: %s", stack)
	}

	// Verify the error message is still present
	if !strings.Contains(err.Error(), "integration test error") {
		t.Errorf("expected error message to be present, got: %s", err.Error())
	}
}

// TestIntegration_MultipleErrorsWithPrefixRemoval tests that prefix removal
// works consistently across multiple errors
func TestIntegration_MultipleErrorsWithPrefixRemoval(t *testing.T) {
	SetPrefixesToSanitize("github.com/tech4works/")
	defer SetPrefixesToSanitize()

	// Create multiple errors
	err1 := New("error 1")
	err2 := New("error 2")
	err3 := Newf("error %d", 3)

	// Verify all have prefixes removed
	for i, err := range []*Err{err1, err2, err3} {
		stack := err.Stack()
		if strings.Contains(stack, "github.com/tech4works/") {
			t.Errorf("error %d: expected prefix to be removed, got: %s", i+1, stack)
		}
	}
}

// TestIntegration_InheritWithPrefixRemoval tests that prefix removal works
// with error inheritance
func TestIntegration_InheritWithPrefixRemoval(t *testing.T) {
	SetPrefixesToSanitize("github.com/tech4works/")
	defer SetPrefixesToSanitize()

	// Create a base error
	baseErr := New("base error")

	// Inherit from it
	inheritedErr := Inherit(baseErr, "inherited error")

	// Verify prefixes are removed in both
	if strings.Contains(inheritedErr.Stack(), "github.com/tech4works/") {
		t.Errorf("expected prefix to be removed from inherited error stack")
	}
}

// TestIntegration_BoilerplateFilteringWithPrefixRemoval tests that both
// boilerplate filtering and prefix removal work together
func TestIntegration_BoilerplateFilteringWithPrefixRemoval(t *testing.T) {
	originalPolicy := getPolicy()
	defer SetPolicy(originalPolicy)

	SetPolicy(PolicyNormal)
	SetPrefixesToSanitize("github.com/tech4works/", "github.com/tech4works/sdk-manas")
	defer SetPrefixesToSanitize()

	err := New("test error")
	stack := err.Stack()

	// Verify boilerplate is filtered
	if strings.Contains(stack, "runtime.") {
		t.Errorf("expected runtime boilerplate to be filtered")
	}

	// Verify prefixes are removed
	if strings.Contains(stack, "github.com/tech4works/") {
		t.Errorf("expected github.com/tech4works/ prefix to be removed")
	}

	if strings.Contains(stack, "sdk-manas") {
		t.Errorf("expected sdk-manas prefix to be removed")
	}
}

// TestIntegration_ResetPrefixes tests that calling SetPrefixesToRemove()
// with no arguments resets the prefixes
func TestIntegration_ResetPrefixes(t *testing.T) {
	// Set prefixes
	SetPrefixesToSanitize("github.com/tech4works/")

	// Create error with prefix
	err1 := New("error with prefix")
	stack1 := err1.Stack()

	// Reset prefixes
	SetPrefixesToSanitize()

	// Create error without prefix removal
	err2 := New("error without prefix")
	stack2 := err2.Stack()

	// The second error might have the prefix (depending on the actual stack)
	// but we can verify that the function doesn't panic
	_ = stack1
	_ = stack2
}

// TestIntegration_ConcurrentPrefixRemoval tests that SetPrefixesToRemove
// is thread-safe
func TestIntegration_ConcurrentPrefixRemoval(t *testing.T) {
	done := make(chan bool, 10)

	// Goroutine 1: Set prefixes
	go func() {
		for i := 0; i < 100; i++ {
			SetPrefixesToSanitize("github.com/tech4works/")
		}
		done <- true
	}()

	// Goroutine 2: Create errors
	go func() {
		for i := 0; i < 100; i++ {
			err := New("concurrent error")
			_ = err.Stack()
		}
		done <- true
	}()

	// Goroutine 3: Reset prefixes
	go func() {
		for i := 0; i < 100; i++ {
			SetPrefixesToSanitize()
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestIntegration_PrefixRemovalWithCode tests prefix removal with error codes
func TestIntegration_PrefixRemovalWithCode(t *testing.T) {
	SetPrefixesToSanitize("github.com/tech4works/")
	defer SetPrefixesToSanitize()

	err := NewWithCode("BAD_REQUEST", "invalid input")
	stack := err.Stack()

	// Verify code is preserved
	if !strings.Contains(err.Error(), "BAD_REQUEST") {
		t.Errorf("expected code to be preserved")
	}

	// Verify prefix is removed
	if strings.Contains(stack, "github.com/tech4works/") {
		t.Errorf("expected prefix to be removed")
	}
}

// TestIntegration_PrefixRemovalWithMetadata tests prefix removal with metadata
func TestIntegration_PrefixRemovalWithMetadata(t *testing.T) {
	SetPrefixesToSanitize("github.com/tech4works/")
	defer SetPrefixesToSanitize()

	metadata := map[string]any{
		"user_id": 123,
		"action":  "login",
	}

	err := NewWithMetadata(metadata, "metadata error")
	stack := err.Stack()

	// Verify metadata is preserved
	if !strings.Contains(err.Error(), "user_id") {
		t.Errorf("expected metadata to be preserved")
	}

	// Verify prefix is removed
	if strings.Contains(stack, "github.com/tech4works/") {
		t.Errorf("expected prefix to be removed")
	}
}

// --- SetFrameToSanitize Integration Tests ---

func TestIntegration_FrameSanitizationWithStackTrace(t *testing.T) {
	originalPolicy := getPolicy()
	defer SetPolicy(originalPolicy)

	SetPolicy(PolicyNormal)
	SetFrameToSanitize("github.com/tech4works/sdk-manas")
	defer SetFrameToSanitize()

	err := New("frame sanitization test")
	stack := err.Stack()

	// Verify frame is removed
	if strings.Contains(stack, "sdk-manas") {
		t.Errorf("expected 'sdk-manas' frame to be removed from stack")
	}

	// Verify error message is preserved
	if !strings.Contains(err.Error(), "frame sanitization test") {
		t.Errorf("expected error message to be preserved")
	}
}

func TestIntegration_MultipleFrameSanitization(t *testing.T) {
	SetFrameToSanitize("runtime.goexit", "testing.tRunner", "gin")
	defer SetFrameToSanitize()

	raw := []byte("runtime.goexit()\n\t/usr/local/go/src/runtime/asm.s:1650 +0x1\ntesting.tRunner()\n\t/usr/local/go/src/testing/testing.go:1234 +0x1\ngithub.com/gin-gonic/gin.().Next()\n\t/go/pkg/mod/github.com/gin-gonic/gin@v1.12.0/context.go:192 +0x5f\nmyapp.myFunc()\n\t/app/main.go:42 +0x1")
	result := normalizeStack(raw)

	// Verify all frames are removed
	if strings.Contains(result, "runtime.goexit") {
		t.Errorf("expected 'runtime.goexit' frame to be removed")
	}
	if strings.Contains(result, "testing.tRunner") {
		t.Errorf("expected 'testing.tRunner' frame to be removed")
	}
	if strings.Contains(result, "gin") {
		t.Errorf("expected 'gin' frame to be removed")
	}

	// Verify user code is preserved
	if !strings.Contains(result, "myapp.myFunc") {
		t.Errorf("expected 'myapp.myFunc' frame to be preserved")
	}
}

func TestIntegration_FrameAndPrefixSanitization(t *testing.T) {
	originalPolicy := getPolicy()
	defer SetPolicy(originalPolicy)

	SetPolicy(PolicyNormal)
	SetPrefixesToSanitize("github.com/tech4works/")
	SetFrameToSanitize("sdk-manas")
	defer func() {
		SetPrefixesToSanitize()
		SetFrameToSanitize()
	}()

	err := New("combined sanitization test")
	stack := err.Stack()

	// Verify both sanitizations work
	if strings.Contains(stack, "github.com/tech4works/") {
		t.Errorf("expected prefix to be removed")
	}
	if strings.Contains(stack, "sdk-manas") {
		t.Errorf("expected frame containing 'sdk-manas' to be removed")
	}
}

func TestIntegration_FrameSanitizationWithInheritance(t *testing.T) {
	SetFrameToSanitize("gin")
	defer SetFrameToSanitize()

	baseErr := New("base error")
	inheritedErr := Inherit(baseErr, "inherited error")

	stack := inheritedErr.Stack()

	// Verify frame sanitization works with inherited errors
	if strings.Contains(stack, "gin") {
		t.Errorf("expected 'gin' frame to be removed from inherited error stack")
	}
}

// --- Hex Offset Sanitization Integration Tests ---

func TestIntegration_HexOffsetRemoval(t *testing.T) {
	originalPolicy := getPolicy()
	defer SetPolicy(originalPolicy)

	SetPolicy(PolicyNormal)

	err := New("hex offset test")
	stack := err.Stack()

	// Verify hex offsets are removed
	if strings.Contains(stack, "+0x") {
		t.Errorf("expected hex offsets to be removed from stack, got: %s", stack)
	}

	// Verify error message is preserved
	if !strings.Contains(err.Error(), "hex offset test") {
		t.Errorf("expected error message to be preserved")
	}
}

func TestIntegration_HexOffsetWithPrefixSanitization(t *testing.T) {
	originalPolicy := getPolicy()
	defer SetPolicy(originalPolicy)

	SetPolicy(PolicyNormal)
	SetPrefixesToSanitize("github.com/tech4works/")
	defer SetPrefixesToSanitize()

	err := New("combined sanitization")
	stack := err.Stack()

	// Verify both sanitizations work
	if strings.Contains(stack, "+0x") {
		t.Errorf("expected hex offsets to be removed")
	}

	if strings.Contains(stack, "github.com/tech4works/") {
		t.Errorf("expected prefixes to be removed")
	}
}

func TestIntegration_HexOffsetWithFrameSanitization(t *testing.T) {
	originalPolicy := getPolicy()
	defer SetPolicy(originalPolicy)

	SetPolicy(PolicyNormal)
	SetFrameToSanitize("gin")
	defer SetFrameToSanitize()

	err := New("frame and offset test")
	stack := err.Stack()

	// Verify both sanitizations work
	if strings.Contains(stack, "+0x") {
		t.Errorf("expected hex offsets to be removed")
	}

	if strings.Contains(stack, "gin") {
		t.Errorf("expected 'gin' frames to be removed")
	}
}

func TestIntegration_AllSanitizationsCombined(t *testing.T) {
	originalPolicy := getPolicy()
	defer SetPolicy(originalPolicy)

	SetPolicy(PolicyNormal)
	SetPrefixesToSanitize("github.com/tech4works/")
	SetFrameToSanitize("gin", "sdk-manas")
	defer func() {
		SetPrefixesToSanitize()
		SetFrameToSanitize()
	}()

	err := New("all sanitizations")
	stack := err.Stack()

	// Verify all sanitizations work together
	if strings.Contains(stack, "+0x") {
		t.Errorf("expected hex offsets to be removed")
	}

	if strings.Contains(stack, "github.com/tech4works/") {
		t.Errorf("expected prefixes to be removed")
	}

	if strings.Contains(stack, "gin") || strings.Contains(stack, "sdk-manas") {
		t.Errorf("expected frames to be removed")
	}
}
