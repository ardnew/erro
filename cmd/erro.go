package cmd

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/ardnew/version"
)

// Type definitions for various errors raised by cmd package.
type (
	ErrParseFlags struct {
		err error
		arg []string
	}
)

// Error returns the string representation of ErrParseFlags
func (e *ErrParseFlags) Error() string {
	a := enquote(e.arg...)
	return fmt.Sprintf("parse: %s: [ %s ]", e.err.Error(), strings.Join(a, ","))
}

func Erro(fs *flag.FlagSet, av ...string) error {

	var (
		argVersion bool
		argChanges bool
	)

	fs.BoolVar(&argVersion, "v", false, "Display version information")
	fs.BoolVar(&argChanges, "V", false, "Display change history")
	if err := fs.Parse(av); nil != err {
		return &ErrParseFlags{err: err, arg: av}
	}

	if argChanges {
		version.PrintChangeLog()
	} else if argVersion {
		fmt.Printf("erro version %s\n", version.String())
	} else {
		// main
	}

	return &ErrParseFlags{err: errors.New("some error"), arg: av}
}

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
