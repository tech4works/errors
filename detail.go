package errors

import (
	"fmt"
	"runtime/debug"
	"strconv"
)

type Detail struct {
	file     string
	line     string
	funcName string
	message  string
	stack    string
}

// New constructs a new error instance with detailed information.
// It builds an error message with the provided arguments, fetches caller information and stacks the debug info.
//
// Parameters:
//   - args: Variadic arguments of any type to be composed into an error message.
//
// Returns:
//   - error: An error instance wrapped with details including file name, line number, function name, message and debug stack.
//
// Example:
//
//	err := New("File not found")
//	fmt.Println(err) // Outputs error detail with message "File not found"
//
//	err = New("File not found", "in directory", "/home/user")
//	fmt.Println(err) // Outputs error detail with message "File not found in directory /home/user"
func New(args ...any) error {
	msg := buildMessage(args...)
	file, line, funcName := callerInfos(2)
	stack := debug.Stack()
	return &Detail{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  msg,
		stack:    string(stack),
	}
}

// Newf constructs a new error instance with detailed information. It accepts a format string and variadic arguments
// to build the error message. It fetches caller information and stacks the debug info.
//
// Parameters:
//   - format: A format as a string.
//   - args: Variadic arguments of any type to be composed into an error message following the provided format.
//
// Returns:
//   - error: An error instance wrapped with details including file name, line number, function name, message and debug stack.
//
// Example:
//
//	err := Newf("%s not found", "File")
//	fmt.Println(err) // Outputs error detail with message "File not found"
//
//	err = Newf("%s not found in %s", "File", "/home/user")
//	fmt.Println(err) // Outputs error detail with message "File not found in /home/user"
//
// Note:
// The function utilizes internal helper functions such as buildMessageByFormat, callerInfos, and debug.Stack()
// to construct the error message.
func Newf(format string, args ...any) error {
	msg := buildMessageByFormat(format, args...)
	file, line, funcName := callerInfos(2)
	stack := debug.Stack()
	return &Detail{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  msg,
		stack:    string(stack),
	}
}

// NewSkipCaller constructs a new Detail structure. It takes a variable amount of parameters
// to construct the error message incorporating caller details such as file, line number and function
// name for providing comprehensive error logging and stack trace details. The elements of `args` are
// formatted as strings using the %v verb and passed to the `buildMessage` function to construct the
// message string.
//
// Parameters:
//   - skipCaller: Integer specifying the number of stack frames to skip
//   - args: Variable arguments used in error message construction
//
// Returns:
//   - error: An error structure which provides comprehensive error logging and stack trace details.
//
// Panic:
//   - If buildMessage function calls a method from an uninitialized structure, a panic can occur.
//
// Example:
//
//	errorVariable := NewSkipCaller(1, "Incorrect operation.")
//	// Causes the Error() method of the error variable to be called producing output like:
//	// [CAUSE]: (file.go:50) funcName: Incorrect operation. [STACK]: Goroutine 23 - file.go:50
//	fmt.Println(errorVariable)
//
//	secondError := NewSkipCaller(5, "Error in processing.")
//	// Causes the Error() method of the error variable to be called producing different output
//	// due to the skipCaller value being different.
//	fmt.Println(secondError)
func NewSkipCaller(skipCaller int, args ...any) error {
	msg := buildMessage(args...)
	file, line, funcName := callerInfos(skipCaller + 1)
	stack := debug.Stack()
	return &Detail{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  msg,
		stack:    string(stack),
	}
}

// NewSkipCallerf creates a new error structured type `ErrorDetail`. It includes detailed logging information such as
// origin file, line number, and the function where the error occurred.
// The provided arguments are formatted and passed to the `buildMessageByFormat` function to finalize the error message.
// The function uses details from the calling stack, skipping a certain number of stack frames as specified by
// `skipCaller`, to inject further details.
//
// Parameters:
//   - skipCaller: Integer representing the number of stack frames to omit.
//   - format: String specifying the format for the error message.
//   - args: Variadic parameters used to construct the error message.
//
// Returns:
//   - error: A newly created error of type `ErrorDetail` with all the tracing details included.
//
// Panic:
//   - If the function `runtime.Caller()` used in the `callerInfos` function doesn't have enough stack frames to skip.
//
// Example:
//
//	err := NewSkipCallerf(1, "Division by zero at %s function.", "divide")
//	// Result:
//	// err.Error() outputs: "[CAUSE]: (file.go:50) divide: Division by zero at divide function. [STACK]: Goroutine 23 - file.go:50"
//	fmt.Println(err)
//
//	err = NewSkipCallerf(2, "Unexpected value in %s function: %v", "logValue", "nil")
//	// Result:
//	// err.Error() outputs: "[CAUSE]: (file.go:61) logValue: Unexpected value in logValue function: <nil>. [STACK]:
//	// Goroutine 24 - file.go:61"
//	fmt.Println(err)
func NewSkipCallerf(skipCaller int, format string, args ...any) error {
	msg := buildMessageByFormat(format, args...)
	file, line, funcName := callerInfos(skipCaller + 1)
	stack := debug.Stack()
	return &Detail{
		file:     file,
		line:     line,
		funcName: funcName,
		message:  msg,
		stack:    string(stack),
	}
}

// Error constructs a detailed error string containing the cause of the error and the
// debug stack. The string is formatted in such a way that it emphasizes the cause
// of the error and the corresponding debug stack for better readability in error
// logs or output.
//
// This method is typically used when detail-rich error messages are needed,
// especially during debugging or while logging errors as it helps identify the
// exact location (file and line number) and the function that might have resulted
// in the error.
//
// No Parameters.
//
// Returns:
//   - string: A detailed string representation of the error, formatted as
//     "[CAUSE]: <cause of the error> [STACK]: <debug stack>"
func (e *Detail) Error() string {
	return fmt.Sprint("[CAUSE]: ", e.Cause(), " [STACK]: ", e.stack)
}

// PrintStackTrace prints the debug stack of the Detail instance.
// This method can be used to output the debug stack for debugging purposes or
// logging the error. The debug stack contains information about the file,
// line number, and function name where the error occurred.
func (e *Detail) PrintStackTrace() {
	fmt.Print(e.stack)
}

// PrintCause prints the cause of the error represented by the Detail instance.
//
// This method can be used to output the cause of the error for debugging purposes
// or logging the error. The cause of the error contains information about the file,
// line number, and function name where the error occurred, along with the error message.
func (e *Detail) PrintCause() {
	fmt.Print(e.Cause())
}

// Cause returns a string representation of the cause of the error, including the file, line number,
// function name, and error message. It is typically used to display the cause of an error for debugging
// or logging purposes. The returned string has the format "(file:line) function: message".
//
// No parameters.
//
// Returns:
//   - string: A string representation of the cause of the error.
func (e *Detail) Cause() string {
	return fmt.Sprint("(", e.file, ":", e.line, ")", " ", e.funcName, ": ", e.message)
}

// Message returns the message associated with the Detail instance.
// This method can be used to extract only the error message from an Detail object.
//
// No parameters.
//
// Returns:
//   - string: The error message.
func (e *Detail) Message() string {
	return e.message
}

// File returns the file name associated with the Detail instance.
// This method can be used to retrieve the file name where the error occurred.
//
// No parameters.
//
// Returns:
//   - string: The file name where the error occurred.
func (e *Detail) File() string {
	return e.file
}

// Line returns the line number where the error occurred. It converts the line number,
// which is stored as a string in the Detail instance, to an integer and returns it.
//
// This method is useful for identifying the exact line number where an error occurred
// and can be used for debugging or logging purposes.
//
// Returns:
//   - int: The line number where the error occurred.
func (e *Detail) Line() int {
	lineNumber, _ := strconv.Atoi(e.line)
	return lineNumber
}

// Func returns the name of the function where the error occurred. This method is useful
// for identifying the function name that caused the error and can be used for debugging
// or logging purposes.
//
// Returns:
//   - string: The name of the function where the error occurred.
func (e *Detail) Func() string {
	return e.funcName
}

// Stack returns the debug stack associated with the Detail instance.
// This method can be used to retrieve the stack trace of the error for debugging
// or logging purposes.
//
// Returns:
//   - string: The debug stack trace of the error.
func (e *Detail) Stack() string {
	return e.stack
}
