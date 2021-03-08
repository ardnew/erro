package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ardnew/erro/cmd"

	"github.com/ardnew/version"
)

func init() {
	version.ChangeLog = []version.Change{{
		Package: "erro",
		Version: "0.1.0",
		Date:    "2021 Mar 08",
		Description: []string{
			"initial implementation",
		},
	}}
}

func main() {
	var code int
	switch err := cmd.Erro(flag.CommandLine, os.Args...).(type) {
	case *cmd.ErrParseFlags:
		fmt.Println(err.Error())
		code = 1
	default:
		code = 0
	}
	os.Exit(code)
}
