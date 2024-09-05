package errors

import (
	"errors"
	"regexp"
	"runtime/debug"
	"strings"
)

const regex = `\[CAUSE]: \(([^:]+):(\d+)\) ([^:]+): (.+?) \[STACK]:\s*([\s\S]+)`

// Is checks if the target error is the same as the error passed to it. If either err or target
// is of type Detail, it extracts the message and compares it with the other error. Returns
// true if err is not nil and the same as target, and false otherwise.
//
// Parameters:
//   - err: The actual error to be checked.
//   - target: The target error to compare with.
//
// Returns:
//   - bool: A boolean value indicating if err matches target error.
//
// Example:
//
//	err := errors.New("test error")
//	target := errors.New("test error")
//	fmt.Println(Is(err, target)) // true
//
//	errDetail := New("test error")
//	targetDetail := New("test error")
//	fmt.Println(Is(errDetail, targetDetail)) // true
//
//	fmt.Println(Is(nil, nil)) // false
func Is(err, target error) bool {
	if IsDetailed(err) {
		errDetails := Details(err)
		err = errors.New(errDetails.Message())
	}

	if IsDetailed(target) {
		errDetails := Details(target)
		target = errors.New(errDetails.Message())
	}

	return err != nil && target != nil && err.Error() == target.Error()
}

// IsNot checks if the target error is different from the error passed to it.
// This function is the inverse of the Is function. If the Is function returns
// true, indicating that the errors are the same, the IsNot function will return
// false, and vice versa.
//
// Parameters:
//   - err: The error that needs to be compared.
//   - target: The error with which the comparison is made.
//
// Returns:
//   - bool: A boolean value indicating if err does not match target error.
//
// Example:
//
//	err1 := errors.New("test error")
//	err2 := errors.New("another error")
//	fmt.Println(IsNot(err1, err2)) // returns: true
//
//	err3 := errors.New("test error")
//	err4 := errors.New("test error")
//	fmt.Println(IsNot(err3, err4)) // returns: false
func IsNot(err, target error) bool {
	return !Is(err, target)
}

// Contains determines whether the error message from the 'err' error is found
// within the error message from the 'target' error. It uses the 'IsDetailed' function
// to check if the errors are detailed, gets their messages using 'Details' function and checks
// if the error message of 'err' contains that of 'target'.
//
// Parameters:
//   - err: The error to be checked.
//   - target: The error to be searched in 'err'.
//
// Returns:
//   - bool: A boolean value indicating whether the error message of 'err' contains that of 'target'.
//
// Example:
//
//	err := errors.New("test")
//	target := New("test")
//	fmt.Println(Contains(err, target)) // true
//
//	errDetail := New("test")
//	targetDetail := New("test2")
//	fmt.Println(Contains(errDetail, targetDetail)) // false
func Contains(err, target error) bool {
	if IsDetailed(err) {
		errDetails := Details(err)
		err = errors.New(errDetails.Message())
	}

	if IsDetailed(target) {
		errDetails := Details(target)
		target = errors.New(errDetails.Message())
	}

	return err != nil && target != nil && strings.Contains(err.Error(), target.Error())
}

// NotContains checks if the error message from the 'err' error is not found within
// the error message from the 'target' error. This function is the inverse of the 'Contains' function.
//
// Parameters:
//   - err: The error to be checked.
//   - target: The error to be checked within 'err'.
//
// Returns:
//   - bool: A boolean value indicating whether the error message of 'err' does not contain that of 'target'.
//
// Example:
//
//	err := errors.New("test")
//	target := errors.New("test2")
//	fmt.Println(NotContains(err, target)) // true
//
//	errAnother := errors.New("sample")
//	targetAnother := errors.New("sample")
//	fmt.Println(NotContains(errAnother, targetAnother)) // false
func NotContains(err, target error) bool {
	return !Contains(err, target)
}

// IsDetailed checks if a given error matches a detailed error regex. If the error is not nil and it matches
// the regex pattern regexErrorDetail, it returns true; otherwise, it returns false.
//
// Parameters:
//   - err: The error to be checked against the detailed error regex.
//
// Returns:
//   - bool: A boolean value indicating whether the given error matches the detailed error regex.
//
// Example:
//
//	err := errors.New("[CAUSE]: (fileName.go:1) FuncName: Error occurred [STACK]: detailed stack")
//	fmt.Println(IsDetailed(err)) // true
//
//	err = errors.New("simple error")
//	fmt.Println(IsDetailed(err)) // false
func IsDetailed(err error) bool {
	regex := regexp.MustCompile(regex)
	return err != nil && regex.MatchString(err.Error())
}

// Details function extracts detailed information from an error
//
// This function works by extracting specific parts of the error message using regular expressions.
// The details extracted include the file name, line number, function name, the error message, and a debug stack trace.
// If the error does not match the expected format, the function uses runtime.Caller and debug.Stack to get the
// file info and debug stack, and creates an error message using buildMessage function. It then creates a new Detail object
// with these details and returns it. If the provided error is nil, function simply returns nil.
//
// Parameters:
//   - err: The error from which the details are to be extracted
//
// Returns:
//   - *Detail: A pointer to a Detail struct containing error details. Returns nil if the provided error is nil.
//
// Example:
//
//	package main
//
//	import (
//		"errors"
//		"fmt"
//	)
//
//	func main() {
//		err := errors.New("test error")
//		detail := Details(err)
//		fmt.Println(detail.Message()) // Outputs: test error
//	}
func Details(err error) *Detail {
	if err == nil {
		return nil
	}

	var file string
	var line string
	var funcName string
	var message string
	var stack string

	rg := regexp.MustCompile(regex)
	matches := rg.FindStringSubmatch(err.Error())

	if len(matches) > 0 {
		file = matches[1]
		line = matches[2]
		funcName = matches[3]
		message = matches[4]
		stack = matches[5]
	} else {
		file, line, funcName = callerInfos(2)
		stack = string(debug.Stack())
		message = buildMessage(err.Error())
	}

	return &Detail{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  message,
		stack:    stack,
	}
}

func Join(errs []error, sep string) (result string) {
	for i, err := range errs {
		dt := Details(err)
		result += dt.message
		if i < len(errs)-1 {
			result += sep
		}
	}
	return result
}

func JoinToErr(errs []error, sep string) error {
	msg := Join(errs, sep)
	return errors.New(msg)
}
