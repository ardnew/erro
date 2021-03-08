package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ardnew/version"
)

// Type definitions for various errors raised by cmd package.
type (
	ErrParseFlags struct {
		err error
		arg []string
	}
	ErrIOError string
)

// Error returns the string representation of ErrParseFlags
func (e ErrParseFlags) Error() string {
	a := enquote(e.arg...)
	return fmt.Sprintf("parse: %s: [%s]", e.err.Error(), strings.Join(a, ","))
}

func (e ErrIOError) Error() string {
	return fmt.Sprintf("failed to write to stderr: %s", string(e))
}

func usage(set *flag.FlagSet, separated bool) {
	exe := filepath.Base(executablePath())
	if separated {
		fmt.Fprintln(os.Stderr, "--")
	}
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  "+exe, "[options]", "[args ...]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "options:")
	fmt.Fprintln(os.Stderr, "  -v")
	fmt.Fprintln(os.Stderr, "  	Display version information")
	fmt.Fprintln(os.Stderr, "  -V")
	fmt.Fprintln(os.Stderr, "  	Display change history")
	fmt.Fprintln(os.Stderr, "  -n")
	fmt.Fprintln(os.Stderr, "  	Do not output a trailing newline")
	fmt.Fprintln(os.Stderr, "  -e")
	fmt.Fprintln(os.Stderr, "  	Enable interpretation of backslash escapes")
	fmt.Fprintln(os.Stderr, "  -E")
	fmt.Fprintln(os.Stderr, "  	Disable interpretation of backslash escapes (default true)")
	fmt.Fprintln(os.Stderr, "  -f format")
	fmt.Fprintln(os.Stderr, "  	Format output string according to format, where \"{N}\" represents argument N")
}

// Erro writes the given arguments to stderr according to command-line flags.
func Erro(set *flag.FlagSet, arg ...string) error {

	var (
		argVersion bool
		argChanges bool
		argNewline bool
		argEscapes bool
		argLiteral bool
		argFormats string
	)

	args := []string{}
	if len(arg) > 1 {
		args = append(args, arg[1:]...)
	}

	set.BoolVar(&argVersion, "v", false, "")
	set.BoolVar(&argChanges, "V", false, "")
	set.BoolVar(&argNewline, "n", false, "")
	set.BoolVar(&argEscapes, "e", false, "")
	set.BoolVar(&argLiteral, "E", true, "")
	set.StringVar(&argFormats, "f", "", "")
	set.Usage = func() { usage(set, false) }
	if err := set.Parse(args); nil != err {
		return ErrParseFlags{err: err, arg: args}
	}

	if argChanges {
		version.PrintChangeLog()
	} else if argVersion {
		fmt.Printf("erro version %s\n", version.String())
	} else {

		given := flagsProvided(set)

		_, escGiven := given["e"]
		_, litGiven := given["E"]

		literal := argEscapes == argLiteral && (litGiven && escGiven) &&
			argFormats == ""

		var output string
		if literal {
			output = strings.Join(set.Args(), " ")
		} else {
			if "" != argFormats {
				output = argFormats
				for i, s := range set.Args() {
					output = strings.ReplaceAll(output, fmt.Sprintf("{%d}", i), s)
				}
			} else {
				output = fmt.Sprintf("%s", strings.Join(set.Args(), " "))
			}
			output = fmt.Sprintf("%s", output) // expand any remaining escapes
		}

		var err error
		if argNewline {
			_, err = fmt.Fprint(os.Stderr, output)
		} else {
			_, err = fmt.Fprintln(os.Stderr, output)
		}
		if nil != err {
			return ErrIOError(err.Error())
		}
	}
	return nil
}

func executablePath() string {
	exe, err := os.Executable()
	if nil != err {
		panic("error: cannot determine executable: " + err.Error())
	}
	return exe
}

func flagsProvided(set *flag.FlagSet) map[string]flag.Value {
	m := map[string]flag.Value{}
	set.Visit(func(f *flag.Flag) { m[f.Name] = f.Value })
	return m
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
