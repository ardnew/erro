package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Usage prints the command's usage summary to stderr.
func Usage(set *flag.FlagSet, hzRule ...HorizontalRule) {
	exe := filepath.Base(ExecutablePath())
	if IsHorizontalRuleInList(hzRule, LeadingHorizontalRule) {
		fmt.Fprintln(os.Stderr, "--")
	}
	for _, ln := range []string{
		`usage:`,
		`  ` + exe + ` [options] [args ...]`,
		``,
		`options:`,
		`  -v       Display version information`,
		`  -V       Display change history`,
		`  -n       Do not output a trailing newline`,
		`  -e       Enable interpretation of backslash escapes`,
		`  -E       Disable interpretation of backslash escapes (default true)`,
		`             Accepted for compatibility with echo, but this flag is ignored.`,
		`  -f fmt   Output string according to fmt (implies -e)`,
		`             See formatting section below for details.`,
		``,
		`formatting:`,
		`  The fmt argument given to flag -f may contain special placeholder symbols`,
		`  of the form "{N}", meaning the N'th command-line argument. The indexing does`,
		`  not include flags or their arguments; it is strictly the N'th argument that`,
		`  would be printed to the standard error stream if no flags were given at all.`,
		``,
		`  Additionally, the -f flag enables interpretation of backslash sequences as`,
		`  if the -e flag was given.`,
		``,
		`  The following backslash escape sequences are recognized:`,
		``,
		`    \\     backslash                        \x5C`,
		`    \a     alert                            \x07  BEL`,
		`    \b     backspace                        \x08  BS`,
		`    \c     produce no further output`,
		`    \e     escape                           \x1B  ESC`,
		`    \f     form feed                        \x0C  FF`,
		`    \n     new line                         \x0A  LF`,
		`    \r     carriage return                  \x0D  CR`,
		`    \t     horizontal tab                   \x09  TAB`,
		`    \v     vertical tab                     \x0B  VT`,
		`    \0NNN  byte with octal value NNN        1 to 3 digits`,
		`    \xHH   byte with hexadecimal value HH   1 to 2 digits`,
	} {
		fmt.Fprintln(os.Stderr, ln)
	}
	if IsHorizontalRuleInList(hzRule, TrailingHorizontalRule) {
		fmt.Fprintln(os.Stderr, "--")
	}
}

// Parse returns a new Config struct from the given FlagSet and arguments.
func Parse(set *flag.FlagSet, arg ...string) (*Config, error) {

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
	set.Usage = func() { Usage(set) }
	if err := set.Parse(args); nil != err {
		return nil, &ErrParseFlags{err: err, arg: args}
	}

	mode := Literal
	if argFormats != "" {
		mode = Formats
	} else if argEscapes {
		mode = Escaped
	}

	return &Config{
		ShowVersion: argVersion,
		ShowChanges: argChanges,
		AddNewline:  !argNewline,
		Mode:        mode,
		Format:      argFormats,
		Args:        set.Args(),
	}, nil
}

// ExecutablePath returns an absolute path to the command executable.
// It will panic if the path cannot be determined.
func ExecutablePath() string {
	exe, err := os.Executable()
	if nil != err {
		panic("error: cannot determine executable: " + err.Error())
	}
	return exe
}

// FlagsProvided returns a map of all flags that were provided in the FlagSet.
// Map elements are the flag's Value (interface) keyed by flag name.
func FlagsProvided(set *flag.FlagSet) map[string]flag.Value {
	m := map[string]flag.Value{}
	set.Visit(func(f *flag.Flag) { m[f.Name] = f.Value })
	return m
}

// HorizontalRule represents an output line that serves as horizontal separator.
type HorizontalRule interface{ String() string }

// horizontalRule is the actual string used to represent a horizontal separator.
const horizontalRule = "--"

// Type definitions for various HorizontalRules whose positions are specified
// relative to other output.
type (
	leadingHorizontalRule  string
	trailingHorizontalRule string
)

// Constant string definitions of the various types of HorizontalRule.
const (
	LeadingHorizontalRule  leadingHorizontalRule  = horizontalRule
	TrailingHorizontalRule trailingHorizontalRule = horizontalRule
)

// String returns the string representation of a leadingHorizontalRule
func (h leadingHorizontalRule) String() string { return string(h) }

// String returns the string representation of a trailingHorizontalRule
func (h trailingHorizontalRule) String() string { return string(h) }

// IsHorizontalRuleInList returns true if and only if the type and value of the
// given rule are non-nil and are both equal to some element in list.
func IsHorizontalRuleInList(list []HorizontalRule, rule HorizontalRule) bool {
	// assume no list contains the nil HorizontalRule
	if nil != rule {
		for _, r := range list {
			if r == rule { // iff same type -and- same value
				return true
			}
		}
	}
	return false
}
