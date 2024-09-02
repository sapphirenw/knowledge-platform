package slogger

import (
	"fmt"
	"log/slog"
	"strings"
)

// some default error types
const (
	// basic
	ErrBasic     = "err-basic"
	ErrServer    = "err-server"
	ErrClient    = "err-client"
	ErrArguments = "err-arguments"
	ErrUnknown   = "err-unknown"

	// database
	ErrDBConn    = "err-db-conn"
	ErrDBExecute = "err-db-execute"
	ErrConflict  = "err-db-conflic"

	// perform custom parsing on the error message to see what went wrong
	ErrParse = "parse"
)

// parses an `ErrType` from an error message
func ParseType(err error) string {
	msg := err.Error()
	switch {
	// basic
	case strings.Contains(msg, "There was an issue parsing the request body"):
		return "err-arguments"

	case strings.Contains(msg, "failed to connect to"):
		return "err-db-conn"
	case strings.Contains(msg, "violates unique constraint"):
		return "err-db-conflict"
	default:
		return "err-unkown"
	}
}

// Type that wraps an error in a linked list with a type for stack history
type Err struct {
	t       string // type of err
	message string // message of this error
	err     error  // message for the err
	prev    *Err   // pointer to prev err
	args    []any  // any internal args that may need to be retrieved
}

// Wraps the error in a new error object and prints the message and args
func NewErr(
	logger *slog.Logger,
	message string,
	err error,
	t string,
	args ...any,
) *Err {
	if logger == nil {
		logger = slog.Default()
	}

	// parse the error if necessary
	if t == "" || t == "parse" {
		t = ParseType(err)
	}

	// handle the stack-based error
	new := &Err{
		t:       t,
		message: message,
		err:     err,
		prev:    nil,
		args:    args,
	}
	new.Print(logger)

	return new
}

// wrap the current error in a new error
func (e *Err) NewErr(
	logger *slog.Logger,
	message string,
	err error,
	t string,
	args ...any,
) *Err {
	new := NewErr(logger, message, err, t, args...)
	new.prev = e
	return new
}

// Checks if the stack contains the requested error type
func (e *Err) Contains(t string) bool {
	curr := e

	// check the ll for the err type
	for {
		if curr == nil {
			return false
		}
		if curr.t == t {
			return true
		}
		curr = e.prev
	}
}

// Prints the current error to the specified logger
func (e *Err) Print(logger *slog.Logger) {
	if logger == nil {
		logger = slog.Default()
	}

	logger.Error(fmt.Sprintf("type=%s message=%s error=%s", e.t, e.message, e.err.Error()), e.args...)
}

// Prints the entire stack without using the args convention in slog, instead prints out logfmt messages
func (e *Err) PrintStack(logger *slog.Logger) {
	if logger == nil {
		logger = slog.Default()
	}

	stack := e.Stack()
	for _, item := range stack {
		logger.Error(item)
	}
}

// implementation of the error interface that writes the entire error stack to a single string, separated by "\n"
func (e *Err) Error() string {

	// get the stack
	stack := e.Stack()

	// write the buffer
	buf := new(strings.Builder)
	for i, item := range stack {
		if i != 0 {
			buf.WriteString("\n")
		}
		buf.WriteString(item)
	}

	return buf.String()
}

// Writes the current Err to a string. Does NOT write the stack.
// If you want a stack of the error messages, use either `e.Error()` or `e.Stack()`
func (e *Err) String() string {
	buf := new(strings.Builder)

	// write the current message
	buf.WriteString(fmt.Sprintf("type=%s message=%s error=%s", e.t, e.message, e.err.Error()))
	if len(e.args) == 0 {
		return buf.String()
	}

	// write the arguments
	// should be true, but catch the case
	if len(e.args)%2 == 0 {
		i := 0

		for {
			if i >= len(e.args) {
				break
			}
			buf.WriteString(fmt.Sprintf(" %v=%v", e.args[i], e.args[i+1]))
			i += 2
		}
	} else {
		buf.WriteString(fmt.Sprintf(" args=%v", e.args))
	}

	return buf.String()
}

// Returns the stack of formatted error messages in errfmt format
func (e *Err) Stack() []string {
	list := make([]string, 0)

	curr := e

	for {
		if curr == nil {
			return list
		}

		list = append(list, curr.String())
		curr = e.prev
	}
}
