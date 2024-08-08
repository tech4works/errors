package errors

import (
	"errors"
	"testing"
)

func TestIs(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{"Error is nil", nil, New("test error"), false},
		{"Target is nil", New("test error"), nil, false},
		{"Both are nil", nil, nil, false},
		{"Errors are the same", New("test error"), New("test error"), true},
		{"Errors are different", New("test error 1"), New("test error 2"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.err, tt.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNot(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{"Error is nil", nil, New("test error"), true},
		{"Target is nil", New("test error"), nil, true},
		{"Both are nil", nil, nil, true},
		{"Errors are the same", New("test error"), New("test error"), false},
		{"Errors are different", New("test error 1"), New("test error 2"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNot(tt.err, tt.target); got != tt.want {
				t.Errorf("IsNot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{"Error is not nil and target is contained", New("test error target"), New("target"), true},
		{"Error is not nil and target is not contained", New("test error"), New("target"), false},
		{"Error is nil", nil, New("target"), false},
		{"Target is nil", New("test error"), nil, false},
		{"Both are nil", nil, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.err, tt.target); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotContains(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{"Error is not nil and target is contained", New("test error target"), New("target"), false},
		{"Error is not nil and target is not contained", New("test error"), New("target"), true},
		{"Error is nil", nil, New("target"), true},
		{"Target is nil", New("test error"), nil, true},
		{"Both are nil", nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotContains(tt.err, tt.target); got != tt.want {
				t.Errorf("NotContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDetailed(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"Error is nil", nil, false},
		{"Error is detailed", New("test error"), true},
		{"Error is not detailed", errors.New("test error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDetailed(tt.err); got != tt.want {
				t.Errorf("IsDetailed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetails(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"Error is nil", nil},
		{"Error is detailed", New("test error")},
		{"Error is not detailed", errors.New("test error")},
		// Add here the test for when the error is detailed once the regex and the way of creating detailed errors is defined
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Details(tt.err); (got != nil) != (tt.err != nil) {
				t.Errorf("Details() should be nil for non-detailed errors and nil errors")
			}
		})
	}
}
