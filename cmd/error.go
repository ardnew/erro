package cmd

import (
	"fmt"
	"strings"
)

// CmdError represents an extended error interface capable of also returning
// an error code.
type CmdError interface {
	Error() string
	Code() int
}

// Type definitions for all errors raised by the cmd package.
type (
	ErrParseFlags struct {
		err error
		arg []string
	}
	ErrPrintStream struct {
		err error
		str string
	}
)

// Constant values enumerating all errors raised by the cmd package.
const (
	errCode = 10 // starting index of all error codes

	errParseFlags = errCode + iota
	errPrintStream
)

// Error returns the string representation of ErrParseFlags.
func (e *ErrParseFlags) Error() string {
	a := enquote(e.arg...)
	return fmt.Sprintf("parse: %s: [%s]", e.err.Error(), strings.Join(a, ","))
}

// Code returns the error code representing ErrParseFlags.
func (e *ErrParseFlags) Code() int { return errParseFlags }

// Error returns the string representation of ErrPrintStream.
func (e *ErrPrintStream) Error() string {
	a := enquote(e.str)
	return fmt.Sprintf("output: %s: [%s]", e.err.Error(), a[0])
}

// Code returns the error code representing ErrPrintStream.
func (e *ErrPrintStream) Code() int { return errPrintStream }

func enquote(str ...string) []string {
	return mapEach(func(s string) string {
		return "'" + strings.ReplaceAll(s, "'", "\\'") + "'"
	}, str...)
}

func mapEach(fn func(string) string, str ...string) []string {
	out := make([]string, len(str))
	for i, s := range str {
		out[i] = fn(s)
	}
	return out
}
