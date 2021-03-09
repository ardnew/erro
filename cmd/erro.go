package cmd

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ardnew/version"
)

// Erro writes the given arguments to stderr using options provided in the given
// set of flags.
func Erro(set *flag.FlagSet, arg ...string) error {

	config, err := Parse(set, arg...)
	if nil != err {
		return err
	}

	if config.ShowChanges {
		version.PrintChangeLog()
	} else if config.ShowVersion {
		fmt.Printf("erro version %s\n", version.String())
	} else {

		var output string
		var cancel bool
		switch config.Mode {
		case Literal:
			output = strings.Join(config.Args, " ")
		case Escaped:
			output, cancel = escape(strings.Join(config.Args, " "))
		case Formats:
			output = config.Format
			for i, s := range config.Args {
				output = strings.ReplaceAll(output, fmt.Sprintf("{%d}", i), s)
			}
			output, cancel = escape(output)
		}

		print := fmt.Fprint
		if config.AddNewline && !cancel {
			print = fmt.Fprintln
		}

		if _, err := print(os.Stderr, output); nil != err {
			return &ErrPrintStream{err: err, str: output}
		}
	}

	return nil
}

// escape replaces all recognized escape sequences with their evaluated result.
// If the special escape sequence "\c" is found, all trailing text is truncated
// and the bool return value will be true. otherwise, the bool will be false.
func escape(str string) (string, bool) {
	escSeq := regexp.MustCompile(`\\[\\abcefnrtv]|\\0[0-7]{1,3}|\\x[0-9a-fA-F]{1,2}`)
	idx := escSeq.FindAllStringIndex(str, -1)
	if idx != nil {
		var out string
		for i, loc := range idx {
			beg := 0
			if i > 0 {
				beg = idx[i-1][1]
			}
			out += str[beg:loc[0]]
			switch str[loc[0]+1] {
			case '0':
				if val, err := strconv.ParseUint(str[loc[0]+2:loc[1]], 8, 8); nil == err {
					out += string([]byte{byte(val)})
				}
			case 'x':
				if val, err := strconv.ParseUint(str[loc[0]+2:loc[1]], 16, 8); nil == err {
					out += string([]byte{byte(val)})
				}
			case 'c':
				return out, true
			case '\\':
				out += "\\"
			case 'a':
				out += "\a"
			case 'b':
				out += "\b"
			case 'e':
				out += "\x1B"
			case 'f':
				out += "\f"
			case 'n':
				out += "\n"
			case 'r':
				out += "\r"
			case 't':
				out += "\t"
			case 'v':
				out += "\v"
			}
			if i+1 == len(idx) {
				out += str[loc[1]:]
			}
		}
		return out, false
	}
	return str, false
}
