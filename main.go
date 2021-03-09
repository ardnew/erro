package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ardnew/erro/cmd"

	"github.com/ardnew/version"
)

// variables exported by Makefile via Go linker
var (
	PROJECT   string
	VERSION   string
	BRANCH    string
	REVISION  string
	BUILDTIME string
	PLATFORM  string
)

func init() {
	version.ChangeLog = []version.Change{{
		Package: "erro",
		Version: "0.1.0",
		Date:    "2021-03-09T20:05:09Z",
		Description: []string{
			"initial implementation",
		},
	}, {
		Package: PROJECT,
		Version: VERSION,
		Date:    BUILDTIME,
		Description: []string{
			"implemented proper support for escape sequences",
		},
	}}
}

func main() {
	var code int
	switch err := cmd.Erro(flag.CommandLine, os.Args...).(type) {
	case *cmd.ErrParseFlags:
		fmt.Fprintln(os.Stderr, err.Error())
		code = err.Code()
	case *cmd.ErrPrintStream:
		fmt.Fprintln(os.Stderr, err.Error())
		code = err.Code()
	}
	os.Exit(code)
}
